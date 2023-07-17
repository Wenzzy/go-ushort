package models

import (
	"time"
)

type LinkModel struct {
	ID             uint       `gorm:"column:id;primary_key"`
	Name           *string    `gorm:"column:name"`
	RealUrl        string     `gorm:"column:real_url;unique"`
	GeneratedAlias string     `gorm:"column:generated_alias;unique"`
	CreatedAt      *time.Time `gorm:"column:registered_at;not null;default:CURRENT_TIMESTAMP"`
	UserID         uint       `gorm:"column:user_id"`
}

// TableName is Database TableName of this model
func (e *LinkModel) TableName() string {
	return "link"
}
