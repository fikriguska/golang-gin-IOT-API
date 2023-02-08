

package config

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Configuration struct {
	Port   int    `env:"PORT" binding:"required"`
	DBHost string `env:"DB_HOST" binding:"required"`
	DBUser string `env:"DB_USER" binding:"required"`
	DBPass string `env:"DB_PASS" binding:"required"`
	DBName string `env:"DB_NAME" binding:"required"`
	DBPort int    `env:"DB_PORT" binding:"required"`
}

func Setup() Configuration {
	projectName := regexp.MustCompile(`^(.*` + "golang-gin-IOT-API" + `)`)
	currentWorkDirectory, _ := os.Getwd()
	fmt.Println(currentWorkDirectory)
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env`)

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
