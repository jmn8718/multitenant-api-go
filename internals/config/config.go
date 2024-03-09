package config

import (
	"multitenant-api-go/internals/constants"
	"os"
)

type Config struct {
	Host             string
	Port             string
	Environment      string
	DatabaseUrl      string
	LogLevel         string
	JwtSecret        string
	EnableMigrations bool
	EnableSignup     bool
}

func getEnv(key string, defaultValue string) string {
	println(key, defaultValue)
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	println("default")
	return defaultValue
}

func InitializeConfiguration() Config {
	var conf Config

	conf.Host = getEnv("HOST", "localhost")
	conf.Port = getEnv("PORT", "5000")
	conf.Environment = getEnv("ENV", constants.EnvironmentDevelopment)
	conf.LogLevel = getEnv("LOG_LEVEL", "info")
	conf.DatabaseUrl = getEnv("DB_URL", "")
	conf.JwtSecret = getEnv("JWT_SECRET_KEY", "myjwtsecretkey")

	conf.EnableMigrations = getEnv("ENABLE_MIGRATIONS", "false") == "true"
	conf.EnableSignup = getEnv("ENABLE_SIGNUP", "true") == "true"
	return conf
}
