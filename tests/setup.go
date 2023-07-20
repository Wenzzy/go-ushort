package tests

import (
	"go-ushort/app/common/database"
	"go-ushort/app/common/logger"
	"go-ushort/app/config"
	"log"
)

func Setup() {
	if err := config.SetupConfig("../configs/tests.env"); err != nil {
		log.Fatalf("config SetupConfig() error: %s", err)
	}
	logger.InitLogger()
	DBDSN := config.GetDbConfiguration()
	_, err := database.DbConnection(DBDSN)
	if err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}

}
