package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	GRPCPort   string
	GrinexURL  string
}

func LoadConfig() *Config {
	cfg := &Config{}

	cfg.DBHost = getEnv("DB_HOST", "localhost")
	cfg.DBPort, _ = strconv.Atoi(getEnv("DB_PORT", "5432"))
	cfg.DBUser = getEnv("DB_USER", "postgres")
	cfg.DBPassword = getEnv("DB_PASSWORD", "password")
	cfg.DBName = getEnv("DB_NAME", "usdtrate")
	cfg.GRPCPort = getEnv("GRPC_PORT", "50051")
	cfg.GrinexURL = getEnv("GRINEX_URL", "https://grinex.io/api/v2/depth")

	flag.StringVar(&cfg.DBHost, "dbhost", cfg.DBHost, "database host")
	flag.IntVar(&cfg.DBPort, "dbport", cfg.DBPort, "database port")
	flag.StringVar(&cfg.DBUser, "dbuser", cfg.DBUser, "database user")
	flag.StringVar(&cfg.DBPassword, "dbpassword", cfg.DBPassword, "database password")
	flag.StringVar(&cfg.DBName, "dbname", cfg.DBName, "database name")
	flag.StringVar(&cfg.GRPCPort, "grpcport", cfg.GRPCPort, "grpc server port")
	flag.StringVar(&cfg.GrinexURL, "grinexurl", cfg.GrinexURL, "grinex API url")

	flag.Parse()
	return cfg
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
