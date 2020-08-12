package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/middlewares"
	"github.com/nimble-link/backend/pkg/ginutils"
)

func Register(g *gin.Engine) *gin.Engine {
	g.Use(gin.Recovery())
	g.Use(gin.Logger())

	store := cookie.NewStore([]byte("secret"))
	g.Use(sessions.Sessions("backend", store))

	g.Use(middlewares.CurrentUser())

	r := &ginutils.ApplicationRouter{Router: g}

	registerWeb(r)

	return g
}
