package middlewares

import (
	"multitenant-api-go/internals/constants"
	utils_auth "multitenant-api-go/internals/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func BearerTokenMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		if bearerToken == "" {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "code": "missing_token"})
			c.Abort()
			return
		}
		tokenParts := strings.Split(bearerToken, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "code": "invalid_authentication_token"})
			c.Abort()
			return
		}
		claims, reason := utils_auth.ValidateJwt(secret, tokenParts[1])
		if reason != "" {
			c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "code": reason})
			c.Abort()
			return
		}
		c.Set(constants.UserId, claims.UserId)
		c.Next()
	}
}
