package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/services/authentication"
)

func CurrentUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		authentication.SaveCurrentUserToContext(c)

		c.Next()
	}
}
