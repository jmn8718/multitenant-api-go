package utils_auth

import (
	"multitenant-api-go/internals/constants"
	"multitenant-api-go/internals/models"

	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthClaims struct {
	UserId string `json:"userId"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

type UserClaims struct {
	ID      string
	IsAdmin bool
}

func GenerateJwt(secret string, user UserClaims) (*models.JwtResponse, error) {
	var secretKey = []byte(secret)
	expirationTime := time.Now().Add(60 * time.Minute).Unix()
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime,
	}
	tokenClaims := AuthClaims{
		UserId:         user.ID,
		Role:           constants.USER,
		StandardClaims: *claims,
	}
	if user.IsAdmin {
		tokenClaims.Role = constants.ADMIN
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	accessToken, err := token.SignedString(secretKey)
	if err != nil {
		return nil, err
	}
	return &models.JwtResponse{AccessToken: accessToken, ExpiresAt: expirationTime}, nil
}

func ValidateJwt(secret string, tokenString string) (*AuthClaims, string) {
	claims := &AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(parsedToken *jwt.Token) (interface{}, error) {
		if parsedToken.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("invalid_signing_method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "token is expired") {
				return nil, "token_expired"
			}
			return nil, err.Error()
		}
		return nil, "invalid_token"
	}
	// validate claims
	return token.Claims.(*AuthClaims), ""
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
