package domain

import (
	"time"

	"gorm.io/gorm"
)

type Coupon struct {
	gorm.Model
	Code           string    `gorm:"unique;not null"`
	DiscountAmount float64   `gorm:"not null"`
	MinimumAmount  float64   `gorm:"not null"`
	ExpiryDate     time.Time `gorm:"not null"`
	UsageLimit     int       `gorm:"not null"`
	UsedCount      int       `gorm:"default:0"`
	IsActive       bool      `gorm:"default:true"`
}