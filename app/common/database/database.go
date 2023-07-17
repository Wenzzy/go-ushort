package database

import (
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
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

//	func Migrate() {
//		var migrationModels = []any{&models.Example{}}
//
//		err := database.DB.AutoMigrate(migrationModels...)
//		if err != nil {
//			return
//		}
//	}

func DbConnection(DSN string) error {
	var db = DB
	logMode := viper.GetBool("DB_LOG_MODE")
	//debug := viper.GetBool("DEBUG")

	logLevel := logger.Silent

	if logMode {
		logLevel = logger.Info
	}

	db, err = gorm.Open(postgres.Open(DSN), &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logLevel),
	})

	// if need use replicas:
	//if !debug {
	//		db.Use(dbresolver.Register(dbresolver.Config{
	//			Replicas: []gorm.Dialector{
	//				postgres.Open(replicaDSN),
	//			},
	//			Policy: dbresolver.RandomPolicy{},
	//		}))
	//	}

	if err != nil {
		log.Fatalf("Db connection error")
		return err
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Error set dialect for goose")
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {

	}
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error while getting CWD")
		return err
	}

	if err := goose.Up(sqlDB, filepath.Join(cwd, "app/common/database/migrations")); err != nil {
		log.Fatalf("Error apply migrations")
		return err
	}

	DB = db
	return nil

}

func GetDB() *gorm.DB {
	return DB
}
