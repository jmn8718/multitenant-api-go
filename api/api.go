package api

import (
	api_auth "multitenant-api-go/api/auth"
	api_tenants "multitenant-api-go/api/tenants"
	api_user "multitenant-api-go/api/user"
	"multitenant-api-go/internals/constants"
	"multitenant-api-go/internals/globals"
	"multitenant-api-go/internals/middlewares"

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

	router.Use(middlewares.ErrorHandler(logger))

	router.GET("/health", HealthCheckHandler)

	jwtMiddleware := middlewares.BearerTokenMiddleware(globals.Conf.JwtSecret)

	auth := router.Group("/auth")
	api_auth.RegisterRoutes(auth, db)

	// jwt authenticated routes
	api := router.Group("/api")
	api.Use(jwtMiddleware)
	api_user.RegisterRoutes(api, db)
	api_tenants.RegisterRoutes(api, db)

	return router
}

type HealthCheckResponse struct {
	Status string `json:"status" binding:"required"`
}

func HealthCheckHandler(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, HealthCheckResponse{Status: "ok"})
}
