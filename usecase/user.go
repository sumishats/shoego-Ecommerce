package usecase

import (
	"errors"
	"log"
	"net/mail"
	"shoego/domain"
	"shoego/helper"
	"shoego/models"
	"shoego/repository"
	"strings"
	"time"

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

	hashedPassword, err := helper.PasswordHashing(user.Password)
	if err != nil {
		return err
	}

	otp := helper.GenerateOTP()
	expiry := time.Now().Add(2 * time.Minute)

	// temporary save full info in OTP table before signup
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

		log.Println("OTP save error", err)
		return err
	}
	err = helper.SendOTPEmail(user.Email, otp)
	if err != nil {
		return err
	}

	log.Println("OTP sent", otp)
	return nil
}

func UserLogged(user models.UserLogin) (*models.TokenUser, error) {

	// Validate email format
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return &models.TokenUser{}, errors.New("EMAIL SHOULD BE CORRECT FORMAT")
	}

	// Get user from users table
	userDetails, err := repository.FindUserDetailByEmail(user)
	if err != nil {
		return &models.TokenUser{}, err
	}
	if userDetails.ID == 0 {
		return &models.TokenUser{}, models.ErrEmailNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(user.Password))
	if err != nil {
		log.Println("password mismatch")

		return &models.TokenUser{}, errors.New("hashed password not matching")
	}

	//create user response
	userResp := models.SignupDetailResponse{
		ID:    int(userDetails.ID),
		Name:  userDetails.Name,
		Email: userDetails.Email,
		Phone: userDetails.Phone,
	}

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

func GetUserProfile(userID uint) (*models.UserProfileResponse, error) {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	addresses, err := repository.GetAddressesByUserID(userID)
	if err != nil {
		return nil, err
	}

	//address convert to response
	var addressResp []models.AddressResponse
	for _, addr := range addresses {
		addressResp = append(addressResp, models.AddressResponse{
			ID:        addr.ID,
			Name:      addr.Name,
			Phone:     addr.Phone,
			HouseName: addr.HouseName,
			Street:    addr.Street,
			City:      addr.City,
			State:     addr.State,
			Pincode:   addr.Pincode,
			IsDefault: addr.IsDefault,
		})
	}

	return &models.UserProfileResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Phone:        user.Phone,
		ProfileImage: user.ProfileImage,
		Addresses:    addressResp,
	}, nil
}

func EditUserProfile(userID uint, req models.EditProfileRequest) error {
	return repository.UpdateUserProfile(userID, req.Name, req.Phone, req.ProfileImage)
}

func ChangePassword(userID uint, req models.ChangePasswordRequest) error {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return errors.New("old password is incorrect")
	}

	hashedPassword, err := helper.PasswordHashing(req.NewPassword)
	if err != nil {
		return err
	}

	return repository.UpdateUserPassword(userID, hashedPassword)
}

func RequestEmailChange(userID uint, newEmail string) error {

	otp := helper.GenerateOTP()

	err := repository.SaveOTPFull(domain.OTPVerification{
		Email: newEmail,
		OTP:   otp,
		Type:  "email_change",
	})
	if err != nil {
		return err
	}

	return helper.SendOTPEmail(newEmail, otp)
}

// verify email , if valid update email in db
func VerifyEmailChange(userID uint, req models.VerifyEmailChangeRequest) error {
	_, err := repository.VerifyOTP(req.NewEmail, req.OTP, "email_change")
	if err != nil {
		return err
	}

	err = repository.UpdateUserEmail(userID, req.NewEmail)
	if err != nil {
		return err
	}

	_ = repository.DeleteOTP(req.NewEmail, "email_change")
	return nil
}

func AddUserAddress(userID uint, req models.AddAddressRequest) error {
	//check new address is default or not
	if req.IsDefault {
		_ = repository.ClearDefaultAddresses(userID)
	}

	address := domain.Address{
		UserID:    userID,
		Name:      req.Name,
		Phone:     req.Phone,
		HouseName: req.HouseName,
		Street:    req.Street,
		City:      req.City,
		State:     req.State,
		Pincode:   req.Pincode,
		IsDefault: req.IsDefault,
	}

	return repository.AddAddress(address)
}

func EditUserAddress(userID uint, addressID uint, req models.EditAddressRequest) error {
	address, err := repository.GetAddressByID(addressID)
	if err != nil {
		return err
	}

	if address.UserID != userID {
		return errors.New("address does not belong to this user")
	}

	if req.IsDefault {
		_ = repository.ClearDefaultAddresses(userID)
	}

	return repository.UpdateAddress(addressID, map[string]interface{}{
		"name":       req.Name,
		"phone":      req.Phone,
		"house_name": req.HouseName,
		"street":     req.Street,
		"city":       req.City,
		"state":      req.State,
		"pincode":    req.Pincode,
		"is_default": req.IsDefault,
	})
}

func DeleteUserAddress(userID uint, addressID uint) error {
	address, err := repository.GetAddressByID(addressID)
	if err != nil {
		return err
	}

	if address.UserID != userID {
		return errors.New("address does not belong to this user")
	}

	return repository.DeleteAddress(addressID)
}

func Logout(authHeader string) error {
	if authHeader == "" {
		return errors.New("authorization header is empty")
	}

	if !strings.HasPrefix(authHeader, "Bearer") {
		return errors.New("invalid authorization format")
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return errors.New("token is empty")
	}
	err := repository.SaveBlacklistToken(token)
	if err != nil {
		return err
	}
	return nil
}
