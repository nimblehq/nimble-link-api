package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/config"
)

func SetupGin(g *gin.Engine) {
	gin.SetMode(config.AppConfig.RunMode)
}
