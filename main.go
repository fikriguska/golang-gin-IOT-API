package main

import (
	"net/http"
	"src/cfg"
	"src/models"

	"github.com/gin-gonic/gin"
)

func init() {
	cfg.Setup()
	models.Setup()
	models.Cx()
}

func main() {

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
