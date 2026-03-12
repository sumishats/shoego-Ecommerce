package domain

import "time"

type OTPVerification struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string
	OTP       string
	Name      string
	Phone     string
	Password  string
	Type      string
	ExpiresAt time.Time
	CreatedAt time.Time
}