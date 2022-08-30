package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Configuration struct {
	Port   int    `env:"PORT", required`
	DBHost string `env:"DB_HOST", required`
	DBUser string `env:"DB_USER", required`
	DBPass string `env:"DB_PASS", required`
	DBName string `env:"DB_NAME", required`
	DBPort int    `env:"DB_PORT", required`
}

func Setup() Configuration {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := Configuration{}

	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("Unable to parse env")
	}

	fmt.Println(cfg)
	return cfg
}
