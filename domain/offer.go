package domain

import (
	"time"

	"gorm.io/gorm"
)

type ProductOffer struct {
	gorm.Model

	ProductID uint    `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID"`

	OfferName string

	DiscountPercentage float64

	StartDate time.Time
	EndDate   time.Time

	IsActive bool `gorm:"default:true"`
}

type CategoryOffer struct {
	gorm.Model

	CategoryID uint
	Category   Category `gorm:"foreignKey:CategoryID"`

	OfferName string

	DiscountPercentage float64

	StartDate time.Time
	EndDate   time.Time

	IsActive bool
}
type Referral struct {
	gorm.Model
	UserID       uint
	ReferralCode string
	ReferredByID *uint
	IsRewarded   bool
}