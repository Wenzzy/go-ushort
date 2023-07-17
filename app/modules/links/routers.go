package links

import (
	"github.com/gin-gonic/gin"
	"go-ushorter/app/routers/middlewares"
)

func Router(route *gin.RouterGroup) {
	route.Use(middlewares.AuthMiddleware(true))
	route.GET("/", GetAllLinks)

}
