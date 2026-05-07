package domain

import "gorm.io/gorm"

type Wishlist struct {
	gorm.Model
	UserID    uint    `gorm:"not null;index"`
	User      User    `gorm:"foreignKey:UserID"`
	ProductID uint    `gorm:"not null;index"`
	Product   Product `gorm:"foreignKey:ProductID"`
}
