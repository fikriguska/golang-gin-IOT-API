package controller

import (
	"log"
	"net/http"
	e "src/error"
	"src/middleware"
	"src/models"
	"src/service/channel_service"
	"src/service/node_service"

	"github.com/gin-gonic/gin"
)

func ChannelRoute(r *gin.Engine) {
	authorized := r.Group("/channel", middleware.JwtMiddleware.MiddlewareFunc())

	authorized.POST("", AddChannel)
}

func AddChannel(c *gin.Context) {
	var json models.ChannelAdd

	// Check required parameter
	if err := c.BindJSON(&json); err != nil {
		log.Println(err)
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
		return
	}

	nodeService := node_service.Node{
		Node: models.Node{
			Id: json.Id_node,
		},
	}

	idUser, isAdmin := extractJwt(c)

	exist, owner := nodeService.IsExistAndOwner(idUser)

	current_node := nodeService.GetNodeOnly()

	if !exist {
		errorResponse(c, http.StatusNotFound, e.ErrNodeIdNotFound)
		return
	} else if !owner && !isAdmin {
		errorResponse(c, http.StatusForbidden, e.ErrUseNodeNotPermitted)
		return
	}

	for i, v := range json.Value {

		if v != nil {
			if current_node.Field_sensor[i] == nil {
				errorResponse(c, http.StatusBadRequest, e.ErrFieldIsEmpty)
				return
			}
		}
	}

	channelService := channel_service.Channel{
		Channel: models.Channel{
			Value:   json.Value,
			Id_node: json.Id_node,
		},
	}

	channelService.Add()

	successResponse(c, http.StatusCreated, "Success add channel")

}
