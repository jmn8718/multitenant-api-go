package api_user

import (
	"multitenant-api-go/internals/constants"
	"multitenant-api-go/internals/models"
	"multitenant-api-go/internals/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// getMeHandler godoc
//
//	@Summary	Get Authenticate a user
//	@Schemes
//	@Description	Get authenticated user information
//	@Tags			user
//	@Security		JwtAuth
//	@Produce		json
//	@Success		200	{object}	models.UserMeResponse
//	@Failure		401	{object}	errors.HttpError "Unauthorized"
//	@Failure		500	{object}	errors.HttpError "Internal error"
//	@Router			/api/me [get]
func getMeHandler(usersRepository repositories.UsersRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		UserId, _ := c.Get(constants.UserId)

		dbUser, err := usersRepository.GetUser(UserId.(string))
		if err != nil {
			c.Error(err)
			return
		}

		c.IndentedJSON(http.StatusOK, models.UserMeResponse{Id: dbUser.ID, Email: dbUser.Email, Name: dbUser.Name, Image: dbUser.Image, EmailVerified: dbUser.EmailVerified})
	}
}

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) {
	var usersRepository = repositories.NewUsersRepository(db)
	router.GET("/me", getMeHandler(usersRepository))
}
