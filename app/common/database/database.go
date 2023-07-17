package database

import (
	"github.com/pressly/goose/v3"
	"go-ushorter/app/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path/filepath"
)

var (
	DB  *gorm.DB
	err error
)

func DbConnection(DSN string) (*gorm.DB, error) {
	var db = DB
	cfg := config.GetCfg()

	logLevel := logger.Silent

	if cfg.Database.LogMode || cfg.Server.IsDebug {
		logLevel = logger.Info
	}

	db, err = gorm.Open(postgres.Open(DSN), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logLevel),
	})

	// if need use replicas:
	//if cfg.Database.LogMode || cfg.Server.IsDebug {
	//		db.Use(dbresolver.Register(dbresolver.Config{
	//			Replicas: []gorm.Dialector{
	//				postgres.Open(replicaDSN),
	//			},
	//			Policy: dbresolver.RandomPolicy{},
	//		}))
	//	}

	if err != nil {
		log.Fatalf("Db connection error")
		return nil, err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Message set dialect for goose")
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {

	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Message while getting CWD")
		return nil, err
	}

	if err := goose.Up(sqlDB, filepath.Join(cwd, "app/common/database/migrations")); err != nil {
		log.Fatalf("Message apply migrations")
		return nil, err
	}

	DB = db
	return db, nil

}

func GetDB() *gorm.DB {
	return DB
}
