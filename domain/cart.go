package domain

import "gorm.io/gorm"

// type Cart struct {
// 	gorm.Model
// 	UserID uint
// 	Items  []CartItem `gorm:"foreignKey:CartID"`

// 	AppliedCouponID *uint
// 	DiscountAmount  float64
// 	FinalAmount     float64
// }

type Cart struct {
	gorm.Model
	UserID        uint
	CouponCode    string
	DiscountAmount float64
	Items         []CartItem
}

type CartItem struct {
	gorm.Model
	CartID    uint
	ProductID uint
	Quantity  int
	Product Product `gorm:"foreignKey:ProductID"`
}
