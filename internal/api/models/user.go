package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null;type:varchar(255)"`
	Name      string    `json:"name" gorm:"type:varchar(255)"`
	GoogleID  string    `json:"google_id" gorm:"uniqueIndex;type:varchar(255)"`
	Picture   string    `json:"picture" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}
