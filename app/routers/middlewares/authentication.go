package middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go-ushorter/app/common/constants/emsgs"
	"go-ushorter/app/common/database"
	"go-ushorter/app/common/logger"
	"go-ushorter/app/common/utils"
	"go-ushorter/app/config"
	"go-ushorter/app/models"
	"net/http"
	"strings"
)

func UpdateContextUserModel(c *gin.Context, user_id uint) {
	var userModel models.UserModel
	if user_id != 0 {
		db := database.GetDB()
		db.First(&userModel, user_id)
	}
	c.Set("user_id", user_id)
	c.Set("user_model", userModel)
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
		logger.Infof("%s \n", tokenStr)
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			cfg := config.GetCfg()
			return []byte(cfg.Server.JwtAccessSecret), nil
		})
		if err != nil {
			if abortWithError {
				c.JSON(http.StatusUnauthorized, utils.NewError(emsgs.Unauthorized))
				c.Abort()
			}
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userId := uint(claims["id"].(float64))
			fmt.Printf("%T, %s, %v", userId, userId, userId)
			UpdateContextUserModel(c, userId)
		}
	}
}
