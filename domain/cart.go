package domain

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID     uint
	CouponCode string
	Items []CartItem
}


type CartItem struct {
	gorm.Model

	CartID    uint
	ProductID uint
	VariantID *uint
	Quantity int
	Product Product
	Variant ProductVariant `gorm:"foreignKey:VariantID"`
}
