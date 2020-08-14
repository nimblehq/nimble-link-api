package routes

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/middlewares"
	"github.com/nimble-link/backend/pkg/ginutils"
)

func Register(g *gin.Engine) *gin.Engine {
	g.Use(gin.Recovery())
	g.Use(gin.Logger())

	config := cors.DefaultConfig()
	config.AddAllowHeaders("X-Requested-With")
	config.AddAllowHeaders("Authorization")
	config.AllowOriginFunc = func(origin string) bool {
		return strings.HasPrefix(origin, "http://localhost") || strings.HasSuffix(origin, "nimble-link.vercel.app")
	}
	g.Use(cors.New(config))

	g.Use(middlewares.CurrentUser())

	r := &ginutils.ApplicationRouter{Router: g}

	registerWeb(r)
	registerApi(r)

	return g
}
