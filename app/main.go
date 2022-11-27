package main

import (
	"runtime"
	"src/config"
	"src/controller"
	"src/models"
	_ "net/http/pprof"
	"github.com/felixge/fgprof"
	"github.com/gin-contrib/pprof"
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
	runtime.SetBlockProfileRate(1)

	pprof.Register(r)
	r.GET("/debug/fgprof", gin.WrapH(fgprof.Handler()))
	r.Run()
}
