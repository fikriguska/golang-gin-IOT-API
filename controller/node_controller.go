package controller

import (
	"fmt"
	"log"
	"net/http"
	e "src/error"
	"src/middleware"
	"src/models"
	"src/service/hardware_service"
	"src/service/node_service"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func NodeRoute(r *gin.Engine) {
	authorized := r.Group("/node", middleware.JwtMiddleware.MiddlewareFunc())

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
		log.Println(err)
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
		return
	}

	jwt.ExtractClaims(c)
	idUser, _ := extractJwt(c)

	nodeService := node_service.Node{
		Node: models.Node{
			Id_user:            idUser,
			Name:               json.Name,
			Location:           json.Location,
			Id_hardware_sensor: json.Id_hardware_sensor,
			Field_sensor:       json.Field_sensor,
		},
	}

	if json.Is_public != nil {
		nodeService.Is_public = *json.Is_public
	} else {
		nodeService.Is_public = false
	}

	// check is id_harware passed in request
	if json.Id_hardware_node != nil {
		nodeService.Id_hardware_node = *json.Id_hardware_node
		hardwareService := hardware_service.Hardware{
			Hardware: models.Hardware{
				Id: *json.Id_hardware_node,
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
		nodeService.Id_hardware_node = -1
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
	idUser, isAdmin := extractJwt(c)

	exist, owner := nodeService.IsExistAndOwner(idUser)

	if !exist {
		errorResponse(c, http.StatusNotFound, e.ErrNodeIdNotFound)
		return
	} else if !owner && !isAdmin {
		errorResponse(c, http.StatusForbidden, e.ErrSeeNodeNotPermitted)
		return
	}

	node := nodeService.Get()
	c.IndentedJSON(http.StatusOK, node)

}

func ListNode(c *gin.Context) {
	nodeService := node_service.Node{}

	idUser, isAdmin := extractJwt(c)

	nodes := nodeService.GetAll(idUser, isAdmin)

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

	idUser, isAdmin := extractJwt(c)

	exist, owner := nodeService.IsExistAndOwner(idUser)

	if !exist {
		errorResponse(c, http.StatusNotFound, e.ErrNodeIdNotFound)
		return
	} else if !owner && !isAdmin {
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

	idUser, _ := extractJwt(c)

	exist, owner := nodeService.IsExistAndOwner(idUser)

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
