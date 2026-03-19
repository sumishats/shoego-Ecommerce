package domain

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string
	Description string
	BrandID     uint
	SKU         string
	Price       float64
	Stock       int
	CategoryID  uint
	Category    Category       `gorm:"foreignKey:CategoryID"`
	IsListed    bool 
	Images      []ProductImage `gorm:"foreignKey:ProductID"`
}

type ProductImage struct {
	gorm.Model
	ProductID uint
	ImageURL  string
}

type Category struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string
	IsListed    bool `gorm:"default:true"`
}