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
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrInvalidParams.Error(),
		})
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
			c.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"data":   e.ErrHardwareNotFound.Error(),
			})
			return
		}

	} else {
		nodeService.Id_hardware = -1
	}

	nodeService.Add()

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
		"data":   "success add node",
	})
}

func DeleteNode(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrInvalidParams.Error(),
		})
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
		c.JSON(http.StatusNotFound, gin.H{
			"status": "error",
			"data":   e.ErrNodeNotFound.Error(),
		})
		return
	} else if !owner {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "error",
			"data":   e.ErrDeleteNodeNotPermitted.Error(),
		})
		return
	}

	nodeService.Delete()

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   "success delete node",
	})

}
