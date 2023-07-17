package db_utils

import (
	"errors"
	"go-ushorter/app/common/constants/emsgs"
	"go-ushorter/app/common/database"
	"go-ushorter/app/common/utils"
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
