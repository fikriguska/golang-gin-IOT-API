package controller

import (
	"net/http"
	e "src/error"
	"src/middleware"
	"src/models"
	"src/service/channel_service"
	"src/service/sensor_service"

	"github.com/gin-gonic/gin"
)

func ChannelRoute(r *gin.Engine) {
	authorized := r.Group("/channel", middleware.BasicAuth())

	authorized.POST("", AddChannel)
}

func AddChannel(c *gin.Context) {
	var json models.ChannelAdd

	// Check required parameter
	if err := c.BindJSON(&json); err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
		return
	}

	sensorService := sensor_service.Sensor{
		Sensor: models.Sensor{
			Id: json.Id_sensor,
		},
	}
	id_user, _ := c.Get("id_user")

	exist, owner := sensorService.IsExistAndOwner(id_user.(int))

	if !exist {
		errorResponse(c, http.StatusNotFound, e.ErrSensorIdNotFound)
		return
	} else if !owner {
		errorResponse(c, http.StatusForbidden, e.ErrUseSensorNotPermitted)
		return
	}

	channelService := channel_service.Channel{
		Channel: models.Channel{
			Value:     json.Value,
			Id_sensor: json.Id_sensor,
		},
	}

	channelService.Add()

	successResponse(c, http.StatusCreated, "Success add channel")

}
