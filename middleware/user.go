package middleware

import (
	"fmt"
	"net/http"
	"shoego/helper"
	"shoego/repository"
	"shoego/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Request Headerss:", c.Request.Header)

		
		authheader := c.GetHeader("Authorization")

		tokenString := helper.GetTokenFromHeader(authheader)

		//validate token and extract user id
		if tokenString == "" {
			var err error
			tokenString, err = c.Cookie("Authorization")
			if err != nil {

				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}

		// check user is logout or not 
		blacklisted, err := repository.IsTokenBlacklist(tokenString)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if blacklisted {
			errRes := response.ClientResponse(http.StatusUnauthorized, "token is invalid or already logged out", nil, "blacklisted token")
			c.JSON(http.StatusUnauthorized, errRes)
			c.Abort()
			return
		}

		userId, userEmail, err := helper.ExtractUserIDFromToken(tokenString)
		fmt.Println("userId", userId, "userEmail", userEmail)

		if err != nil {

			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		
		c.Set("user_id", uint(userId))
		c.Set("user_email", userEmail)

		c.Next()
	}
}
