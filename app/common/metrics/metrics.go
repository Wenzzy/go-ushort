package metrics

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/wenzzyx/go-ushort/app/config"
)

func InitMetrics(r *gin.Engine) {

	cfg := config.GetCfg()

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"live": "ok"})
	})

	if cfg.Server.IsEnableProm {
		m := ginmetrics.GetMonitor()

		m.SetMetricPath("/metrics")
		m.SetSlowTime(10)
		m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})

		m.Use(r)
	}

}
