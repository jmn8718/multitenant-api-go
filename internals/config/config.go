package config

import (
	"multitenant-api-go/internals/constants"
	"os"
)

type Config struct {
	Host        string
	Port        string
	Environment string
	DatabaseUrl string
	LogLevel    string
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}

func InitializeConfiguration() Config {
	var conf Config

	conf.Host = getEnv("HOST", "localhost")
	conf.Port = getEnv("PORT", "5000")
	conf.Environment = getEnv("ENV", constants.EnvironmentDevelopment)
	conf.LogLevel = getEnv("LOG_LEVEL", "info")
	conf.DatabaseUrl = getEnv("DB_URL", "")

	return conf
}
