package routers

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/wenzzyx/go-ushort/app/modules/auth"
	"github.com/wenzzyx/go-ushort/app/modules/links"
	"github.com/wenzzyx/go-ushort/docs"
)

func RegisterRoutes(r *gin.Engine) {

	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/api-docs", func(c *gin.Context) {
		c.Redirect(302, "/api-docs/index.html")
	})
	r.GET("/api-docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	v1 := r.Group("/api/v1")

	// api routes
	auth.Router(v1.Group("/auth"))
	links.Router(v1.Group("/links"))

	// alias for short url
	r.GET("/:alias", links.Redirect)
}
