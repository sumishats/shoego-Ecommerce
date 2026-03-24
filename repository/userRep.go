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
	err := database.DB.Raw(`SELECT * FROM users WHERE email = ? AND blocked = false LIMIT 1`, user.Email).Scan(&UserDetails).Error

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

func GetUserByID(userID uint) (*domain.User, error) {
	//get user info by id from db
	var user domain.User
	err := database.DB.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUserProfile(userID uint, name, phone, profileImage string) error {
	//send edited  info to db
	return database.DB.Model(&domain.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"name":          name,
		"phone":         phone,
		"profile_image": profileImage,
	}).Error
}

// if email is valid store to db
func UpdateUserEmail(userID uint, newEmail string) error {
	return database.DB.Model(&domain.User{}).Where("id = ?", userID).Update("email", newEmail).Error
}

// send update password info to db
func UpdateUserPassword(userID uint, hashedPassword string) error {
	return database.DB.Model(&domain.User{}).Where("id = ?", userID).Update("password", hashedPassword).Error
}

// insert new Address  in db
func AddAddress(address domain.Address) error {
	return database.DB.Create(&address).Error
}

func GetAddressesByUserID(userID uint) ([]domain.Address, error) {
	var addresses []domain.Address
	err := database.DB.Where("user_id = ?", userID).Find(&addresses).Error
	if err != nil {
		return nil, err
	}
	return addresses, nil
}

func GetAddressByID(addressID uint) (*domain.Address, error) {
	var address domain.Address
	err := database.DB.First(&address, addressID).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func UpdateAddress(addressID uint, data map[string]interface{}) error {
	return database.DB.Model(&domain.Address{}).Where("id = ?", addressID).Updates(data).Error
}

func DeleteAddress(addressID uint) error {
	return database.DB.Delete(&domain.Address{}, addressID).Error
}

// user set new default address ,clear old default address
func ClearDefaultAddresses(userID uint) error {

	return database.DB.Model(&domain.Address{}).Where("user_id = ?", userID).Update("is_default", false).Error
}
