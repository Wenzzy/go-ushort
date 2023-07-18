package db_utils

import (
	"errors"
	"fmt"
	"go-ushorter/app/common/constants/emsgs"
	"go-ushorter/app/common/database"
	"go-ushorter/app/common/utils"
	"go-ushorter/app/models"
	"gorm.io/gorm"
)

func SaveOne(data any, descriptionStrings ...string) *utils.CommonError {
	db := database.GetDB()
	err := db.Save(data).Error
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return utils.NewError(emsgs.ObjectAlreadyExists, descriptionStrings...)
	}
	return utils.NewError(emsgs.Internal)
}

func uniqueLinkAliasGenerator(db *gorm.DB, aliasLength int) string {
	var model *models.LinkModel

	randomString := utils.GenRandomString(aliasLength)
	if err := db.Where(&models.LinkModel{GeneratedAlias: randomString}).First(&model); err == nil {
		return uniqueLinkAliasGenerator(db, aliasLength)
	}
	return randomString
}

func GenUniqueLinkAlias() string {
	db := database.GetDB()
	var count int64
	db.Model(&models.LinkModel{}).Count(&count)

	alias := uniqueLinkAliasGenerator(db, len(fmt.Sprintf("%d", count))+1)
	return alias
}
