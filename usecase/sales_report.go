package usecase

import (
	"errors"
	"shoego/helper"
	"shoego/models"
	"shoego/repository"
	"time"
)

func GetSalesReport(filter, startDateStr, endDateStr string) (*models.SalesReportResponse, error) {

	var startDate time.Time
	var endDate time.Time

	today := time.Now()

	switch filter {
	case "", "all":

		startDate = time.Date(2000, 1, 1, 0, 0, 0, 0, today.Location())
		endDate = time.Date(today.Year(), today.Month(), today.Day(), 23, 59, 59, 999999999, today.Location())

	case "daily":

		startDate = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

		endDate = startDate.Add(24*time.Hour - time.Nanosecond)

	case "weekly":

		weekday := int(today.Weekday())

		if weekday == 0 {
			weekday = 7
		}

		startDate = time.Date(today.Year(), today.Month(), today.Day()-weekday+1, 0, 0, 0, 0, today.Location())

		endDate = startDate.AddDate(0, 0, 7).Add(-time.Nanosecond)

	case "monthly":

		startDate = time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, today.Location())

		endDate = startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

	case "yearly":

		startDate = time.Date(today.Year(), 1, 1, 0, 0, 0, 0, today.Location())

		endDate = time.Date(today.Year(), 12, 31, 23, 59, 59, 0, today.Location())

	case "custom":

		var err error

		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return nil, errors.New("invalid start date")
		}

		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return nil, errors.New("invalid end date")
		}

		endDate = endDate.Add(24*time.Hour - time.Nanosecond)

	default:
		return nil, errors.New("invalid filter")
	}

	report, err := repository.GetSalesReport(startDate, endDate)
	if err != nil {
		return nil, err
	}

	response := &models.SalesReportResponse{
		FromDate:       startDate.Format("2006-01-02"),
		ToDate:         endDate.Format("2006-01-02"),
		TotalOrders:    report.TotalSalesCount,
		GrossSales:     report.TotalOrderAmount,
		OfferDiscount:  report.TotalOfferDiscount,
		CouponDiscount: report.TotalCouponDiscount,
		TotalDiscount:  report.TotalOfferDiscount + report.TotalCouponDiscount,
		NetSales:       report.NetSales,
	}

	for _, order := range report.Orders {

		totalProducts := 0

		for _, item := range order.OrderItems {
			totalProducts += item.Quantity
		}

		response.TotalProducts += totalProducts

		response.Orders = append(response.Orders, models.SalesReportOrderResponse{
			OrderID:        order.OrderID,
			OrderDate:      order.CreatedAt,
			CustomerName:   order.User.Name,
			PaymentMethod:  order.PaymentMethod,
			OrderStatus:    order.OrderStatus,
			TotalProducts:  totalProducts,
			OrderAmount:    order.Subtotal,
			OfferDiscount:  order.OfferDiscount,
			CouponDiscount: order.CouponDiscount,
			TotalDiscount:  order.OfferDiscount + order.CouponDiscount,
			NetAmount:      order.FinalAmount,
		})
	}

	return response, nil
}

func DownloadSalesReportPDF(period, from, to string) ([]byte, error) {

	report, err := GetSalesReport(period, from, to)
	if err != nil {
		return nil, err
	}

	pdfBytes, err := helper.GenerateSalesReportPDF(report)
	if err != nil {
		return nil, err
	}

	return pdfBytes, nil
}
