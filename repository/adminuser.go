package repository

import (
	"shoego/database"
	"shoego/domain"
)

func GetUsers(search string, limit int, offset int) ([]domain.User, error) {

	var users []domain.User

	query := database.DB.Model(&domain.User{})

	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	err := query.Order("created_at desc").Limit(limit).Offset(offset).Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func CountUsers(search string) (int64, error) {

	var count int64

	query := database.DB.Model(&domain.User{})

	if search != "" {
		query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func FindUserByID(id uint) (*domain.User, error) {
	//fetch user by id in db
	var user domain.User

	err := database.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

//update user blocked  in db
func BlockUser(id uint) error {
	return database.DB.Model(&domain.User{}).Where("id = ?", id).Update("blocked", true).Error
}

//update user unblock 
func UnblockUser(id uint) error {
	return database.DB.Model(&domain.User{}).Where("id = ?", id).Update("blocked", false).Error
}
