package domain

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name          string
	Description   string
	BrandID       uint
	SKU           string
	Price         float64
	Stock         int
	CategoryID    uint
	Category      Category `gorm:"foreignKey:CategoryID"`
	IsListed      bool
	Images        []ProductImage `gorm:"foreignKey:ProductID"`
	ProductOffers []ProductOffer `gorm:"foreignKey:ProductID"`
	Variants []ProductVariant `gorm:"foreignKey:ProductID"`
}

type ProductImage struct {
	gorm.Model
	ProductID uint
	ImageURL  string
}

type ProductVariant struct {
	gorm.Model

	ProductID uint
	Size string
	SKU string
	Stock int
	Product Product `gorm:"foreignKey:ProductID"`
}

type Category struct {
	gorm.Model
	Name           string `gorm:"unique;not null"`
	Description    string
	IsListed       bool            `gorm:"default:true"`
	CategoryOffers []CategoryOffer `gorm:"foreignKey:CategoryID"`
}
