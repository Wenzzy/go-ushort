package auth

import (
	"github.com/gin-gonic/gin"
)

func Router(route *gin.RouterGroup) {
	route.POST("/registration", Registration)
	route.POST("/login", Login)
	route.POST("/refresh", Refresh)
	route.POST("/logout", Logout)
}
