package repository

import (
	"errors"
	"fmt"
	"shoego/database"
	"shoego/domain"
	"time"

	"gorm.io/gorm"
)

func SaveOTPFull(data domain.OTPVerification) error {

	err := database.DB.Create(&data).Error
	if err != nil {
		fmt.Println("DB Insert Error:", err)
		return err
	}
	
	return nil

}

func CheckOTPResendAllowed(email string) (bool, error) {

	var otp domain.OTPVerification

	err := database.DB.
		Where("email = ?", email).
		Order("created_at desc").
		First(&otp).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}

		return false, err
	}

	fmt.Println("Last OTP created at:", otp.CreatedAt)
	fmt.Println("Time since:", time.Since(otp.CreatedAt))

	if time.Since(otp.CreatedAt) < 30*time.Second {
		return false, nil
	}

	return true, nil
}

func VerifyOTP(email, otp, otpType string) (domain.OTPVerification, error) {

	var otpData domain.OTPVerification

	err := database.DB.
		Where("email = ? AND otp = ? AND type = ?", email, otp, otpType).
		First(&otpData).Error

	return otpData, err
}
func GetSignupDataByEmail(email string) (domain.OTPVerification, error) {

	var otpData domain.OTPVerification

	err := database.DB.
		Where("email = ?", email).
		Order("created_at desc").
		First(&otpData).Error

	return otpData, err
}

func GetSignupDataFromOTP(email string) (*domain.OTPVerification, error) {
	var otp domain.OTPVerification

	err := database.DB.
		Where("email = ? AND type = ?", email, "signup").
		Order("created_at desc").
		First(&otp).Error

	if err != nil {
		return nil, err
	}

	return &otp, nil
}

func DeleteOTP(email string, otpType string) error {
	return database.DB.Exec("DELETE FROM otp_verifications WHERE email = ? AND type = ?", email, otpType).Error
}

func UpdatePassword(email, newPassword string) (*domain.User, error) {
	var user domain.User
	result := database.DB.Where("email=?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	user.Password = newPassword

	if err := database.DB.Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
