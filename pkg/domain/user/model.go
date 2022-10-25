package user

import (
	"encoding/json"
	"fin-asg/pkg/domain/comment"
	"fin-asg/pkg/domain/photo"
	"fin-asg/pkg/domain/social"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type dob datatypes.Date

var _ json.Unmarshaler = &dob{}

func (mt *dob) UnmarshalJSON(bs []byte) error {
	var s string
	err := json.Unmarshal(bs, &s)
	if err != nil {
		return err
	}
	t, err := time.ParseInLocation("2006-01-02", s, time.UTC)
	if err != nil {
		return err
	}
	*mt = dob(t)
	return nil
}

type User struct {
	ID        uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Username  string    `json:"username" gorm:"column:first_name"`
	Email     string    `json:"email" gorm:"column:email"`
	Password  string    `json:"password"`
	Dob       dob       `json:"dob" gorm:"column:dob"`
	Age       int       `json:"age" gorm:"column:age"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;default:null"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;default:null"`
	Photo     []photo.Photo
	Social    []social.Social
	Comment   []comment.Comment
}
