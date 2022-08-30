package main

import (
	"src/config"
	"src/controller"
	"src/models"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Setup()
	models.Setup(cfg)
	// UserRepository := repository.NewUserRepository(database)
	// UserService := service.NewUserService(&UserRepository)
	// UserController := controller.NewUserController(&UserService)
	r := gin.Default()
	controller.UserRoute(r)
	controller.HardwareRoute(r)
	r.Run()
}
