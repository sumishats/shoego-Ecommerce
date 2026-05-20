package domain

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	OrderID            string  `gorm:"uniqueIndex;not null"`
	UserID             uint    `gorm:"not null"`
	User               User    `gorm:"foreignKey:UserID"`
	AddressID          uint    `gorm:"not null"`
	Address            Address `gorm:"foreignKey:AddressID"`
	OrderStatus        string  `gorm:"type:varchar(50);default:'placed'"`
	PaymentMethod      string  `gorm:"type:varchar(50)"`
	PaymentStatus      string  `gorm:"type:varchar(50);default:'pending'"`
	Subtotal           float64
	TaxAmount          float64
	DiscountAmount     float64
	ShippingCharge     float64
	FinalAmount        float64
	CancellationReason string      `gorm:"type:text"`
	ReturnReason       string      `gorm:"type:text"`
	OrderItems         []OrderItem `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	gorm.Model
	OrderID            uint    `gorm:"not null"`
	ProductID          uint    `gorm:"not null"`
	Product            Product `gorm:"foreignKey:ProductID"`
	Quantity           int     `gorm:"not null"`
	Price              float64 `gorm:"not null"`
	TotalPrice         float64 `gorm:"not null"`
	ItemStatus         string  `gorm:"type:varchar(50);default:'placed'"`
	CancellationReason string  `gorm:"type:text"`
	ReturnReason       string  `gorm:"type:text"`
}
