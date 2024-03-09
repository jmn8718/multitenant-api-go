package main

import (
	"multitenant-api-go/api"
	"multitenant-api-go/internals/config"
	"multitenant-api-go/internals/database"
	"multitenant-api-go/internals/globals"
	"multitenant-api-go/internals/logging"
	"os"
	"runtime"

	"fmt"

	"go.uber.org/zap"
)

// ConfigRuntime sets the number of operating system threads.
func ConfigRuntime(logger *zap.SugaredLogger) {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	logger.Debugf("Running with %d CPUs\n", nuCPU)
}

//	@title			Multitenant API
//	@version		1.0
//	@description	Multitenant API server.
//	@termsOfService

//	@contact.name	API Support
//	@contact.url
//	@contact.email

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:5000
//	@BasePath	/

//	@securityDefinitions.apikey	JwtAuth
//	@in							header
//	@name						Authorization

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						X-API-Key

// @externalDocs.description
// @externalDocs.url
func main() {
	globals.Conf = config.InitializeConfiguration()

	logger := logging.NewLogger(logging.LoggerConfig{
		Environment: globals.Conf.Environment,
		LogLevel:    globals.Conf.LogLevel,
	})

	ConfigRuntime(logger)

	defer func() {
		if r := recover(); r != nil {
			logger.Errorln("Recovered error", r)
			os.Exit(3)
		}
	}()

	db, dbErr := database.ConnectDatabase(logger)

	if dbErr != nil {
		panic(dbErr)
	}

	router := api.InitRouter(db, logger)

	if err := router.Run(fmt.Sprintf(":%s", globals.Conf.Port)); err != nil {
		logger.Panicw("error: %s", err)
	}
}
