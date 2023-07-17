package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-ushorter/app/common/logger"
	"go-ushorter/app/routers/middlewares"
)

func SetupRoute() *gin.Engine {
	environment := viper.GetBool("DEBUG")
	if environment {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	allowedHosts := viper.GetString("ALLOWED_HOSTS")
	r := gin.Default()

	if err := r.SetTrustedProxies([]string{allowedHosts}); err != nil {
		logger.Errorf("Can't set trusted proxies")
	}

	r.Use(gin.Recovery())
	r.Use(middlewares.CORSMiddleware())

	RegisterRoutes(r)

	return r
}
