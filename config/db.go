package config

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/xo/dburl"
)

type dbConfig struct {
	Connection string

	URL   string
	Debug bool
}

func newDBConfig() *dbConfig {
	connection := os.Getenv("DB_CONNECTION")

	var databaseURL string

	if gin.Mode() == gin.ReleaseMode {
		url := os.Getenv("DATABASE_URL")
		u, err := dburl.Parse(url)
		if err != nil {
			log.Fatal(err)
		}

		databaseURL = u.DSN
		connection = u.Driver
	} else {
		username := os.Getenv("DB_USERNAME")
		password := os.Getenv("DB_PASSWORD")
		host := os.Getenv("DB_HOST")
		port := os.Getenv("DB_PORT")

		database := os.Getenv("DB_DATABASE")
		database = database + "_" + AppConfig.RunMode

		databaseURL = createDatabaseURL(username, password, host, port, database)
	}

	return &dbConfig{
		Connection: connection,
		URL:        databaseURL,
	}
}

func createDatabaseURL(username string, password string, host string, port string, database string) string {
	return fmt.Sprintf("sslmode=disable host=%s port=%s dbname=%s user=%s password=%s",
		host,
		port,
		database,
		username,
		password,
	)
}
