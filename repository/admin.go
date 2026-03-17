package repository

import (
	"errors"
	"shoego/database"
	"shoego/domain"
	"shoego/models"
)

func AdminLogin(admin models.AdminLogin) (domain.Admin, error) {
	var admindomain domain.Admin

	if err := database.DB.Raw("select * from users where email= ? and is_admin= true ", admin.Email).Scan(&admindomain).Error; err != nil {
		return domain.Admin{}, errors.New("admin email is not available on database")
	}

	return admindomain, nil
}
