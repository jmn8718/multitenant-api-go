package models

type SignInUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type SignUpUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type JwtResponse struct {
	AccessToken string `json:"access_token" binding:"required"`
	ExpiresAt   int64  `json:"expires_at" binding:"required"`
}
