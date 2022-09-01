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
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrInvalidParams.Error(),
		})
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
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrInvalidHardwareType.Error(),
		})
		return
	}

	hardwareService.Add()

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
		"data":   "Success add new hardware",
	})
}

func DeleteHardware(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrInvalidParams.Error(),
		})
		return
	}

	hardwareService := hardware_service.Hardware{
		Hardware: models.Hardware{
			Id: id,
		},
	}

	exist := hardwareService.IsExist()

	if !exist {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"data":   e.ErrHardwareNotFound.Error(),
		})
		return
	}

	hardwareService.Delete()

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   "Success add new hardware",
	})
}
