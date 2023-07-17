package auth

import (
	"github.com/gin-gonic/gin"
)

func AuthR(route *gin.RouterGroup) {
	route.POST("/registration", Registration)

}
