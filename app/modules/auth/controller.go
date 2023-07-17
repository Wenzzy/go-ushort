package auth

import (
	"github.com/gin-gonic/gin"
	"go-ushorter/app/common/db-utils"
	"go-ushorter/app/common/utils"
	"go-ushorter/app/modules/auth/serializers"
	"go-ushorter/app/modules/auth/validators"
	"net/http"
)

// Registration godoc
//
//	@Summary	User registration
//	@Schemes
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Success	200		{object}	serializers.RegisterResponse
//	@Param		request	body		validators.RegisterValidator	true	"Request Body"
//	@Router		/auth/registration [get]
func Registration(c *gin.Context) {
	validator := validators.NewRegisterValidator()
	if err := validator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.NewValidatorError(err))
		return
	}

	if err := db_utils.SaveOne(&validator.UserModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	c.Set("user_model", validator.UserModel)
	sz := serializers.RegisterSerializer{c}
	c.JSON(http.StatusCreated, sz.Response())
}
