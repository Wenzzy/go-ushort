package links

import (
	"github.com/gin-gonic/gin"
	"github.com/wenzzyx/go-ushort/app/routers/middlewares"
)

func Router(route *gin.RouterGroup) {
	route.Use(middlewares.AuthMiddleware(true))
	route.GET("/", GetAll)
	route.POST("/", Create)
	route.PATCH("/:linkId", Update)

}
