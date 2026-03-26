package domain

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserID uint
	Items  []CartItem `gorm:"foreignKey:CartID"`
}

type CartItem struct {
	gorm.Model
	CartID    uint
	ProductID uint
	Quantity  int
	Product Product `gorm:"foreignKey:ProductID"`
}
