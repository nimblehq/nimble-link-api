package config

import (
	"os"
)

type appConfig struct {
	Name      string
	RunMode   string
	Addr      string
	SecretKey string

	PublicPath string
	URL        string
}

func newAppConfig() *appConfig {
	return &appConfig{
		Name:      "backend",
		RunMode:   os.Getenv("APP_RUN_MODE"),
		Addr:      os.Getenv("APP_ADDR"),
		SecretKey: os.Getenv("APP_SECRET_KEY"),
	}
}
