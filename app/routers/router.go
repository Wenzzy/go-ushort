package routers

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/wenzzyx/go-ushort/app/common/logger"
	"github.com/wenzzyx/go-ushort/app/common/metrics"
	"github.com/wenzzyx/go-ushort/app/config"
)

func SetupRouter() *gin.Engine {
	cfg := config.GetCfg()
	debugMode := cfg.Server.IsDebug
	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	r.Use(gin.Recovery())

	metrics.InitMetrics(r)

	if cfg.Server.IsDebug {
		r.Use(gin.Logger())
	} else {
		r.Use(logger.JsonLoggerMiddleware())
	}

	if err := r.SetTrustedProxies([]string{cfg.Server.AllowedHosts}); err != nil {
		logger.Errorf("Can't set trusted proxies")
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(cfg.Server.AllowedOrigins, ","),
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	RegisterRoutes(r)

	return r
}
