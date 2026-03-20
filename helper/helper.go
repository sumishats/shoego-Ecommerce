package helper

import (
	"errors"
	"shoego/config"
	"shoego/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthCustomClaims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func PasswordHashing(password string) (string, error) {


	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return "", errors.New("hash Server issue")
	}

	hash := string(hashedPassword)

	return hash, nil

}

func GenerateTokenUsers(userID int, userEmail string, expirationTime time.Time) (string, error) {
	cfg, _ := config.LoadConfig()

	claims := &AuthCustomClaims{
		Id:    userID,
		Email: userEmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.KEY))


	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateAccessToken(user models.SignupDetailResponse) (string, error) {

	expirationTime := time.Now().Add(15 * time.Minute) //valid 15 minute 
	tokenString, err := GenerateTokenUsers(user.ID, user.Email, expirationTime) 
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func GenerateRefreshToken(user models.SignupDetailResponse) (string, error) {

	expirationTime := time.Now().Add(24 * 90 * time.Hour) //valid 90 day
	tokeString, err := GenerateTokenUsers(user.ID, user.Email, expirationTime) //use for gew new access token when access token expired
	if err != nil {
		return "", err
	}
	return tokeString, nil

}
