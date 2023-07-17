package routers

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-ushorter/app/docs"
	"go-ushorter/app/modules/auth"
	"net/http"
)

func RegisterRoutes(r *gin.Engine) {

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"live": "ok"})
	})
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/api", func(c *gin.Context) {
		c.Redirect(302, "/api/index.html")
	})
	r.GET("/api/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	v1 := r.Group("/api/v1")

	// api routes
	auth.AuthR(v1.Group("/auth"))
}
