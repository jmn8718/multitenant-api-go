package api_tenants

import (
	"multitenant-api-go/internals/constants"
	"multitenant-api-go/internals/errors"
	"multitenant-api-go/internals/middlewares"
	"multitenant-api-go/internals/models"
	"multitenant-api-go/internals/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// getMyTenantsHandler godoc
//
//	@Summary	Get tenants
//	@Schemes
//	@Description	Get user's tenants
//	@Tags			tenant
//	@Security		JwtAuth
//	@Produce		json
//	@Success		200	{object}	models.MyTenantsDataResponse
//	@Failure		401	{object}	errors.HttpError	"Unauthorized"
//	@Failure		500	{object}	errors.HttpError	"Internal error"
//	@Router			/api/tenants [get]
func getMyTenantsHandler(tenantsRepository repositories.TenantsRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, _ := c.Get(constants.UserId)
		data, err := tenantsRepository.GetMyTenants(userId.(string))

		if err != nil {
			c.Error(err)
			return
		}
		c.IndentedJSON(http.StatusCreated, data)
	}
}

// createUserTenantHandler godoc
//
//	@Summary	Create tenant
//	@Schemes
//	@Description	Create tenant for the authenticated user
//	@Tags			tenant
//	@Security		JwtAuth
//	@Produce		json
//	@Param			tenant	body		models.TenantCreateRequest	true	"Tenant data"
//	@Success		201		{object}	models.TenantCreateResponse
//	@Failure		400		{object}	errors.HttpError	"Bad request"
//	@Failure		401		{object}	errors.HttpError	"Unauthorized"
//	@Failure		500		{object}	errors.HttpError	"Internal error"
//	@Router			/api/tenants [post]
func createUserTenantHandler(tenantsRepository repositories.TenantsRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body models.TenantCreateRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Error(errors.NewHttpErrorWithReason(400, "bad_request", err.Error(), true))
			return
		}
		userId, _ := c.Get(constants.UserId)
		tenant, err := tenantsRepository.CreateTenantWithUser(body, userId.(string))

		if err != nil {
			c.Error(err)
			return
		}
		c.IndentedJSON(http.StatusCreated, tenant)
	}
}

// getMyTenantUsersHandler godoc
//
//	@Summary	Get tenant's users
//	@Schemes
//	@Description	Get tenant's users
//	@Tags			users, tenant
//	@Security		JwtAuth
//	@Produce		json
//	@Param			tenantId	path		string	true	"tenant id"
//	@Success		200			{object}	models.TenantUsersResponse
//	@Failure		401			{object}	errors.HttpError	"Unauthorized"
//	@Failure		500			{object}	errors.HttpError	"Internal error"
//	@Router			/api/tenants/:tenantId/users [get]
func getMyTenantUsersHandler(tenantsRepository repositories.TenantsRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantId, _ := c.Get(constants.TenantId)
		userId, _ := c.Get(constants.UserId)

		data, err := tenantsRepository.GetMyTenantUsers(userId.(string), tenantId.(string))

		if err != nil {
			c.Error(err)
			return
		}
		c.IndentedJSON(http.StatusCreated, data)
	}
}

// postMyTenantAddUserHandler godoc
//
//	@Summary	Add user to tenant
//	@Schemes
//	@Description	Add new user to tenant
//	@Tags			users, tenant
//	@Security		JwtAuth
//	@Produce		json
//	@Param			tenantId	path		string						true	"tenant id"
//	@Param			user		body		models.TenantAddUserRequest	true	"User data"
//	@Success		200			{object}	models.SuccessResponse
//	@Failure		400			{object}	errors.HttpError	"Invalid request"
//	@Failure		401			{object}	errors.HttpError	"Unauthorized"
//	@Failure		403			{object}	errors.HttpError	"Forbidden"
//	@Failure		500			{object}	errors.HttpError	"Internal error"
//	@Router			/api/tenants/:tenantId/users [post]
func postMyTenantAddUserHandler(tenantsRepository repositories.TenantsRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body models.TenantAddUserRequest
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Error(errors.NewHttpErrorWithReason(400, "bad_request", err.Error(), true))
			return
		}

		tenantId, _ := c.Get(constants.TenantId)
		userId, _ := c.Get(constants.UserId)
		err := tenantsRepository.AddUserToMyTenant(userId.(string), tenantId.(string), body.Email)

		if err != nil {
			c.Error(err)
			return
		}
		c.IndentedJSON(http.StatusCreated, models.SuccessResponse{Success: true})
	}
}

// getMyTenantUsersHandler godoc
//
//	@Summary	Get api key
//	@Schemes
//	@Description	Get tenant's api key
//	@Tags			tenant
//	@Security		JwtAuth
//	@Produce		json
//	@Param			tenantId	path		string	true	"tenant id"
//	@Success		200			{object}	models.TenantApiKeyResponse
//	@Failure		401			{object}	errors.HttpError	"Unauthorized"
//	@Failure		500			{object}	errors.HttpError	"Internal error"
//	@Router			/api/tenants/:tenantId/keys [get]
func getMyTenantApiKeyHandler(tenantsRepository repositories.TenantsRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantId, _ := c.Get(constants.TenantId)
		data, err := tenantsRepository.GetMyTenantApiKey(tenantId.(string))
		if err != nil {
			c.Error(err)
			return
		}
		c.IndentedJSON(http.StatusOK, data)
	}
}

// updateMyTenantApiKeyHandler godoc
//
//	@Summary	Update api key
//	@Schemes
//	@Description	Update (replace) tenant's api key
//	@Tags			tenant
//	@Security		JwtAuth
//	@Produce		json
//	@Param			tenantId	path		string	true	"tenant id"
//	@Success		200			{object}	models.TenantApiKeyResponse
//	@Failure		401			{object}	errors.HttpError	"Unauthorized"
//	@Failure		500			{object}	errors.HttpError	"Internal error"
//	@Router			/api/tenants/:tenantId/keys [patch]
func updateMyTenantApiKeyHandler(tenantsRepository repositories.TenantsRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantId, _ := c.Get(constants.TenantId)
		data, err := tenantsRepository.UpdateMyTenantApiKey(tenantId.(string))

		if err != nil {
			c.Error(err)
			return
		}
		c.IndentedJSON(http.StatusOK, data)
	}
}

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) {
	var tenantsRepository = repositories.NewTenantsRepository(db)
	var userParamTenantAccessMiddleware = middlewares.UserParamTenantAccessMiddleware(tenantsRepository, "")
	var userAdminParamTenantAccessMiddleware = middlewares.UserParamTenantAccessMiddleware(tenantsRepository, constants.ADMIN)
	router.GET("/tenants", getMyTenantsHandler(tenantsRepository))
	router.POST("/tenants", createUserTenantHandler(tenantsRepository))
	router.GET("/tenants/:tenantId/users", userParamTenantAccessMiddleware, getMyTenantUsersHandler(tenantsRepository))
	router.POST("/tenants/:tenantId/users", userAdminParamTenantAccessMiddleware, postMyTenantAddUserHandler(tenantsRepository))
	router.GET("/tenants/:tenantId/keys", userParamTenantAccessMiddleware, getMyTenantApiKeyHandler(tenantsRepository))
	router.PATCH("/tenants/:tenantId/keys", userAdminParamTenantAccessMiddleware, updateMyTenantApiKeyHandler(tenantsRepository))
}
