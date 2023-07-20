package db_utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-ushort/app/common/constants/emsgs"
	"go-ushort/app/common/database"
	"go-ushort/app/common/utils"
	"go-ushort/app/models"
	"gorm.io/gorm"
	"strconv"
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

type Pagination struct {
	Take   int `json:"take"`
	Page   int `json:"page"`
	Offset int `json:"offset"`
}

func GenPagination(c *gin.Context) Pagination {
	take, _ := strconv.Atoi(c.DefaultQuery("take", "25"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	if take > 100 {
		take = 100
	}
	return Pagination{
		Take:   take,
		Page:   page,
		Offset: take*page - take,
	}

}
