package controller

import (
	"net/http"
	e "src/error"
	"src/middleware"
	"src/models"
	"src/service/hardware_service"
	"src/service/node_service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NodeRoute(r *gin.Engine) {
	authorized := r.Group("/node", middleware.BasicAuth())

	authorized.POST("/", AddNode)
	authorized.DELETE("/:id", DeleteNode)
}

type AddNodeStruct struct {
	Name        string `json:"name" binding:"required"`
	Location    string `json:"location" binding:"required"`
	Id_hardware *int   `json:"id_hardware"`
}

func AddNode(c *gin.Context) {
	var json AddNodeStruct

	// Check required parameter
	if err := c.BindJSON(&json); err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrUserExist)
		return
	}

	id_user, _ := c.Get("id_user")

	nodeService := node_service.Node{
		Node: models.Node{
			Id_user:  id_user.(int),
			Name:     json.Name,
			Location: json.Location,
		},
	}

	// check is id_harware passed in request
	if json.Id_hardware != nil {
		nodeService.Id_hardware = *json.Id_hardware
		hardwareService := hardware_service.Hardware{
			Hardware: models.Hardware{
				Id: *json.Id_hardware,
			},
		}
		hardwareExist := hardwareService.IsExist()
		if !hardwareExist {
			errorResponse(c, http.StatusNotFound, e.ErrHardwareNotFound)
			return
		}

	} else {
		nodeService.Id_hardware = -1
	}

	nodeService.Add()

	successResponse(c, http.StatusCreated, "success add new node")

}

func DeleteNode(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidEmail)
		return
	}

	nodeService := node_service.Node{
		Node: models.Node{
			Id: id,
		},
	}

	id_user, _ := c.Get("id_user")

	exist, owner := nodeService.IsExistAndOwner(id_user.(int))

	if !exist {
		errorResponse(c, http.StatusNotFound, e.ErrNodeNotFound)
		return
	} else if !owner {
		errorResponse(c, http.StatusForbidden, e.ErrDeleteNodeNotPermitted)
		return
	}

	nodeService.Delete()

	successResponse(c, http.StatusOK, "success delete node")

}
