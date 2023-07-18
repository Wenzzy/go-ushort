package models

import (
	"go-ushorter/app/common/constants/emsgs"
	"go-ushorter/app/common/database"
	"go-ushorter/app/common/utils"
	"time"
)

type LinkModel struct {
	ID             uint       `gorm:"column:id;primary_key"`
	Name           *string    `gorm:"column:name"`
	RealUrl        string     `gorm:"column:real_url;unique"`
	GeneratedAlias string     `gorm:"column:generated_alias;unique"`
	CreatedAt      *time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
	UserID         uint       `gorm:"column:user_id"`
}

// TableName is Database TableName of this model
func (e *LinkModel) TableName() string {
	return "link"
}

func FindAllLinks(condition any) ([]LinkModel, *utils.CommonError) {
	db := database.GetDB()
	var models []LinkModel
	err := db.Where(condition).Find(&models).Error
	if err == nil {
		return models, nil
	}
	return models, utils.NewError(emsgs.Internal)
}
