package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DBUrl         string
	GRPCPort      string
	GrinexBaseURL string
	GrinexPath    string
	GrinexMarket  string
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Could not read .env file, fallback to ENV only: %v", err)
	}

	cfg := &Config{
		DBUrl:         viper.GetString("DB_URL"),
		GRPCPort:      viper.GetString("GRPC_PORT"),
		GrinexBaseURL: viper.GetString("GRINEX_API_BASE"),
		GrinexPath:    viper.GetString("GRINEX_DEPTH_PATH"),
		GrinexMarket:  viper.GetString("GRINEX_MARKET"),
	}

	return cfg
}
