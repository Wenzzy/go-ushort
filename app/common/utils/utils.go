package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go-ushorter/app/common/constants/emsgs"
	"go-ushorter/app/config"
	"strings"
	"time"
)

func GenAuthTokens(id uint) (string, string) {
	cfg := config.GetCfg()
	access_token_duration, _ := time.ParseDuration(cfg.Server.JwtAccessExpTime)
	refresh_token_duration, _ := time.ParseDuration(cfg.Server.JwtRefreshExpTime)

	access_token := jwt.New(jwt.GetSigningMethod("HS256"))
	access_token.Claims = jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(access_token_duration).Unix(),
	}

	signed_access_token, _ := access_token.SignedString([]byte(cfg.Server.JwtAccessSecret))

	refresh_token := jwt.New(jwt.GetSigningMethod("HS256"))
	refresh_token.Claims = jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(refresh_token_duration).Unix(),
	}

	signed_refresh_token, _ := access_token.SignedString([]byte(cfg.Server.JwtRefreshSecret))

	return signed_access_token, signed_refresh_token
}

func msgForValidationTag(tag string) string {
	switch tag {
	case "required":
		return emsgs.IsEmpty
	case "email":
		return emsgs.NotEmail

	}
	return fmt.Sprintf("unhandled-validation-tag:%s", tag)
}

type CommonError struct {
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
}

type CommonValidationError struct {
	Errors map[string]string `json:"validation-errors"`
}

func (ce *CommonError) Error() string {
	return ce.Message
}

func NewValidatorError(err error) CommonValidationError {
	res := CommonValidationError{}
	res.Errors = make(map[string]string)
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		res.Errors[v.Field()] = msgForValidationTag(v.Tag())
	}
	return res
}

func NewError(errMsg string, descriptionStrings ...string) *CommonError {
	res := CommonError{}
	res.Message = errMsg
	if len(descriptionStrings) > 0 {
		res.Description = strings.Join(descriptionStrings[:], ",")
	}

	return &res
}

// Changed the c.MustBindWith() ->  c.ShouldBindWith().
// I don't want to auto return 400 when error happened.
// origin function is here: https://github.com/gin-gonic/gin/blob/master/context.go
func Bind(c *gin.Context, obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindWith(obj, b)
}
