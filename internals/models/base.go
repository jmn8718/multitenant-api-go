package models

import (
	"time"

	"github.com/lucsky/cuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        string    `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	base.ID = cuid.New()
	return nil
}

type SuccessResponse struct {
	Success bool `json:"success"`
}
