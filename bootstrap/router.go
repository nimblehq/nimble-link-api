package bootstrap

import "github.com/gin-gonic/gin"

func SetupRouter(g *gin.Engine) {
	g.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{})
	})
}
