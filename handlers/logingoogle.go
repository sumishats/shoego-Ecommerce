package handlers

import (
	"net/http"
	"shoego/config"
	"shoego/domain"
	"shoego/helper"
	"shoego/models"
	"shoego/repository"
	"shoego/usecase"

	"github.com/gin-gonic/gin"
)

func GoogleLogin(c *gin.Context) {

	url := config.GoogleOAuthConfig.AuthCodeURL("state")

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code not found"})
		return
	}

	googleUser, err := usecase.GetGoogleUser(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := repository.FindUserByEmail(googleUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var userData *domain.User

	if existingUser == nil {
		newUser := domain.User{
			Name:     googleUser.Name,
			Email:    googleUser.Email,
			Phone:    "",
			Password: "",
			Blocked:  false,
			IsAdmin:  false,
		}

		createdUser, err := repository.CreateGoogleUser(newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		userData = createdUser
	} else {
		userData = existingUser
	}

	userResp := models.SignupDetailResponse{
		ID:    int(userData.ID),
		Name:  userData.Name,
		Email: userData.Email,
		Phone: userData.Phone,
	}

	accessToken, err := helper.GenerateAccessToken(userResp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, err := helper.GenerateRefreshToken(userResp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, models.TokenUser{
		User:         userResp,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
