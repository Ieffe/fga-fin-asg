package photo

import (
	"fin-asg/pkg/domain/comment"
	"time"

	"gorm.io/gorm"
)

type Photo struct {
	ID        uint64            `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string            `gorm:"not null" json:"title" valid:"required~photo title is required"`
	Caption   string            `json:"caption"`
	Url       string            `gorm:"not null" json:"url" valid:"required~photo url is required,url~invalid url format"`
	UserID    uint64            `json:"user_id"`
	CreatedAt time.Time         `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time         `json:"updated_at" gorm:"column:updated_at;default:null"`
	DeletedAt *gorm.DeletedAt   `json:"deleted_at" gorm:"column:deleted_at;default:null"`
	Comments  []comment.Comment `json:"comments"`
}
