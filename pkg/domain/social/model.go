package social

import (
	"time"

	"gorm.io/gorm"
)

type Social struct {
	ID        uint64          `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name      string          `gorm:"not null" json:"name" valid:"required~social media name is required"`
	URL       string          `gorm:"not null" json:"url" valid:"required~social media url is required,url~invalid url format"`
	UserID    uint64          `json:"user_id"`
	CreatedAt time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"column:updated_at;default:null"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;default:null"`
}
