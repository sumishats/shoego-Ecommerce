package domain

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID           uint   `json:"id" gorm:"unique;not null"`
	Name         string `json:"name"`
	Email        string `json:"email" gorm:"unique;not null"`
	Password     string `json:"password"`
	Phone        string `json:"phone"`
	ProfileImage string `json:"profile_image"`
	Blocked      bool   `json:"blocked" gorm:"default:false"`
	IsAdmin      bool   `json:"is_admin" gorm:"default:false"`
}

type Address struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `json:"user_id"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	HouseName string `json:"house_name"`
	Street    string `json:"street"`
	City      string `json:"city"`
	State     string `json:"state"`
	Pincode   string `json:"pincode"`
	IsDefault bool   `json:"is_default"`
	
}
type BlacklistToken struct {
	gorm.Model
	Token string `json:"token" gorm:"unique;not null"`
}
