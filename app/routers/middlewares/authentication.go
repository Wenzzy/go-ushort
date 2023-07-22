package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/wenzzyx/go-ushort/app/common/constants/emsgs"
	"github.com/wenzzyx/go-ushort/app/common/database"
	"github.com/wenzzyx/go-ushort/app/common/utils"
	"github.com/wenzzyx/go-ushort/app/config"
	"github.com/wenzzyx/go-ushort/app/models"
)

func UpdateContextUserModel(c *gin.Context, userId uint) {
	var userModel models.UserModel
	if userId != 0 {
		db := database.GetDB()
		db.First(&userModel, userId)
	}
	c.Set("userId", userId)
	c.Set("userModel", userModel)
}

func extractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	rawBearerToken := c.Request.Header.Get("Authorization")
	bearerTokenSplited := strings.Split(rawBearerToken, " ")
	if len(bearerTokenSplited) == 2 {
		return bearerTokenSplited[1]
	}
	return ""
}

func AuthMiddleware(abortWithError bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		UpdateContextUserModel(c, 0)
		tokenStr := extractToken(c)
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			cfg := config.GetCfg()
			return []byte(cfg.Server.JwtAccessSecret), nil
		})
		if err != nil {
			if abortWithError {
				c.JSON(utils.NewError(http.StatusUnauthorized, emsgs.Unauthorized).H())
				c.Abort()
			}
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId := uint(claims["id"].(float64))
			UpdateContextUserModel(c, userId)
		}
	}
}
