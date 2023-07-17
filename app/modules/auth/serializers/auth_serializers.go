package serializers

import (
	"github.com/gin-gonic/gin"
	"go-ushorter/app/common/utils"
	"go-ushorter/app/models"
)

type RegisterSerializer struct {
	C *gin.Context
}

type RegisterResponse struct {
	ID    uint   `json:"ID"`
	Token string `json:"token"`
}

func (self *RegisterSerializer) Response() RegisterResponse {
	model := self.C.MustGet("user_model").(models.UserModel)
	r := RegisterResponse{
		ID:    model.ID,
		Token: utils.GenToken(model.ID),
	}
	return r
}
