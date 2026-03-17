package repository

import (
	"errors"
	"shoego/database"
	"shoego/domain"
	"shoego/models"

	"gorm.io/gorm"
)

func CheckingEmailValidation(email string) (*domain.User, error) {

	var user domain.User

	result := database.DB.Where(&domain.User{Email: email}).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil

		}
		return nil, result.Error
	}

	return &user, nil
}

func CheckingPhoneExists(phone string) (*domain.User, error) {

	var user domain.User
	result := database.DB.Where(&domain.User{Phone: phone}).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, result.Error
	}
	return &user, nil

}

func SignupInsert(user models.SignupDetail) (models.SignupDetailResponse, error) {
	dbUser := domain.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
		Phone:    user.Phone,
	}

	err := database.DB.Create(&dbUser).Error
	if err != nil {
		return models.SignupDetailResponse{}, err
	}

	return models.SignupDetailResponse{
		ID:    int(dbUser.ID),
		Name:  dbUser.Name,
		Email: dbUser.Email,
		Phone: dbUser.Phone,
	}, nil
}

func FindUserDetailByEmail(user models.UserLogin) (models.UserLoginResponse, error) {
	var UserDetails models.UserLoginResponse

	//get user for db
	err := database.DB.Raw(`SELECT * FROM users WHERE email = ? AND blocked = false LIMIT 1`, user.Email,).Scan(&UserDetails).Error

	if err != nil {
		return models.UserLoginResponse{}, errors.New("error searching users by email")
	}

	return UserDetails, nil
}

func FindUserByEmail(email string) (*domain.User, error) {
	var User domain.User

	result := database.DB.Where(&domain.User{Email: email}).First(&User)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {  
			return nil, nil
		}
		return nil, result.Error
	}
	return &User, nil
}

func CreateGoogleUser(user domain.User) (*domain.User, error) {
	//save user info into db 
	err := database.DB.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
