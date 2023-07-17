package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
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
