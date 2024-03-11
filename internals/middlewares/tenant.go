package middlewares

import (
	"multitenant-api-go/internals/constants"
	"multitenant-api-go/internals/errors"
	"multitenant-api-go/internals/repositories"

	"github.com/gin-gonic/gin"
)

func UserParamTenantAccessMiddleware(tenantsRepository repositories.TenantsRepository, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get(constants.UserId)
		if exists {
			tenantId := c.Param("tenantId")
			if tenantId != "" {
				data, err := tenantsRepository.HasUserAccessToTenant(userId.(string), tenantId)
				if err != nil {
					c.Error(err)
					c.Abort()
					return
				}
				if data != nil && data.Role != "" {
					if role != "" && data.Role != role {
						c.Error(errors.NewHttpError(403, "forbidden"))
						c.Abort()
						return
					}
					c.Set(constants.TenantId, tenantId)
					c.Set(constants.UserRole, data.Role)
					return
				}
			}
		}
		c.Error(errors.NewHttpError(401, "unauthorized"))
		c.Abort()
		return
	}
}
