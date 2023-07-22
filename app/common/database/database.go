package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pressly/goose/v3"
	"github.com/wenzzyx/go-ushort/app/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	err error
)

func DbConnection(DSN string) (*gorm.DB, error) {
	var db = DB
	cfg := config.GetCfg()

	logLevel := gormLogger.Silent

	if cfg.Database.LogMode || cfg.Server.IsDebug {
		logLevel = gormLogger.Info
	}

	db, err = gorm.Open(postgres.Open(DSN), &gorm.Config{
		TranslateError: true,
		Logger:         gormLogger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Fatalf("Db connection error")
		return nil, err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Set dialect for goose error")
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {

	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error while getting CWD")
		return nil, err
	}

	if err := goose.Up(sqlDB, filepath.Join(cwd, cfg.Database.MigrationsPath)); err != nil {
		log.Fatalf("Apply migrations error")
		return nil, err
	}

	DB = db
	return db, nil

}

func GetDB() *gorm.DB {
	return DB
}
