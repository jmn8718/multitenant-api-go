package models

import "time"

type User struct {
	BaseModel
	Email         string    `json:"email" gorm:"unique,required"`
	Password      string    `json:"password"`
	Name          string    `json:"name"`
	EmailVerified bool      `json:"emailVerified" gorm:"default:false"`
	VerifiedAt    time.Time `json:"verifiedAt"`
	Image         string    `json:"image"`
	IsAdmin       bool      `json:"-"  gorm:"default:false"`
}

type UserMeResponse struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Image         string `json:"image"`
	EmailVerified bool   `json:"emailVerified"`
}
