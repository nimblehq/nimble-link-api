package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/routes"
)

func SetupRouter(g *gin.Engine) {
	routes.Register(g)
}
