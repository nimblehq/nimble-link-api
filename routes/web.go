package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/controllers"
	"github.com/nimble-link/backend/middlewares"
	"github.com/nimble-link/backend/pkg/ginutils"
)

func registerWeb(r *ginutils.ApplicationRouter, mids ...gin.HandlerFunc) {
	r.Middlewares(mids...)

	r.Register("GET", "/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	})

	r.Register("GET", "/dashboard", middlewares.Authenticated(), func(c *gin.Context) {
		c.JSON(http.StatusOK, c.Keys["current_user"])
	})

	r.Register("GET", "/oauth2/callback", controllers.LoginCallback)

	r.Register("GET", "/login", middlewares.Guest(), controllers.Login)
	r.Register("POST", "/logout", middlewares.Authenticated(), controllers.Logout)
}
