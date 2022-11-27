package main

import (
	"runtime"
	"src/config"
	"src/controller"
	"src/models"

	"github.com/gin-contrib/pprof"
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
	runtime.SetBlockProfileRate(1)
	pprof.Register(r)
	r.Run()
}
