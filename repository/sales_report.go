package repository

import (
	"shoego/database"
	"shoego/domain"
	"time"
)

type SalesReport struct {
	Orders              []domain.Order
	TotalSalesCount     int64
	TotalOrderAmount    float64
	TotalOfferDiscount  float64
	TotalCouponDiscount float64
	NetSales            float64
}

func GetSalesReport(startDate, endDate time.Time) (*SalesReport, error) {
	var orders []domain.Order
	err := database.DB.Preload("User").Preload("OrderItems").Where("created_at BETWEEN ? AND ?", startDate, endDate).Where("payment_status = ?", "paid").Order("created_at DESC").Find(&orders).Error

	if err != nil {
		return nil, err
	}

	report := &SalesReport{
		Orders: orders,
	}

	for _, order := range orders {
		report.TotalSalesCount++
		report.TotalOrderAmount += order.Subtotal
		report.TotalOfferDiscount += order.OfferDiscount
		report.TotalCouponDiscount += order.CouponDiscount
		report.NetSales += order.FinalAmount
	}
	return report, nil
}
