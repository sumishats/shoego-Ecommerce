package domain

import "shoego/models"

type TokenAdmin struct {
	Admin models.AdminDetailsResponse
	Token string
}

type Admin struct {
	ID       uint   `json:"id" gorm:"unique;not null"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"validate:required"`
	Password string `json:"password" gorm:"validate:required"`
}
