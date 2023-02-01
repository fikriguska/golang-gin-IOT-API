package controller

import (
	"fmt"
	"net/http"
	e "src/error"
	"src/middleware"
	"src/models"
	"src/service/hardware_service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HardwareRoute(r *gin.Engine) {

	authorized := r.Group("/hardware", middleware.JwtMiddleware.MiddlewareFunc())
	authorized.POST("", AddHardware)
	authorized.GET("", ListHardware)
	authorized.GET("/:id", GetHardware)
	authorized.PUT("/:id", UpdateHardware)
	authorized.DELETE("/:id", DeleteHardware)
}

func AddHardware(c *gin.Context) {

	_, isAdmin := extractJwt(c)

	if !isAdmin {
		errorResponse(c, http.StatusUnauthorized, e.ErrNotAdministrator)
		return
	}

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

	successResponse(c, http.StatusCreated, "Success add new hardware")
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
		errorResponse(c, http.StatusNotFound, e.ErrHardwareIdNotFound)
		return
	}

	hardware := hardwareService.Get()
	c.IndentedJSON(http.StatusOK, hardware)
}

func UpdateHardware(c *gin.Context) {

	_, isAdmin := extractJwt(c)

	if !isAdmin {
		errorResponse(c, http.StatusUnauthorized, e.ErrNotAdministrator)
		return
	}

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
		errorResponse(c, http.StatusNotFound, e.ErrHardwareIdNotFound)
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

	successResponse(c, http.StatusOK, "Success edit hardware")

}

func DeleteHardware(c *gin.Context) {

	_, isAdmin := extractJwt(c)

	if !isAdmin {
		errorResponse(c, http.StatusUnauthorized, e.ErrNotAdministrator)
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
		errorResponse(c, http.StatusNotFound, e.ErrHardwareIdNotFound)
		return
	}

	// stillUsed := hardwareService.IsStillUsed()

	err = hardwareService.Delete()

	if err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrHardwareStillUsed)
		return
	}

	successResponse(c, http.StatusOK, fmt.Sprintf("Delete hardware, id: %d", id))

}
