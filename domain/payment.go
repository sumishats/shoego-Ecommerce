package domain

import "gorm.io/gorm"

type Payment struct {
	gorm.Model

	OrderID uint  `gorm:"not null"`
	Order   Order `gorm:"foreignKey:OrderID"`
	UserID  uint 

	Amount float64 `gorm:"not null"`

	PaymentMethod string `gorm:"type:varchar(50)"`
	Status        string `gorm:"type:varchar(50);default:'pending'"`

	RazorpayOrderID   string
	RazorpayPaymentID string
	RazorpaySignature string
}
