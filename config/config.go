package config

import (
	"os"

	"github.com/joho/godotenv"
)

const (
	RunModeDebug   = "debug"
	RunModeRelease = "release"
	RunModeTest    = "test"
)

var (
	AppConfig *appConfig
	DBConfig  *dbConfig
)

func InitConfig() {
	loadEnv()

	AppConfig = newAppConfig()
	DBConfig = newDBConfig()
}

func loadEnv() {
	env := os.Getenv("APP_ENV")
	if "" == env {
		env = RunModeDebug
	}

	godotenv.Load(".env." + env + ".local")
	if "test" != env {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load()
}
