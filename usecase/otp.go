package usecase

import (
	"errors"
	"fmt"
	"shoego/domain"
	"shoego/helper"
	"shoego/models"
	"shoego/repository"
	"strings"
	"time"
)

func VerifyOTPAndCreateUser(data models.VerifyOTP) (*models.TokenUser, error) {

	data.OTP = strings.TrimSpace(data.OTP)

	// 1. Get OTP record directly
	otpData, err := repository.VerifyOTP(data.Email, data.OTP, "signup")
	if err != nil {
		return nil, err
	}

	// 2. Check if user already exists
	userData, err := repository.FindUserByEmail(data.Email)
	if err != nil {
		return nil, err
	}

	// 3. If user does not exist, create user
	if userData == nil {

		userInsert, err := repository.SignupInsert(models.SignupDetail{
			Name:     otpData.Name,
			Email:    otpData.Email,
			Password: otpData.Password,
			Phone:    otpData.Phone,
		})
		if err != nil {
			return nil, err
		}

		userData = &domain.User{
			ID:    uint(userInsert.ID),
			Name:  userInsert.Name,
			Email: userInsert.Email,
			Phone: userInsert.Phone,
		}
	}

	// 4. Delete OTP AFTER user creation
	_ = repository.DeleteOTP(data.Email, "signup")

	// 5. Prepare response
	userResp := models.SignupDetailResponse{
		ID:    int(userData.ID),
		Name:  userData.Name,
		Email: userData.Email,
		Phone: userData.Phone,
	}

	// 6. Generate tokens
	accessToken, _ := helper.GenerateAccessToken(userResp)
	refreshToken, _ := helper.GenerateRefreshToken(userResp)

	return &models.TokenUser{
		User:         userResp,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func ResendOTP(email string) error {
	allowed, err := repository.CheckOTPResendAllowed(email)
	if err != nil {
		return err
	}

	if !allowed {
		return errors.New("wait before requesting new OTP")
	}

	otp := helper.GenerateOTP()
	expiry := time.Now().Add(2 * time.Minute)

	userData, err := repository.GetSignupDataFromOTP(email) // make sure this fetches full OTP info
	if err != nil {
		return errors.New("cannot find signup info for this email")
	}

	// save new OTP with full info
	userData.OTP = otp
	userData.ExpiresAt = expiry
	err = repository.SaveOTPFull(*userData)
	if err != nil {
		return err
	}

	err = helper.SendOTPEmail(email, otp)
	if err != nil {
		return err
	}

	fmt.Println("Resend OTP:", otp)
	return nil
}

func ForgotPassword(email string) error {
	user, err := repository.FindUserByEmail(email)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("email not found")
	}

	otp := helper.GenerateOTP()

	expiry := time.Now().Add(2 * time.Minute)

	otpData := domain.OTPVerification{
		Email:     email,
		OTP:       otp,
		Type:      "forgot_password",
		ExpiresAt: expiry,
	}
	if err := repository.SaveOTPFull(otpData); err != nil {
		return err
	}
	return helper.SendOTPEmail(email, otp)

}

func ResetPassword(email, otp, newPassword string) (*models.TokenUser, error) {
	// 1️⃣ Verify OTP
	otpData, err := repository.VerifyOTP(email, otp, "forgot_password")
	if err != nil {
		return nil, err
	}

	// Optional: check if OTP is expired
	if otpData.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("OTP expired")
	}
	// 2️⃣ Hash new password
	hashedPass, err := helper.PasswordHashing(newPassword)
	if err != nil {
		return nil, err
	}

	// 3️⃣ Update password in user table
	user, err := repository.UpdatePassword(email, hashedPass)
	if err != nil {
		return nil, err
	}

	// 4️⃣ Delete OTP after verification
	_ = repository.DeleteOTP(email, "forgot_password")

	// 5️⃣ Generate tokens
	userResp := models.SignupDetailResponse{
		ID:    int(user.ID),
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}

	accessToken, _ := helper.GenerateAccessToken(userResp)
	refreshToken, _ := helper.GenerateRefreshToken(userResp)

	return &models.TokenUser{
		User:         userResp,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
