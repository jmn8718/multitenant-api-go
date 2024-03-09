package middlewares

import (
	"net/http"

	"multitenant-api-go/internals/errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandler(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		for _, err := range c.Errors {
			switch e := err.Err.(type) {
			case errors.HttpError:
				c.AbortWithStatusJSON(e.StatusCode, e)
			default:
				logger.Errorln(err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Service Unavailable"})
			}
		}
	}
}
