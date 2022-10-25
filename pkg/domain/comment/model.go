package comment

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	UserID    uint64          `json:"user_id"`
	PhotoID   uint64          `json:"photo_id" valid:"required~photo id is required"`
	Message   string          `gorm:"not null" json:"message" valid:"required~comment message is required"`
	CreatedAt time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"column:updated_at;default:null"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;default:null"`
}