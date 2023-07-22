package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/wenzzyx/go-ushort/app/common/logger"
	"github.com/wenzzyx/go-ushort/app/common/metrics"
	"github.com/wenzzyx/go-ushort/app/config"
	"github.com/wenzzyx/go-ushort/app/routers/middlewares"
)

func SetupRouter() *gin.Engine {
	cfg := config.GetCfg()
	debugMode := cfg.Server.IsDebug
	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	allowedHosts := cfg.Server.AllowedHosts
	r := gin.New()

	r.Use(gin.Recovery())

	metrics.InitMetrics(r)

	if cfg.Server.IsDebug {
		r.Use(gin.Logger())
	} else {
		r.Use(logger.JsonLoggerMiddleware())
	}

	if err := r.SetTrustedProxies([]string{allowedHosts}); err != nil {
		logger.Errorf("Can't set trusted proxies")
	}

	r.Use(middlewares.CORSMiddleware())

	RegisterRoutes(r)

	return r
}
