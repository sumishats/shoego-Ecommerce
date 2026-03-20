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

// @Summary Google Login
// @Description Redirect user to Google authentication page
// @Tags User Authentication
// @Accept json
// @Produce json
// @Success 302 {string} string "Redirect to Google login page"
// @Router /auth/google/login [get]

func GoogleLogin(c *gin.Context) {

	url := config.GoogleOAuthConfig.AuthCodeURL("state")

	c.Redirect(http.StatusTemporaryRedirect, url)
}

// @Summary Google Login Callback
// @Description Google OAuth callback to authenticate user
// @Tags User Authentication
// @Accept json
// @Produce json
// @Param code query string true "Authorization Code"
// @Param state query string true "State"
// @Success 200 {object} response.Response{}
// @Failure 500 {object} response.Response{}
// @Router /auth/google/callback [get]

func GoogleCallback(c *gin.Context) {
	
	code := c.Query("code") // get code from url query params 
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

	//if user not exist  ,create user
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
