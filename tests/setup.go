package tests

import (
	"log"

	"github.com/wenzzyx/go-ushort/app/common/database"
	"github.com/wenzzyx/go-ushort/app/common/logger"
	"github.com/wenzzyx/go-ushort/app/config"
)

func Setup() {
	if err := config.SetupConfig("../../configs/tests.env"); err != nil {
		log.Fatalf("config SetupConfig() error: %s", err)
	}
	logger.InitLogger()
	DBDSN := config.GetDbConfiguration()
	_, err := database.DbConnection(DBDSN)
	if err != nil {
		logger.Fatalf("database DbConnection error: %s", err)
	}

}
