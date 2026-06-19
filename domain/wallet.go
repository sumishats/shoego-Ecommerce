package domain

import "gorm.io/gorm"

type Wallet struct{
	gorm.Model
	UserID uint `gorm:"uniqueIndex;not null"`
	Balance float64 `gorm:"default:0"`

}

type WalletTransaction struct{
	gorm.Model

	WalletID uint
	Amount float64
	Type string
	Description string

}