package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/controllers"
	"github.com/nimble-link/backend/middlewares"
	"github.com/nimble-link/backend/pkg/ginutils"
)

func registerApi(r *ginutils.ApplicationRouter, mids ...gin.HandlerFunc) {
	r.Middlewares(mids...)

	r.Register("GET", "api/v1/links/:alias", controllers.GetLink)
	r.Register("POST", "api/v1/links/:alias", controllers.GetLinkWithPassword)
	r.Register("POST", "api/v1/links", middlewares.Authenticated(), controllers.CreateLink)
	r.Register("GET", "api/v1/links", middlewares.Authenticated(), controllers.GetLinks)
	r.Register("DELETE", "api/v1/links/:id", middlewares.Authenticated(), controllers.DeleteLink)
}
