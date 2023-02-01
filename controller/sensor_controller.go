package controller

// import (
// 	"fmt"
// 	"net/http"
// 	e "src/error"
// 	"src/middleware"
// 	"src/models"
// 	"src/service/hardware_service"
// 	"src/service/node_service"
// 	"src/service/sensor_service"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// func SensorRoute(r *gin.Engine) {
// 	authorized := r.Group("/sensor", middleware.JwtMiddleware.MiddlewareFunc())

// 	authorized.POST("", AddSensor)
// 	authorized.GET("", ListSensor)
// 	authorized.GET("/:id", GetSensor)
// 	authorized.PUT("/:id", UpdateSensor)
// 	authorized.DELETE("/:id", DeleteSensor)
// }

// func AddSensor(c *gin.Context) {

// 	var json models.SensorAdd
// 	// var id_user int
// 	// var isAdmin bool

// 	// Check required parameter
// 	if err := c.BindJSON(&json); err != nil {
// 		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
// 		return
// 	}

// 	nodeService := node_service.Node{
// 		Node: models.Node{
// 			Id: json.Id_Node,
// 		},
// 	}

// 	idUser, isAdmin := extractJwt(c)

// 	exist, owner := nodeService.IsExistAndOwner(idUser)

// 	if !exist {
// 		errorResponse(c, http.StatusNotFound, e.ErrNodeIdNotFound)
// 		return
// 	} else if !isAdmin && !owner {
// 		errorResponse(c, http.StatusForbidden, e.ErrUseNodeNotPermitted)
// 		return
// 	}

// 	sensorService := sensor_service.Sensor{
// 		Sensor: models.Sensor{
// 			Name:    json.Name,
// 			Unit:    json.Unit,
// 			Id_node: json.Id_Node,
// 		},
// 	}

// 	// check is id_harware passed in request
// 	if json.Id_hardware != nil {
// 		hardwareService := hardware_service.Hardware{
// 			Hardware: models.Hardware{
// 				Id: *json.Id_hardware,
// 			},
// 		}
// 		hardwareExist := hardwareService.IsExist()
// 		if !hardwareExist {
// 			errorResponse(c, http.StatusNotFound, e.ErrHardwareIdNotFound)
// 			return
// 		}
// 		isSensor := hardwareService.CheckHardwareType("sensor")

// 		if !isSensor {
// 			errorResponse(c, http.StatusBadRequest, e.ErrHardwareMustbeSensor)
// 			return
// 		}
// 		sensorService.Id_hardware = *json.Id_hardware

// 	} else {
// 		sensorService.Id_hardware = -1
// 	}

// 	sensorService.Add()

// 	successResponse(c, http.StatusCreated, "Success add new sensor")

// }

// func GetSensor(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))

// 	if err != nil {
// 		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
// 		return
// 	}

// 	sensorService := sensor_service.Sensor{
// 		Sensor: models.Sensor{
// 			Id: id,
// 		},
// 	}

// 	idUser, isAdmin := extractJwt(c)

// 	exist, owner := sensorService.IsExistAndOwner(idUser)

// 	if !exist {
// 		errorResponse(c, http.StatusNotFound, e.ErrSensorIdNotFound)
// 		return
// 	} else if !owner && !isAdmin {
// 		errorResponse(c, http.StatusForbidden, e.ErrSeeSensorNotPermitted)
// 		return
// 	}

// 	sensor := sensorService.Get()
// 	c.IndentedJSON(http.StatusOK, sensor)
// }

// func ListSensor(c *gin.Context) {
// 	sensorService := sensor_service.Sensor{}

// 	idUser, isAdmin := extractJwt(c)

// 	sensors := sensorService.GetAll(idUser, isAdmin)

// 	c.IndentedJSON(http.StatusOK, sensors)
// }

// func UpdateSensor(c *gin.Context) {
// 	var json models.SensorUpdate

// 	// Check required parameter
// 	if err := c.ShouldBindJSON(&json); err != nil {
// 		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
// 		return
// 	}

// 	id, err := strconv.Atoi(c.Param("id"))

// 	if err != nil {
// 		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
// 		return
// 	}

// 	sensorService := sensor_service.Sensor{
// 		Sensor: models.Sensor{
// 			Id: id,
// 		},
// 	}

// 	idUser, isAdmin := extractJwt(c)

// 	exist, owner := sensorService.IsExistAndOwner(idUser)

// 	if !exist {
// 		errorResponse(c, http.StatusNotFound, e.ErrNodeIdNotFound)
// 		return
// 	} else if !owner && !isAdmin {
// 		errorResponse(c, http.StatusForbidden, e.ErrEditSensorNotPermitted)
// 		return
// 	}

// 	sensorService.Update(json)

// 	successResponse(c, http.StatusOK, "Success edit sensor data")

// }

// func DeleteSensor(c *gin.Context) {
// 	id, err := strconv.Atoi(c.Param("id"))

// 	if err != nil {
// 		errorResponse(c, http.StatusBadRequest, e.ErrInvalidEmail)
// 		return
// 	}

// 	sensorService := sensor_service.Sensor{
// 		Sensor: models.Sensor{
// 			Id: id,
// 		},
// 	}

// 	idUser, isAdmin := extractJwt(c)

// 	exist, owner := sensorService.IsExistAndOwner(idUser)

// 	if !exist {
// 		errorResponse(c, http.StatusNotFound, e.ErrSensorIdNotFound)
// 		return
// 	} else if !owner && !isAdmin {
// 		errorResponse(c, http.StatusForbidden, e.ErrDeleteSensorNotPermitted)
// 		return
// 	}

// 	sensorService.Delete()

// 	successResponse(c, http.StatusOK, fmt.Sprintf("Success delete sensor data, id: %d", id))

// }
