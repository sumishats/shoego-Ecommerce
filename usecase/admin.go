package usecase

import (
	"errors"
	"shoego/domain"
	"shoego/helper"
	"shoego/models"
	"shoego/repository"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func AdminLogin(adminmodel models.AdminLogin) (domain.TokenAdmin, error) {

	// getting admin details using email
	AdminDetail, err := repository.AdminLogin(adminmodel)
	if err != nil {
		return domain.TokenAdmin{}, errors.New("admin email not found")
	}

	if AdminDetail.Password == "" {
		return domain.TokenAdmin{}, errors.New("admin password is empty")
	}

	err = bcrypt.CompareHashAndPassword([]byte(AdminDetail.Password), []byte(adminmodel.Password))
	if err != nil {
		return domain.TokenAdmin{}, models.PasswordIsNotCorrect
	}

	//  response struct
	adminDetailsResponse := models.AdminDetailsResponse{
		ID:    AdminDetail.ID,
		Name:  AdminDetail.Name,
		Email: AdminDetail.Email,
	}

	tokenString, err := helper.GenerateTokenAdmin(adminDetailsResponse)
	if err != nil {
		return domain.TokenAdmin{}, errors.New("error generating token")
	}

	return domain.TokenAdmin{
		Admin: adminDetailsResponse,
		Token: tokenString,
	}, nil
}

func AdminLogout(authHeader string) error {
	if authHeader == "" {
		return errors.New("authorization header is empty")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
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
