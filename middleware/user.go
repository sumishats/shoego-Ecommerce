package middleware

import (
	"net/http"
	"shoego/helper"
	"shoego/repository"
	"shoego/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := helper.GetTokenFromHeader(authHeader)

		// if header token missing, check cookie
		if tokenString == "" {
			cookieToken, err := c.Cookie("Authorization")
			if err != nil {
				errRes := response.ClientResponse(
					http.StatusUnauthorized,
					"authorization token missing",
					nil,
					"token not found in header or cookie",
				)
				c.JSON(http.StatusUnauthorized, errRes)
				c.Abort()
				return
			}
			tokenString = cookieToken
		}

		// check token is blacklisted or not
		blacklisted, err := repository.IsTokenBlacklist(tokenString)
		if err != nil {
			errRes := response.ClientResponse(
				http.StatusInternalServerError,
				"failed to validate token",
				nil,
				err.Error(),
			)
			c.JSON(http.StatusInternalServerError, errRes)
			c.Abort()
			return
		}

		if blacklisted {
			errRes := response.ClientResponse(
				http.StatusUnauthorized,
				"token is invalid or already logged out",
				nil,
				"blacklisted token",
			)
			c.JSON(http.StatusUnauthorized, errRes)
			c.Abort()
			return
		}

		userID, userEmail, err := helper.ExtractUserIDFromToken(tokenString)
		if err != nil {
			errRes := response.ClientResponse(http.StatusUnauthorized, "invalid token", nil, err.Error())
			c.JSON(http.StatusUnauthorized, errRes)
			c.Abort()
			return
		}
		
		c.Set("user_id", userID)
		c.Set("user_email", userEmail)

		

		c.Next()
	}
}
