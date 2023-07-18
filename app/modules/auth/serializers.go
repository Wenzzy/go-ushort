package auth

import (
	"github.com/gin-gonic/gin"
	"go-ushorter/app/common/utils"
	"go-ushorter/app/models"
)

type AuthSerializer struct {
	C *gin.Context
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (self *AuthSerializer) Response() AuthResponse {
	model := self.C.MustGet("userModel").(models.UserModel)
	access, refresh := utils.GenAuthTokens(model.ID)
	r := AuthResponse{
		access,
		refresh,
	}
	return r
}
