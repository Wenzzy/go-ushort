package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go-ushorter/app/common/constants/emsgs"
	"go-ushorter/app/common/db-utils"
	"go-ushorter/app/common/utils"
	"go-ushorter/app/config"
	"go-ushorter/app/models"
	"net/http"
)

// Registration godoc
//
//	@Summary	User registration
//	@Schemes
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Success	200		{object}	AuthResponse
//	@Failure	422		{object}	utils.CommonValidationError
//	@Param		request	body		RegisterValidator	true	"Request Body"
//	@Router		/auth/registration [post]
func Registration(c *gin.Context) {
	validator := NewRegisterValidator()
	if err := validator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.NewValidatorError(err))
		return
	}

	if err := db_utils.SaveOne(&validator.UserModel, "user"); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	c.Set("user_model", validator.UserModel)
	sz := AuthSerializer{c}
	res := sz.Response()

	if err := SetRefreshToContext(c, res.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewError(emsgs.Internal, "Can't set refresh token"))
		return
	}
	c.JSON(http.StatusCreated, res)
}

// Login godoc
//
//	@Summary	User login
//	@Schemes
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Success	200		{object}	AuthResponse
//	@Failure	400		{object}	utils.CommonError
//	@Failure	422		{object}	utils.CommonValidationError
//	@Param		request	body		LoginValidator	true	"Request Body"
//	@Router		/auth/login [post]
func Login(c *gin.Context) {
	validator := NewLoginValidator()
	if err := validator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.NewValidatorError(err))
		return
	}
	user, err := models.FindOneUser(&models.UserModel{Email: validator.Email})

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.NewError(emsgs.WrongCredentials))
		return
	}
	if user.ComparePassword(validator.Password) != nil {
		c.JSON(http.StatusBadRequest, utils.NewError(emsgs.WrongCredentials))
		return
	}

	c.Set("user_model", user)
	sz := AuthSerializer{c}
	res := sz.Response()
	if err := SetRefreshToContext(c, res.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewError(emsgs.Internal, "Can't set refresh token"))
		return
	}
	c.JSON(http.StatusCreated, res)
}

// Refresh godoc
//
//	@Summary	Refresh token pair
//	@Schemes
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Success	200		{object}	AuthResponse
//	@Failure	401		{object}	utils.CommonError
//	@Param		cookie	header		string	true	"Cookie"	default(refresh_token="...")
//	@Router		/auth/refresh [post]
func Refresh(c *gin.Context) {
	refresh_token, err := c.Request.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.NewError(emsgs.Unauthorized))
		return
	}

	token, err := jwt.Parse(refresh_token.Value, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		cfg := config.GetCfg()
		return []byte(cfg.Server.JwtRefreshSecret), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.NewError(emsgs.Unauthorized))
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, utils.NewError(emsgs.Unauthorized))
		return
	}
	userId := uint(claims["id"].(float64))
	user, cerr := models.FindOneUser(&models.UserModel{ID: userId})
	if cerr != nil {
		c.JSON(http.StatusBadRequest, utils.NewError(emsgs.WrongCredentials))
		return
	}

	c.Set("user_model", user)
	sz := AuthSerializer{c}
	res := sz.Response()
	if err := SetRefreshToContext(c, res.RefreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, utils.NewError(emsgs.Internal, "Can't set refresh token"))
		return
	}
	c.JSON(http.StatusCreated, res)
}

// Logout godoc
//
//	@Summary	Logout (clear refresh token from cookies)
//	@Schemes
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Success	204		"Empty response"
//	@Param		cookie	header	string	false	"Required Cookie"	default(refresh_token="...")
//	@Router		/auth/logout [delete]
func Logout(c *gin.Context) {
	ClearRefreshInContext(c)
	c.Status(http.StatusNoContent)
}
