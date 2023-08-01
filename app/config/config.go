package config

import (
	"fmt"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type ServerConfiguration struct {
	IsProduction         bool   `mapstructure:"IS_PRODUCTION"`
	IsEnableProm         bool   `mapstructure:"IS_ENABLE_PROM"`
	IsDebug              bool   `mapstructure:"IS_DEBUG"`
	AllowedHosts         string `mapstructure:"ALLOWED_HOSTS"`
	AllowedOrigins       string `mapstructure:"ALLOWED_ORIGINS"`
	Domain               string `mapstructure:"DOMAIN" validate:"required"`
	Host                 string `mapstructure:"SERVER_HOST"`
	Port                 string `mapstructure:"SERVER_PORT"`
	Timezone             string `mapstructure:"SERVER_TIMEZONE"`
	JwtAccessSecret      string `mapstructure:"JWT_ACCESS_SECRET" validate:"required"`
	JwtAccessExpTime     string `mapstructure:"JWT_ACCESS_EXP_TIME" validate:"required"`
	JwtRefreshSecret     string `mapstructure:"JWT_REFRESH_SECRET" validate:"required"`
	JwtRefreshExpTime    string `mapstructure:"JWT_REFRESH_EXP_TIME" validate:"required"`
	LimitCountPerRequest int64  `mapstructure:"LIMIT_COUNT_PER_REQUEST"`
}

type DatabaseConfiguration struct {
	Name           string `mapstructure:"DB_NAME" validate:"required"`
	User           string `mapstructure:"DB_USER" validate:"required"`
	Pass           string `mapstructure:"DB_PASS" validate:"required"`
	Host           string `mapstructure:"DB_HOST"`
	Port           string `mapstructure:"DB_PORT"`
	SslMode        string `mapstructure:"DB_SSL_MODE"`
	LogMode        bool   `mapstructure:"DB_LOG_MODE"`
	MigrationsPath string `mapstructure:"MIGRATIONS_PATH"`
}

type Configuration struct {
	Server   ServerConfiguration   `mapstructure:",squash"`
	Database DatabaseConfiguration `mapstructure:",squash"`
}

var (
	Cfg *Configuration
)

func SetupConfig(configPath string) error {
	var configuration *Configuration
	viper.SetConfigFile(configPath)
	viper.SetDefault("IS_PRODUCTION", true)
	viper.SetDefault("IS_DEBUG", false)
	viper.SetDefault("IS_ENABLE_PROM", true)
	viper.SetDefault("ALLOWED_HOSTS", "0.0.0.0")
	viper.SetDefault("ALLOWED_ORIGINS", "*")
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", "8000")
	viper.SetDefault("SERVER_TIMEZONE", "Europe/Berlin")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("MIGRATIONS_PATH", "./app/common/database/migrations")
	viper.BindEnv("DOMAIN")
	viper.BindEnv("JWT_ACCESS_SECRET")
	viper.BindEnv("JWT_ACCESS_EXP_TIME")
	viper.BindEnv("JWT_REFRESH_SECRET")
	viper.BindEnv("JWT_REFRESH_EXP_TIME")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASS")

	viper.AutomaticEnv()

	if configPath != "" {
		// Check if the file exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Printf("Config file '%s' does not exist, reading configuration from environment variables.", configPath)
		} else {
			viper.SetConfigFile(configPath)

			if err := viper.ReadInConfig(); err != nil {
				log.Printf("Error reading config file: %v", err)
			}
		}
	}

	if err := viper.Unmarshal(&configuration); err != nil {
		log.Fatalf("error to decode, %v", err)
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

func GetDbConfiguration() string {
	DBDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		Cfg.Database.Host,
		Cfg.Database.User,
		Cfg.Database.Pass,
		Cfg.Database.Name,
		Cfg.Database.Port,
		Cfg.Database.SslMode,
	)

	return DBDSN
}

func GetRunServerConfig() string {
	appServer := fmt.Sprintf("%s:%s", Cfg.Server.Host, Cfg.Server.Port)
	log.Print("Server Running at: ", appServer)
	return appServer
}
