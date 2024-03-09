package api_auth

import (
	"multitenant-api-go/internals/errors"
	"multitenant-api-go/internals/globals"
	"multitenant-api-go/internals/models"
	"multitenant-api-go/internals/repositories"
	utils_auth "multitenant-api-go/internals/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SignInHandler godoc
//
//	@Summary	Authenticate a user
//	@Schemes
//	@Description	Authenticates a user using username and password, returns a JWT token if successful
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.SignInUser	true	"User login object"
//	@Success		200		{object}	models.JwtResponse
//	@Failure		400		{object}	errors.HttpError
//	@Failure		401		{object}	errors.HttpError
//	@Failure		500		{object}	errors.HttpError
//	@Router			/auth/signin [post]
func SignInHandler(usersRepository repositories.UsersRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body models.SignInUser
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Error(errors.NewHttpErrorWithReason(400, "bad_request", err.Error(), true))
			return
		}

		newUser, err := usersRepository.SignInUser(body)
		if err != nil {
			c.Error(err)
			return
		}

		data, error := utils_auth.GenerateJwt(globals.Conf.JwtSecret, utils_auth.UserClaims{
			ID: newUser.ID, IsAdmin: newUser.IsAdmin,
		})
		if error != nil {
			c.Error(error)
			return
		}
		c.IndentedJSON(http.StatusOK, data)
	}
}

// SignUpHandler godoc
//
//	@Summary		Register a new user
//	@Schemes		http
//	@Description	Registers a new user with the given username and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.SignUpUser	true	"User registration object"
//	@Success		200		{object}	models.JwtResponse
//	@Failure		400		{object}	errors.HttpError
//	@Failure		500		{object}	errors.HttpError
//	@Router			/auth/signup [post]
func SignUpHandler(usersRepository repositories.UsersRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !globals.Conf.EnableSignup {
			c.Error(errors.NewHttpError(400, "signup_not_enabled"))
			return
		}
		var body models.SignUpUser
		if err := c.ShouldBindJSON(&body); err != nil {
			c.Error(errors.NewHttpErrorWithReason(400, "bad_request", err.Error(), true))
			return
		}
		newUser, err := usersRepository.CreateUser(body)
		if err != nil {
			c.Error(err)
			return
		}

		data, error := utils_auth.GenerateJwt(globals.Conf.JwtSecret, utils_auth.UserClaims{
			ID: newUser.ID, IsAdmin: newUser.IsAdmin,
		})
		if error != nil {
			c.Error(error)
			return
		}
		c.IndentedJSON(http.StatusOK, data)
	}
}

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) {
	var usersRepository = repositories.NewUsersRepository(db)
	router.POST("/signup", SignUpHandler(usersRepository))
	router.POST("/signin", SignInHandler(usersRepository))
}
