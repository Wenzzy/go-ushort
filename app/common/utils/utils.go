package utils

import (
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/wenzzyx/go-ushort/app/common/constants/emsgs"
	"github.com/wenzzyx/go-ushort/app/config"
	"github.com/xhit/go-str2duration/v2"
)

func GenAuthTokens(id uint) (string, string) {
	cfg := config.GetCfg()
	accessTokenDuration, _ := str2duration.ParseDuration(cfg.Server.JwtAccessExpTime)
	refreshTokenDuration, _ := str2duration.ParseDuration(cfg.Server.JwtRefreshExpTime)

	access_token := jwt.New(jwt.GetSigningMethod("HS256"))
	access_token.Claims = jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(accessTokenDuration).Unix(),
	}

	signed_access_token, _ := access_token.SignedString([]byte(cfg.Server.JwtAccessSecret))

	refresh_token := jwt.New(jwt.GetSigningMethod("HS256"))
	refresh_token.Claims = jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(refreshTokenDuration).Unix(),
	}

	signed_refresh_token, _ := refresh_token.SignedString([]byte(cfg.Server.JwtRefreshSecret))

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
	StatusCode  int    `json:"-"`
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
}

type CommonValidationError struct {
	Errors map[string]string `json:"validation-errors"`
}

func (ce *CommonError) Error() string {
	return ce.Message
}

func (ce *CommonError) H() (int, *CommonError) {
	return ce.StatusCode, ce
}

func NewValidatorError(err error) CommonValidationError {

	if reflect.TypeOf(err) == reflect.TypeOf(validator.ValidationErrors{}) {
		res := CommonValidationError{}
		res.Errors = make(map[string]string)
		errs := err.(validator.ValidationErrors)
		for _, v := range errs {
			res.Errors[v.Field()] = msgForValidationTag(v.Tag())
		}
		return res
	}
	return CommonValidationError{
		Errors: map[string]string{"body": emsgs.NotObject},
	}

}

func NewError(statusCode int, errMsg string, descriptionStrings ...string) *CommonError {
	res := CommonError{}
	res.StatusCode = statusCode
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

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

var realRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenRandomString(length int) string {
	res := make([]byte, length)
	for i := range res {
		res[i] = charset[realRand.Intn(len(charset))]
	}
	return string(res)
}

func SetupValidatorOptions() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}
