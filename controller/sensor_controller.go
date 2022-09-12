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
	"strconv"

	"github.com/gin-gonic/gin"
)

func SensorRoute(r *gin.Engine) {
	authorized := r.Group("/sensor", middleware.BasicAuth())

	authorized.POST("/", AddSensor)
	authorized.GET("/", ListSensor)
	authorized.GET("/:id", GetSensor)
	authorized.PUT("/:id", UpdateSensor)
	authorized.DELETE("/:id", DeleteSensor)
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
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
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
		errorResponse(c, http.StatusNotFound, e.ErrNodeNotFound)
		return
	} else if !isAdmin.(bool) && !owner {
		errorResponse(c, http.StatusForbidden, e.ErrUseNodeNotPermitted)
		return
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
			return
		}
		sensorService.Id_hardware = *json.Id_hardware

	} else {
		sensorService.Id_hardware = -1
	}

	sensorService.Add()

	successResponse(c, http.StatusCreated, "success add new sensor")

}

func GetSensor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrUserExist)
		return
	}

	sensorService := sensor_service.Sensor{
		Sensor: models.Sensor{
			Id: id,
		},
	}

	id_user, _ := c.Get("id_user")
	is_admin, _ := c.Get("is_admin")

	exist, owner := sensorService.IsExistAndOwner(id_user.(int))

	if !exist {
		errorResponse(c, http.StatusNotFound, e.ErrSensorNotFound)
		return
	} else if !owner && !is_admin.(bool) {
		errorResponse(c, http.StatusForbidden, e.ErrSeeSensorNotPermitted)
		return
	}

	sensor := sensorService.Get()
	c.IndentedJSON(http.StatusOK, sensor)
}

func ListSensor(c *gin.Context) {
	sensorService := sensor_service.Sensor{}
	id_user, _ := c.Get("id_user")
	is_admin, _ := c.Get("is_admin")

	sensors := sensorService.GetAll(id_user.(int), is_admin.(bool))

	c.IndentedJSON(http.StatusOK, sensors)
}

func UpdateSensor(c *gin.Context) {
	var json models.SensorUpdate

	// Check required parameter
	if err := c.ShouldBindJSON(&json); err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
		return
	}

	sensorService := sensor_service.Sensor{
		Sensor: models.Sensor{
			Id: id,
		},
	}

	id_user, _ := c.Get("id_user")
	is_admin, _ := c.Get("is_admin")

	exist, owner := sensorService.IsExistAndOwner(id_user.(int))

	if !exist {
		errorResponse(c, http.StatusNotFound, e.ErrNodeNotFound)
		return
	} else if !owner && !is_admin.(bool) {
		errorResponse(c, http.StatusForbidden, e.ErrEditSensorNotPermitted)
		return
	}

	sensorService.Update(json)

}

func DeleteSensor(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidEmail)
		return
	}

	sensorService := sensor_service.Sensor{
		Sensor: models.Sensor{
			Id: id,
		},
	}
	id_user, _ := c.Get("id_user")
	isAdmin, _ := c.Get("is_admin")
	exist, owner := sensorService.IsExistAndOwner(id_user.(int))

	if !exist {
		errorResponse(c, http.StatusNotFound, e.ErrSensorNotFound)
		return
	} else if !owner && !isAdmin.(bool) {
		errorResponse(c, http.StatusForbidden, e.ErrDeleteSensorNotPermitted)
		return
	}

	sensorService.Delete()
	successResponse(c, http.StatusOK, "success delete sensor")

}
