package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/services/authentication"
)

func Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, err := authentication.GetCurrentUserFromContext(c)
		if err != nil || currentUser == nil {
			c.AbortWithStatus(http.StatusUnauthorized)

			return
		}

		c.Next()
	}
}
