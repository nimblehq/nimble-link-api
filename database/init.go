package database

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nimble-link/backend/config"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	db, err := gorm.Open(config.DBConfig.Connection, config.DBConfig.URL)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to database %v", err))
	}

	DB = db

	return db
}
