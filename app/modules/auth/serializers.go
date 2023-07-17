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
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (self *AuthSerializer) Response() AuthResponse {
	model := self.C.MustGet("user_model").(models.UserModel)
	access, refresh := utils.GenAuthTokens(model.ID)
	r := AuthResponse{
		access,
		refresh,
	}
	return r
}
