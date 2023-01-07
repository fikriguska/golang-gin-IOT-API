package main

import (
	_ "net/http/pprof"
	"src/config"
	"src/controller"
	"src/models"

	"github.com/felixge/fgprof"
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
	r.GET("/debug/fgprof", gin.WrapH(fgprof.Handler()))
	r.Run()
}
