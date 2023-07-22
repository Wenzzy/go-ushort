package db_utils

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wenzzyx/go-ushort/app/common/constants/emsgs"
	"github.com/wenzzyx/go-ushort/app/common/database"
	"github.com/wenzzyx/go-ushort/app/common/utils"
	"github.com/wenzzyx/go-ushort/app/models"
	"gorm.io/gorm"
)

func SaveOne(data any, descriptionStrings ...string) *utils.CommonError {
	db := database.GetDB()
	err := db.Save(data).Error
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return utils.NewError(http.StatusBadRequest, emsgs.ObjectAlreadyExists)
	}
	return utils.NewError(http.StatusInternalServerError, emsgs.Internal)
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
