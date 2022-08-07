package cfg

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Application struct {
	Port   int    `env:"PORT", required`
	DBHost string `env:"DB_HOST", required`
	DBUser string `env:"DB_USER", required`
	DBPass string `env:"DB_PASS", required`
	DBName string `env:"DB_NAME", required`
	DBPort int    `env:"DB_PORT", required`
}

var App Application

func Setup() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	App = Application{}

	err = env.Parse(&App)
	if err != nil {
		log.Fatalf("Unable to parse env")
	}

	fmt.Println(App)
}
