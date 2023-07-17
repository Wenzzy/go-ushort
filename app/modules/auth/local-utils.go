package auth

import (
	"github.com/gin-gonic/gin"
	"go-ushorter/app/config"
	"time"
)

func SetRefreshToContext(c *gin.Context, refreshToken string) error {
	cfg := config.GetCfg()
	refreshDuration, err := time.ParseDuration(cfg.Server.JwtRefreshExpTime)
	if err != nil {
		return err
	}
	c.SetCookie("refresh_token", refreshToken, int(refreshDuration.Seconds()), "/", cfg.Server.Domain, cfg.Server.IsProduction, true)
	return nil
}

func ClearRefreshInContext(c *gin.Context) {
	cfg := config.GetCfg()
	c.SetCookie("refresh_token", "", 0, "/", cfg.Server.Domain, cfg.Server.IsProduction, true)
}
