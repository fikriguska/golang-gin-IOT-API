package config

import (
	"os"
	"strconv"
)

type Configuration struct {
	Port   int
	DBHost string
	DBUser string
	DBPass string
	DBName string
	DBPort int
}

func Setup() Configuration {

	cfg := Configuration{}

	cfg.Port, _ = strconv.Atoi(getEnv("PORT", "8080"))
	cfg.DBHost = getEnv("DB_HOST", "localhost")
	cfg.DBUser = getEnv("DB_USER", "postgres")
	cfg.DBPass = getEnv("DB_PASS", "postgres")
	cfg.DBName = getEnv("DB_NAME", "iot")
	cfg.DBPort, _ = strconv.Atoi(getEnv("DB_PORT", "5432"))

	return cfg
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
