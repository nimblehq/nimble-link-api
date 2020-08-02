package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/services/authentication"
)

func Guest() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := authentication.GetCurrentUserFromContext(c)

		if err == nil || user != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/")
			c.Abort()

			return
		}

		c.Next()
	}
}
