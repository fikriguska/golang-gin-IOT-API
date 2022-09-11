package controller

import (
	"log"
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
	authorized.GET("/", ListNode)
	authorized.GET("/:id", GetNode)
	authorized.DELETE("/:id", DeleteNode)
}

func AddNode(c *gin.Context) {
	var json models.NodeAdd

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

		isNode := hardwareService.CheckHardwareType("node")

		if !isNode {
			errorResponse(c, http.StatusBadRequest, e.ErrHardwareMustbeSensor)
			return
		}

	} else {
		nodeService.Id_hardware = -1
	}

	nodeService.Add()

	successResponse(c, http.StatusCreated, "success add new node")

}

func GetNode(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrUserExist)
		return
	}

	nodeService := node_service.Node{
		Node: models.Node{
			Id: id,
		},
	}
	id_user, _ := c.Get("id_user")
	is_admin, _ := c.Get("is_admin")
	log.Println(id_user, is_admin)

	exist, owner := nodeService.IsExistAndOwner(id_user.(int))

	if !exist {
		errorResponse(c, http.StatusNotFound, e.ErrNodeNotFound)
		return
	} else if !owner && !is_admin.(bool) {
		errorResponse(c, http.StatusForbidden, e.ErrSeeNodeNotPermitted)
		return
	}

	node := nodeService.Get()
	c.IndentedJSON(http.StatusOK, node)

}

func ListNode(c *gin.Context) {
	nodeService := node_service.Node{}

	id_user, _ := c.Get("id_user")
	is_admin, _ := c.Get("is_admin")

	nodes := nodeService.GetAll(id_user.(int), is_admin.(bool))

	c.IndentedJSON(http.StatusOK, nodes)
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
