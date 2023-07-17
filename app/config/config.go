package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go-ushorter/app/common/logger"
	"log"
)

type ServerConfiguration struct {
	Host                 string `mapstructure:"SERVER_HOST"`
	Port                 string `mapstructure:"SERVER_PORT"`
	Timezone             string `mapstructure:"SERVER_TIMEZONE"`
	JwtSecret            string `mapstructure:"JWT_SECRET" validate:"required"`
	LimitCountPerRequest int64  `mapstructure:"LIMIT_COUNT_PER_REQUEST"`
}

type Configuration struct {
	Server   ServerConfiguration   `mapstructure:",squash"`
	Database DatabaseConfiguration `mapstructure:",squash"`
}

type DatabaseConfiguration struct {
	Name    string `mapstructure:"DB_NAME" validate:"required"`
	User    string `mapstructure:"DB_USER" validate:"required"`
	Pass    string `mapstructure:"DB_PASS" validate:"required"`
	Host    string `mapstructure:"DB_HOST" validate:"required"`
	Port    string `mapstructure:"DB_PORT" validate:"required"`
	SslMode string `mapstructure:"DB_SSL_MODE"`
	LogMode bool   `mapstructure:"DB_LOG_MODE"`
}

var (
	Cfg *Configuration
)

func SetupConfig() error {
	var configuration *Configuration
	viper.SetConfigFile("configs/.env")
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", "8000")
	viper.SetDefault("SERVER_TIMEZONE", "Europe/Berlin")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("Error to reading config file", err)
		return err
	}

	if err := viper.Unmarshal(&configuration); err != nil {
		logger.Errorf("error to decode, %v", err)
		return err
	}

	validate := validator.New()
	if err := validate.Struct(&configuration.Server); err != nil {
		log.Fatalf("Missing required attributes in Server env %v\n", err)
	}
	if err := validate.Struct(&configuration.Database); err != nil {
		log.Fatalf("Missing required attributes in Database env %v\n", err)
	}
	Cfg = configuration
	return nil

}
func GetCfg() *Configuration {
	return Cfg
}
