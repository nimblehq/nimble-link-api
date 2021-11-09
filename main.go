package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nimble-link/backend/bootstrap"
	"github.com/nimble-link/backend/config"
)

func main() {
	config.InitConfig()

	db := bootstrap.SetupDB()
	defer db.Close()

	g := gin.New()
	bootstrap.SetupGin(g)
	bootstrap.SetupRouter(g)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.AppConfig.Addr), g); err != nil {
		log.Fatal(fmt.Sprintf("Failed to start server %v", err))
	}
}
