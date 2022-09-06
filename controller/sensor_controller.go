package controller

import (
	"fmt"
	"net/http"
	e "src/error"
	"src/middleware"
	"src/models"
	"src/service/hardware_service"
	"src/service/node_service"
	"src/service/sensor_service"

	"github.com/gin-gonic/gin"
)

func SensorRoute(r *gin.Engine) {
	authorized := r.Group("/sensor", middleware.BasicAuth())

	authorized.POST("/", AddSensor)
}

type AddSensorStruct struct {
	Name        string `json:"name" binding:"required"`
	Unit        string `json:"unit" binding:"required"`
	Id_Node     int    `json:"id_node" binding:"required"`
	Id_hardware *int   `json:"id_hardware"`
}

func AddSensor(c *gin.Context) {

	var json AddSensorStruct
	// var id_user int
	// var isAdmin bool

	// Check required parameter
	if err := c.BindJSON(&json); err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrUserExist)
		return
	}

	nodeService := node_service.Node{
		Node: models.Node{
			Id: json.Id_Node,
		},
	}

	id_user, _ := c.Get("id_user")
	isAdmin, _ := c.Get("is_admin")

	exist, owner := nodeService.IsExistAndOwner(id_user.(int))

	if !exist {
		fmt.Println("[+] node not exist")
		errorResponse(c, 200, e.ErrNodeNotFound)
		return
	} else if !owner {
		errorResponse(c, http.StatusForbidden, e.ErrDeleteNodeNotPermitted)
		return
	}

	if !isAdmin.(bool) && !owner {
		errorResponse(c, http.StatusForbidden, e.ErrUseNodeNotPermitted)
	}

	sensorService := sensor_service.Sensor{
		Sensor: models.Sensor{
			Name:    json.Name,
			Unit:    json.Unit,
			Id_node: json.Id_Node,
		},
	}

	// check is id_harware passed in request
	if json.Id_hardware != nil {
		hardwareService := hardware_service.Hardware{
			Hardware: models.Hardware{
				Id: *json.Id_hardware,
			},
		}
		hardwareExist := hardwareService.IsExist()
		if !hardwareExist {

			// falcon
			// resp.status = falcon.HTTP_400
			// resp.body = 'Id hardware is invalid'
			fmt.Println("[+] hardware not exist")
			errorResponse(c, http.StatusNotFound, e.ErrHardwareNotFound)
			return
		}
		isSensor := hardwareService.CheckHardwareType("sensor")

		if !isSensor {
			errorResponse(c, http.StatusBadRequest, e.ErrHardwareMustbeSensor)
		}
		sensorService.Id_hardware = *json.Id_hardware

	} else {
		sensorService.Id_hardware = -1
	}

	sensorService.Add()

	successResponse(c, http.StatusCreated, "success add new sensor")

}
