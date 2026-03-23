package middleware

import (
	"net/http"
	"shoego/helper"
	"shoego/repository"
	"shoego/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenheader := c.GetHeader("Authorization")

		if tokenheader == "" {
			response := response.ClientResponse(http.StatusUnauthorized, "no auth header proviededs", nil, nil)
			c.JSON(http.StatusUnauthorized, response)

			c.Abort()
			return
		}

		splitted := strings.Split(tokenheader, " ")

		if len(splitted) != 2 {
			response := response.ClientResponse(http.StatusUnauthorized, "invalid token format", nil, nil)
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return

		}

		tokenpart := splitted[1]

		blacklisted, err := repository.IsTokenBlacklist(tokenpart)
		if err != nil {
			response := response.ClientResponse(http.StatusInternalServerError, "error checking blacklist token", nil, err.Error())
			c.JSON(http.StatusInternalServerError, response)
			c.Abort()
			return
		}

		if blacklisted {
			response := response.ClientResponse(http.StatusUnauthorized, "token is invalid or already logged out", nil, "blacklisted token")
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}
		_, err = helper.ValidateToken(tokenpart)
		if err != nil {
			response := response.ClientResponse(http.StatusUnauthorized, "Invalid token", nil, err.Error())
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		c.Next()

	}
}
