package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/xhit/go-str2duration/v2"
	"go-ushorter/app/config"
)

func SetRefreshToContext(c *gin.Context, refreshToken string) error {
	cfg := config.GetCfg()
	refreshDuration, err := str2duration.ParseDuration(cfg.Server.JwtRefreshExpTime)
	if err != nil {
		return err
	}
	c.SetCookie("refreshToken", refreshToken, int(refreshDuration.Seconds()), "/", cfg.Server.Domain, cfg.Server.IsProduction, true)
	return nil
}

func ClearRefreshInContext(c *gin.Context) {
	cfg := config.GetCfg()
	c.SetCookie("refreshToken", "", 0, "/", cfg.Server.Domain, cfg.Server.IsProduction, true)
}
