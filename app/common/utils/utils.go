package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go-ushorter/app/common/constants/emsgs"
	"go-ushorter/app/config"
	"log"
	"strings"
	"time"
)

func GenToken(id uint) string {
	jwt_token := jwt.New(jwt.GetSigningMethod("HS256"))
	// Set some claims
	jwt_token.Claims = jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	cfg := config.GetCfg()
	token, _ := jwt_token.SignedString([]byte(cfg.Server.JwtSecret))
	return token
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
	Error       string `json:"message"`
	Description string `json:"description,omitempty"`
}

type CommonValidationError struct {
	Errors map[string]string `json:"validation-errors"`
}

func NewValidatorError(err error) CommonValidationError {
	res := CommonValidationError{}
	res.Errors = make(map[string]string)
	errs := err.(validator.ValidationErrors)
	for _, v := range errs {
		res.Errors[v.Field()] = msgForValidationTag(v.Tag())
		//if v.Param() != "" {
		//	res.Errors[v.Field()] = fmt.Sprintf("{%v: %v}", v.Tag, v.Param)
		//} else {
		//	res.Errors[v.Field()] = fmt.Sprintf("{key: %v}", v.Tag)
		//}

	}
	return res
}

// Warp the error info in a object
func NewError(errMsg string, descriptionStrings ...string) *CommonError {
	res := CommonError{}
	res.Error = errMsg
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

func GetDbConfiguration() string {
	cfg := config.GetCfg()

	DBDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Pass,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SslMode,
	)

	return DBDSN
}

func GetRunServerConfig() string {
	cfg := config.GetCfg()

	appServer := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Print("Server Running at: ", appServer)
	return appServer
}
