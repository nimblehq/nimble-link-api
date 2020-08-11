package bootstrap

import (
	"github.com/jinzhu/gorm"
	"github.com/nimble-link/backend/database"
	"github.com/nimble-link/backend/models"
)

func SetupDB() *gorm.DB {
	db := database.InitDB()

	db.LogMode(true)

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Link{})

	return db
}
