package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/controllers"
	"github.com/nimble-link/backend/middlewares"
	"github.com/nimble-link/backend/pkg/ginutils"
)

func registerApi(r *ginutils.ApplicationRouter, mids ...gin.HandlerFunc) {
	r.Middlewares(mids...)

	r.Register("POST", "/links", middlewares.Authenticated(), controllers.CreateLink)
	r.Register("GET", "/links", middlewares.Authenticated(), controllers.GetLinks)
	r.Register("DELETE", "/links/:id", middlewares.Authenticated(), controllers.DeleteLink)
}
