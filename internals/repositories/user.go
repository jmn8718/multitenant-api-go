package repositories

import (
	"multitenant-api-go/internals/errors"
	"multitenant-api-go/internals/models"
	utils_auth "multitenant-api-go/internals/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UsersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) UsersRepository {
	return UsersRepository{db: db}
}

func (usersRepository UsersRepository) GetUser(userId string) (*models.User, error) {
	var dbUser models.User

	if err := usersRepository.db.Where("id = ?", userId).First(&dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewHttpError(404, "not_found")
		} else {
			return nil, errors.NewHttpErrorWithReason(500, "internal_error", err.Error(), false)
		}
	}
	return &dbUser, nil
}

func (usersRepository UsersRepository) CreateUser(payload models.SignUpUser) (*models.User, error) {
	var count int64
	if err := usersRepository.db.Model(&models.User{}).Where("email = ?", payload.Email).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		// once there is verification flow, return success
		return nil, errors.NewHttpErrorWithReason(500, "internal_error", "email_exists", false)
	}

	hashedPassword, err := utils_auth.HashPassword(payload.Password)
	if err != nil {
		return nil, errors.NewHttpErrorWithReason(500, "internal_error", err.Error(), false)
	}

	newUser := models.User{Email: payload.Email, Password: hashedPassword, Name: payload.Email, EmailVerified: true, VerifiedAt: time.Now()}

	if err := usersRepository.db.Create(&newUser).Error; err != nil {
		return nil, errors.NewHttpErrorWithReason(500, "internal_error", err.Error(), false)
	}
	return &newUser, nil
}

func (usersRepository UsersRepository) SignInUser(payload models.SignInUser) (*models.User, error) {
	var dbUser models.User

	if err := usersRepository.db.Where("email = ?", payload.Email).Select("id, password, is_admin").First(&dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NewHttpError(401, "invalid_credentials")
		} else {
			return nil, errors.NewHttpErrorWithReason(500, "internal_error", err.Error(), false)
		}
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(payload.Password)); err != nil {
		return nil, errors.NewHttpError(401, "invalid_credentials")
	}
	return &dbUser, nil
}
