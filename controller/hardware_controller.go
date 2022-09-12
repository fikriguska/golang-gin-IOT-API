package controller

import (
	"net/http"
	e "src/error"
	"src/models"
	"src/service/hardware_service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HardwareRoute(r *gin.Engine) {
	r.POST("/hardware", AddHardware)
	r.GET("/hardware", ListHardware)
	r.GET("/hardware/:id", GetHardware)
	r.PUT("/hardware/:id", UpdateHardware)
	r.DELETE("/hardware/:id", DeleteHardware)
}

func AddHardware(c *gin.Context) {
	var json models.HardwareAdd

	// Check required parameter
	if err := c.ShouldBindJSON(&json); err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
		return
	}

	hardwareService := hardware_service.Hardware{
		Hardware: models.Hardware{
			Name:        json.Name,
			Type:        json.Type,
			Description: json.Description,
		},
	}

	valid := hardwareService.IsTypeValid()

	if !valid {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidHardwareType)
		return
	}

	hardwareService.Add()

	successResponse(c, http.StatusCreated, "success add new hardware")
}

func ListHardware(c *gin.Context) {
	hardwareService := hardware_service.Hardware{}
	hardwares := hardwareService.GetAll()
	c.IndentedJSON(http.StatusOK, hardwares)

}

func GetHardware(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
		return
	}

	hardwareService := hardware_service.Hardware{
		Hardware: models.Hardware{
			Id: id,
		},
	}

	exist := hardwareService.IsExist()

	if !exist {
		errorResponse(c, http.StatusNotFound, e.ErrHardwareNotFound)
		return
	}

	hardware := hardwareService.Get()
	c.IndentedJSON(http.StatusOK, hardware)
}

func UpdateHardware(c *gin.Context) {
	var json models.HardwareUpdate

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

	hardwareService := hardware_service.Hardware{
		Hardware: models.Hardware{
			Id: id,
		},
	}

	exist := hardwareService.IsExist()

	if !exist {
		errorResponse(c, http.StatusNotFound, e.ErrHardwareNotFound)
		return
	}

	if json.Type != nil {
		hardwareService.Type = *json.Type
		valid := hardwareService.IsTypeValid()

		if !valid {
			errorResponse(c, http.StatusBadRequest, e.ErrInvalidHardwareType)
			return
		}
	}

	hardwareService.Update(json)

}

func DeleteHardware(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrUserExist)
		return
	}

	hardwareService := hardware_service.Hardware{
		Hardware: models.Hardware{
			Id: id,
		},
	}

	exist := hardwareService.IsExist()

	if !exist {
		errorResponse(c, http.StatusNotFound, e.ErrHardwareNotFound)
		return
	}

	hardwareService.Delete()

	successResponse(c, http.StatusOK, "success delete hardware")

}
