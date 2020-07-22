package bootstrap

import (
	"github.com/jinzhu/gorm"
	"github.com/nimble-link/backend/database"
)

func SetupDB() *gorm.DB {
	db := database.InitDB()

	// Migration

	return db
}
