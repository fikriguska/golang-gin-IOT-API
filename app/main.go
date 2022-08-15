package main

import (
	"fmt"
	"log"
	"src/config"
	"src/controller"
	"src/models"

	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Setup() config.Configuration {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.Configuration{}

	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("Unable to parse env")
	}

	fmt.Println(cfg)
	return cfg
}

func main() {
	cfg := Setup()
	models.Setup(cfg)
	// UserRepository := repository.NewUserRepository(database)
	// UserService := service.NewUserService(&UserRepository)
	// UserController := controller.NewUserController(&UserService)
	r := gin.Default()
	controller.UserRoute(r)
	r.Run()
}
