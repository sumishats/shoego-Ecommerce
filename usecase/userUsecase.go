package usecase

import (
	"errors"
	"fmt"
	"net/mail"
	"shoego/domain"
	"shoego/helper"
	"shoego/models"
	"shoego/repository"
	"time"

	"github.com/jinzhu/copier"
	"golang.org/x/crypto/bcrypt"
)

func UsersignUp(user models.SignupDetail) error {
	// check email already exists
	existingUser, err := repository.CheckingEmailValidation(user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email already exists")
	}

	// hash password
	hashedPassword, err := helper.PasswordHashing(user.Password)
	if err != nil {
		return err
	}

	// generate OTP
	otp := helper.GenerateOTP()
	expiry := time.Now().Add(2 * time.Minute)

	// save full info in OTP table
	otpData := domain.OTPVerification{
		Email:     user.Email,
		Name:      user.Name,
		Phone:     user.Phone,
		Password:  hashedPassword,
		OTP:       otp,
		Type:      "signup",
		ExpiresAt: expiry,
	}

	err = repository.SaveOTPFull(otpData)
	if err != nil {
		fmt.Println("OTP save error:", err)
		return err
	}
	err = helper.SendOTPEmail(user.Email, otp)
	if err != nil {
		return err
	}

	fmt.Println("OTP sent:", otp)
	return nil
}

func UserLogged(user models.UserLogin) (*models.TokenUser, error) {

	// 1. Validate email format
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return &models.TokenUser{}, errors.New("EMAIL SHOULD BE CORRECT FORMAT")
	}

	// 2. Get user from users table
	userDetails, err := repository.FindUserDetailByEmail(user)
	if err != nil {
		return &models.TokenUser{}, err
	}
	if userDetails.ID == 0 {
		return &models.TokenUser{}, models.ErrEmailNotFound
	}
	fmt.Println("Password from DB (hash):", userDetails.Password)
	fmt.Println("Password from input:", user.Password)

	// 3. Compare input password with hashed password in users table
	err = bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(user.Password))
	if err != nil {
		fmt.Println("Password mismatch")
		return &models.TokenUser{}, errors.New("hashed password not matching")
	}

	// 4. Copy to response struct
	var userResp models.SignupDetailResponse
	err = copier.Copy(&userResp, &userDetails)
	if err != nil {
		return &models.TokenUser{}, errors.New("error in copier")
	}

	// 5. Generate tokens
	accessToken, err := helper.GenerateAccessToken(userResp)
	if err != nil {
		return &models.TokenUser{}, errors.New("could not create access token")
	}
	refreshToken, err := helper.GenerateRefreshToken(userResp)
	if err != nil {
		return &models.TokenUser{}, errors.New("could not create refresh token")
	}

	return &models.TokenUser{
		User:         userResp,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
