package main

import (
	"src/config"
	"src/controller"
	"src/models"
	"src/service/cache_service"

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
	controller.NodeRoute(r)
	controller.SensorRoute(r)
	controller.ChannelRoute(r)
	cache_service.Init()

	r.Run()
}
