package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	e "src/error"
	"src/middleware"
	"src/models"
	"src/service/cache_service"
	"src/service/hardware_service"
	"src/service/node_service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NodeRoute(r *gin.Engine) {
	authorized := r.Group("/node", middleware.BasicAuth())

	authorized.POST("", AddNode)
	authorized.GET("", ListNode)
	authorized.GET("/:id", GetNode)
	authorized.PUT("/:id", UpdateNode)
	authorized.DELETE("/:id", DeleteNode)
}

func AddNode(c *gin.Context) {
	var json models.NodeAdd

	// Check required parameter
	if err := c.BindJSON(&json); err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
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
			errorResponse(c, http.StatusNotFound, e.ErrHardwareIdNotFound)
			return
		}

		isNode := hardwareService.CheckHardwareType("node")

		if !isNode {
			errorResponse(c, http.StatusBadRequest, e.ErrHardwareMustbeNode)
			return
		}

	} else {
		nodeService.Id_hardware = -1
	}

	nodeService.Add()

	successResponse(c, http.StatusCreated, "Success add new node")

}

func GetNode(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
		return
	}

	nodeService := node_service.Node{
		Node: models.Node{
			Id: id,
		},
	}
	id_user, _ := c.Get("id_user")
	is_admin, _ := c.Get("is_admin")

	exist, owner := nodeService.IsExistAndOwner(id_user.(int))

	if !exist {
		errorResponse(c, http.StatusNotFound, e.ErrNodeIdNotFound)
		return
	} else if !owner && !is_admin.(bool) {
		errorResponse(c, http.StatusForbidden, e.ErrSeeNodeNotPermitted)
		return
	}
	key := fmt.Sprintf("%d-node", id)
	nodes_byte, err := cache_service.Cache.Get(key)
	if err != nil {
		node := nodeService.Get()
		nodeJson, _ := json.Marshal(node)
		cache_service.Cache.Set(key, nodeJson)
		// log.Println("not cached")
		c.IndentedJSON(http.StatusOK, node)
	} else {
		// log.Println("cached")
		c.Header("Content-Type", "application/json")
		c.String(http.StatusOK, string(nodes_byte))
	}

}

func ListNode(c *gin.Context) {
	nodeService := node_service.Node{}

	id_user, _ := c.Get("id_user")
	is_admin, _ := c.Get("is_admin")

	nodes := nodeService.GetAll(id_user.(int), is_admin.(bool))

	c.IndentedJSON(http.StatusOK, nodes)
}

func UpdateNode(c *gin.Context) {
	var json models.NodeUpdate

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

	nodeService := node_service.Node{
		Node: models.Node{
			Id: id,
		},
	}

	id_user, _ := c.Get("id_user")
	is_admin, _ := c.Get("is_admin")

	exist, owner := nodeService.IsExistAndOwner(id_user.(int))

	if !exist {
		errorResponse(c, http.StatusNotFound, e.ErrNodeIdNotFound)
		return
	} else if !owner && !is_admin.(bool) {
		errorResponse(c, http.StatusForbidden, e.ErrEditNodeNotPermitted)
		return
	}

	nodeService.Update(json)

	successResponse(c, http.StatusOK, "Success edit node")
}

func DeleteNode(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
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
		errorResponse(c, http.StatusNotFound, e.ErrNodeIdNotFound)
		return
	} else if !owner {
		errorResponse(c, http.StatusForbidden, e.ErrDeleteNodeNotPermitted)
		return
	}

	nodeService.Delete()

	successResponse(c, http.StatusOK, fmt.Sprintf("Success delete node, id: %d", id))

}
