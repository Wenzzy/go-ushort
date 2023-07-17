package routers

import (
	"github.com/gin-gonic/gin"
	"go-ushorter/app/common/logger"
	"go-ushorter/app/config"
	"go-ushorter/app/routers/middlewares"
)

func SetupRoute() *gin.Engine {
	cfg := config.GetCfg()
	debugMode := cfg.Server.IsDebug
	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	allowedHosts := cfg.Server.AllowedHosts
	r := gin.Default()

	if err := r.SetTrustedProxies([]string{allowedHosts}); err != nil {
		logger.Errorf("Can't set trusted proxies")
	}

	r.Use(gin.Recovery())
	r.Use(middlewares.CORSMiddleware())

	RegisterRoutes(r)

	return r
}
