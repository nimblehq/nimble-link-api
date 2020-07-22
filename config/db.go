package config

import (
	"fmt"
	"os"
)

type dbConfig struct {
	Connection string
	Host       string
	Port       string
	Database   string
	Username   string
	Password   string

	URL   string
	Debug bool
}

func newDBConfig() *dbConfig {
	connection := os.Getenv("DB_CONNECTION")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	database := os.Getenv("DB_DATABASE")
	database = database + "_" + AppConfig.RunMode

	url := createDatabaseURL(username, password, host, port, database)

	return &dbConfig{
		Connection: connection,
		Host:       host,
		Port:       port,
		Database:   database,
		Username:   username,
		Password:   password,
		URL:        url,
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
