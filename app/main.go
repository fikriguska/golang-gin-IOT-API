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
	r := gin.Default()
	controller.UserRoute(r)
	controller.HardwareRoute(r)
	controller.NodeRoute(r)
	controller.SensorRoute(r)
	controller.ChannelRoute(r)

	r.Run()
}
