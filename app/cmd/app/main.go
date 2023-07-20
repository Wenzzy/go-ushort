package main

import (
	"go-ushort/app/common/database"
	"go-ushort/app/common/logger"
	"go-ushort/app/common/utils"
	"go-ushort/app/config"
	"go-ushort/app/routers"
	"log"
	"time"
)

//	@title			go-ushort API
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
//	@name						Authorization

func main() {

	if err := config.SetupConfig("configs/.env"); err != nil {
		log.Fatalf("config SetupConfig() error: %s", err)
	}
	logger.InitLogger()

	cfg := config.GetCfg()

	loc, _ := time.LoadLocation(cfg.Server.Timezone)
	time.Local = loc

	DBDSN := config.GetDbConfiguration()
	_, err := database.DbConnection(DBDSN)
	if err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}

	router := routers.SetupRouter()
	utils.SetupValidatorOptions()
	logger.Fatalf("%v", router.Run(config.GetRunServerConfig()))
}
