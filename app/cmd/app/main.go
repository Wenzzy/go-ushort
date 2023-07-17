package main

import (
	"go-ushorter/app/common/database"
	"go-ushorter/app/common/logger"
	"go-ushorter/app/config"
	"go-ushorter/app/routers"
	"time"
)

//	@title			Go-Ushorter API
//	@version		1.0
//	@description	This is a sample server for create short urls.

//	@contact.name	Wenzzy Belkov
//	@contact.url	https://github.com/WenzzyX
//	@contact.email	wenzzy.belkov@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:80
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Access-token

func main() {

	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}

	cfg := config.GetCfg()

	loc, _ := time.LoadLocation(cfg.Server.Timezone)
	time.Local = loc

	DBDSN := config.GetDbConfiguration()
	_, err := database.DbConnection(DBDSN)
	if err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}

	router := routers.SetupRoute()
	logger.Fatalf("%v", router.Run(config.GetRunServerConfig()))
}
