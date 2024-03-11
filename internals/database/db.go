package database

import (
	"errors"
	"time"

	"multitenant-api-go/internals/config"
	"multitenant-api-go/internals/globals"
	"multitenant-api-go/internals/models"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(logger *zap.SugaredLogger, config config.Config) (*gorm.DB, error) {
	var database *gorm.DB
	var err error
	var connected = false
	for i := 1; i <= 3; i++ {
		database, err = gorm.Open(postgres.Open(globals.Conf.DatabaseUrl), &gorm.Config{})
		if err == nil {
			connected = true
			break
		} else {
			logger.Warnln("Attempt %d: Failed to initialize database. Retrying...", i)
			time.Sleep(3 * time.Second)
		}
	}

	if database == nil {
		logger.Panicln("DB not initialized")
		err = errors.New("DB not initialized")
	} else if !connected {
		err = errors.New("DB not connected")
	} else {
		logger.Debugln("Connected to db")
		if config.EnableMigrations {
			logger.Debugln("Running migrations...")
			database.AutoMigrate(&models.User{})
			database.AutoMigrate(&models.Tenant{})
			database.AutoMigrate(&models.UserTenants{})
			logger.Debugln("Migrations completed")
		}
	}

	return database, err
}
