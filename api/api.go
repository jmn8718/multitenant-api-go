package api

import (
	"multitenant-api-go/internals/constants"
	"multitenant-api-go/internals/globals"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, logger *zap.SugaredLogger) *gin.Engine {
	if globals.Conf.Environment == constants.EnvironmentDevelopment {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	}))
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	router.GET("/health", HealthCheckHandler)
	return router
}

type HealthCheckResponse struct {
	Status string `json:"status" binding:"required"`
}

func HealthCheckHandler(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, HealthCheckResponse{Status: "ok"})
}
