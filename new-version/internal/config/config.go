package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/spf13/viper"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
	Pagination  `yaml:"pagination"`
	Security    `yaml:"security"`
	Database    `yaml:"database"`
}

type Database struct {
	Host     string `yaml:"db_host"`
	Port     int    `yaml:"db_port"`
	Name     string `yaml:"db_name"`
	User     string `yaml:"db_user"`
	Password string `yaml:"db_password"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle-timeout" env-default:"60s"`
}

type Pagination struct {
	PageSizeSmall int `yaml:"page_size_small"`
	PageSize      int `yaml:"page_size"`
	PageSizeLarge int `yaml:"page_size_large"`
}

type Security struct {
	PasswordMinLen     int           `yaml:"password_min_len"`
	JwtSecret          string        `yaml:"jwt_secret"`
	AccessTokenExpire  time.Duration `yaml:"access_token_expire"`
	RefreshTokenExpire time.Duration `yaml:"refresh_token_expire"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesn't exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("error during config reading: %v", err)
	}

	return &cfg
}

// load config using viper
func LoadConfig() {
	viper.AddConfigPath("./config")
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("config file does not exist")
		}

		log.Fatalf("error during config reading: %v", err)
	}
}
