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
	r.DELETE("/hardware/:id", DeleteHardware)
}

type AddHardwareStruct struct {
	Name        string
	Type        string
	Description string
}

func AddHardware(c *gin.Context) {
	var json AddHardwareStruct

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
