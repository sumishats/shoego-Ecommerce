// package usecase

// import (
// 	"errors"
// 	"shoego/domain"
// 	"shoego/helper"
// 	"shoego/models"
// 	"shoego/repository"

// 	"golang.org/x/crypto/bcrypt"
// )

// func AdminLogin(adminmodel models.AdminLogin) (domain.TokenAdmin, error) {

// 	// getting details from admin with email
// 	AdminDetail, err := repository.AdminLogin(adminmodel)

// 	if err != nil {

// 		return domain.TokenAdmin{}, errors.New("given mail formate have")

// 	}

// 	if AdminDetail.Password == "" {
// 		return domain.TokenAdmin{}, errors.New("error from admin password")
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(AdminDetail.Password), []byte(adminmodel.Password))
// 	if err != nil {

// 		return domain.TokenAdmin{}, models.PasswordIsNotCorrect

// 	}
// 	var AdminDetailsResponse models.AdminDetailsResponse

// 	tokenString, err := helper.GenerateTokenAdmin(AdminDetailsResponse)

// 	if err != nil {

// 		return domain.TokenAdmin{}, errors.New("demo : repository adminlogin")

// 	}

// 	return domain.TokenAdmin{
// 		Admin: AdminDetailsResponse,
// 		Token: tokenString,
// 	}, nil
// }

package usecase

import (
	"errors"
	"shoego/domain"
	"shoego/helper"
	"shoego/models"
	"shoego/repository"

	"golang.org/x/crypto/bcrypt"
)

func AdminLogin(adminmodel models.AdminLogin) (domain.TokenAdmin, error) {

	// getting admin details using email
	AdminDetail, err := repository.AdminLogin(adminmodel)
	if err != nil {
		return domain.TokenAdmin{}, errors.New("admin email not found")
	}

	// check password exists
	if AdminDetail.Password == "" {
		return domain.TokenAdmin{}, errors.New("admin password is empty")
	}

	// compare hashed password with entered password
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
