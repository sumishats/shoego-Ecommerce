package helper

import (
	"bytes"
	"fmt"
	"shoego/domain"
	"strconv"

	"github.com/jung-kurt/gofpdf"
)

func formatUint(id uint) string {
	return strconv.FormatUint(uint64(id), 10)
}

//create pdf
func GenerateInvoicePDF(order domain.Order) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Shoego Invoice")

	pdf.Ln(12)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(60, 10, fmt.Sprintf("Order ID: %s", order.OrderID))
	pdf.Ln(8)
	pdf.Cell(60, 10, fmt.Sprintf("Order Date: %s", order.CreatedAt.Format("2006-01-02")))
	pdf.Ln(8)
	pdf.Cell(60, 10, fmt.Sprintf("Payment Method: %s", order.PaymentMethod))
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(70, 10, "Product")
	pdf.Cell(30, 10, "Qty")
	pdf.Cell(40, 10, "Price")
	pdf.Cell(40, 10, "Total")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	for _, item := range order.OrderItems {
		if item.ItemStatus=="cancelled"{
			continue
		}
		pdf.Cell(70, 10, item.Product.Name)
		pdf.Cell(30, 10, fmt.Sprintf("%d", item.Quantity))
		pdf.Cell(40, 10, fmt.Sprintf("%.2f", item.Price))
		pdf.Cell(40, 10, fmt.Sprintf("%.2f", item.TotalPrice))
		pdf.Ln(10)
	}

	pdf.Ln(10)
	pdf.Cell(60, 10, fmt.Sprintf("Subtotal: %.2f", order.Subtotal))
	pdf.Ln(8)
	pdf.Cell(60, 10, fmt.Sprintf("Tax: %.2f", order.TaxAmount))
	pdf.Ln(8)
	pdf.Cell(60, 10, fmt.Sprintf("Discount: %.2f", order.CouponDiscount))
	pdf.Ln(8)
	pdf.Cell(60, 10, fmt.Sprintf("Shipping: %.2f", order.ShippingCharge))
	pdf.Ln(8)
	pdf.Cell(60, 10, fmt.Sprintf("Final Amount: %.2f", order.FinalAmount))

	//pdt covert to bytes
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
