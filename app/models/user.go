package models

import (
	"errors"
	"go-ushorter/app/common/constants/emsgs"
	"go-ushorter/app/common/database"
	"go-ushorter/app/common/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type UserModel struct {
	ID           uint        `gorm:"column:id;primary_key"`
	Email        string      `gorm:"column:email;unique_index"`
	PasswordHash string      `gorm:"column:password;not null"`
	Role         string      `gorm:"column:role;default:user"`
	RegisteredAt *time.Time  `gorm:"column:registered_at;not null;default:CURRENT_TIMESTAMP"`
	Links        []LinkModel `gorm:"foreignKey:user_id;references:id"`
}

func (e *UserModel) TableName() string {
	return "user"
}

func (u *UserModel) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty!")
	}
	bytePassword := []byte(password)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.PasswordHash = string(passwordHash)
	return nil
}

func (u *UserModel) ComparePassword(password string) error {
	byteHashedPassword := []byte(u.PasswordHash)
	bytePassword := []byte(password)
	err := bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
	return err
}

func FindOneUser(condition any) (UserModel, *utils.CommonError) {
	db := database.GetDB()
	var model UserModel
	err := db.Where(condition).First(&model).Error
	if err == nil {
		return model, nil
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return model, utils.NewError(emsgs.ObjectAlreadyExists, "user")
	}
	return model, utils.NewError(emsgs.Internal)
}
