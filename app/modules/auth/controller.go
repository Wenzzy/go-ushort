package auth

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/wenzzyx/go-ushort/app/common/constants/emsgs"
	db_utils "github.com/wenzzyx/go-ushort/app/common/db-utils"
	"github.com/wenzzyx/go-ushort/app/common/utils"
	"github.com/wenzzyx/go-ushort/app/config"
	"github.com/wenzzyx/go-ushort/app/models"
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
		c.JSON(err.H())
		return
	}
	c.Set("userModel", validator.UserModel)
	sz := AuthSerializer{c}
	res := sz.Response()

	if err := SetRefreshToContext(c, res.RefreshToken); err != nil {

		c.JSON(utils.NewError(http.StatusInternalServerError, emsgs.Internal, "Can't set refresh token").H())
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
		c.JSON(utils.NewError(http.StatusBadRequest, emsgs.WrongCredentials).H())
		return
	}
	if user.ComparePassword(validator.Password) != nil {
		c.JSON(utils.NewError(http.StatusBadRequest, emsgs.WrongCredentials).H())
		return
	}

	c.Set("userModel", user)
	sz := AuthSerializer{c}
	res := sz.Response()
	if err := SetRefreshToContext(c, res.RefreshToken); err != nil {
		c.JSON(utils.NewError(http.StatusInternalServerError, emsgs.Internal, "Can't set refresh token").H())
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
//	@Router		/auth/refresh [get]
func Refresh(c *gin.Context) {
	refreshToken, err := c.Request.Cookie("refreshToken")
	if err != nil {
		c.JSON(utils.NewError(http.StatusUnauthorized, emsgs.Unauthorized).H())
		return
	}

	token, err := jwt.Parse(refreshToken.Value, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		cfg := config.GetCfg()
		return []byte(cfg.Server.JwtRefreshSecret), nil
	})
	if err != nil {
		fmt.Println(err)
		c.JSON(utils.NewError(http.StatusUnauthorized, emsgs.Unauthorized).H())
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(utils.NewError(http.StatusUnauthorized, emsgs.Unauthorized).H())
		return
	}
	userId := uint(claims["id"].(float64))
	user, cerr := models.FindOneUser(&models.UserModel{ID: userId})
	if cerr != nil {
		c.JSON(utils.NewError(http.StatusBadRequest, emsgs.WrongCredentials).H())
		return
	}

	c.Set("userModel", user)
	sz := AuthSerializer{c}
	res := sz.Response()
	if err := SetRefreshToContext(c, res.RefreshToken); err != nil {
		c.JSON(utils.NewError(http.StatusInternalServerError, emsgs.Internal, "Can't set refresh token").H())
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
