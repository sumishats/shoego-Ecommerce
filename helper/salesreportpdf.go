package helper

import (
	"bytes"
	"fmt"
	"shoego/models"

	"github.com/jung-kurt/gofpdf"
)

func GenerateSalesReportPDF(report *models.SalesReportResponse) ([]byte, error) {

	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.AddPage()

	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(190, 10, "Shoego Sales Report")

	pdf.Ln(15)

	pdf.SetFont("Arial", "", 12)

	pdf.Cell(45, 8, "From Date")
	pdf.Cell(80, 8, ": "+report.FromDate)
	pdf.Ln(8)

	pdf.Cell(45, 8, "To Date")
	pdf.Cell(80, 8, ": "+report.ToDate)

	pdf.Ln(12)

	// Summary

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 8, "Summary")

	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)

	pdf.Cell(60, 8, fmt.Sprintf("Total Orders : %d", report.TotalOrders))
	pdf.Ln(8)

	pdf.Cell(60, 8, fmt.Sprintf("Total Products : %d", report.TotalProducts))
	pdf.Ln(8)

	pdf.Cell(60, 8, fmt.Sprintf("Gross Sales : %.2f", report.GrossSales))
	pdf.Ln(8)

	pdf.Cell(60, 8, fmt.Sprintf("Offer Discount : %.2f", report.OfferDiscount))
	pdf.Ln(8)

	pdf.Cell(60, 8, fmt.Sprintf("Coupon Discount : %.2f", report.CouponDiscount))
	pdf.Ln(8)

	pdf.Cell(60, 8, fmt.Sprintf("Total Discount : %.2f", report.TotalDiscount))
	pdf.Ln(8)

	pdf.Cell(60, 8, fmt.Sprintf("Net Sales : %.2f", report.NetSales))

	pdf.Ln(15)

	// Table Header 

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(55, 8, "Order ID")
	pdf.Cell(40, 8, "Customer")
	pdf.Cell(25, 8, "Amount")
	pdf.Cell(25, 8, "Discount")
	pdf.Cell(20, 8, "Net")
	pdf.Cell(25, 8, "Status")

	pdf.Ln(8)

	//Table Data

	pdf.SetFont("Arial", "", 9)

	for _, order := range report.Orders {

		
		pdf.Cell(55, 8, order.OrderID)
		pdf.Cell(40, 8, order.CustomerName)
		pdf.Cell(25, 8, fmt.Sprintf("%.2f", order.OrderAmount))
		pdf.Cell(25, 8, fmt.Sprintf("%.2f", order.TotalDiscount))
		pdf.Cell(20, 8, fmt.Sprintf("%.2f", order.NetAmount))

		pdf.Cell(25, 8, order.OrderStatus)

		pdf.Ln(8)
	}

	// Convert PDF to bytes

	var buffer bytes.Buffer

	err := pdf.Output(&buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
