package controller

// import (
// 	"net/http"
// 	e "src/error"
// 	"src/service/node_service"

// 	"github.com/gin-gonic/gin"
// )

// func NodeRoute(r *gin.Engine) {
// 	r.POST("/node", AddNode)
// }

// type AddNodeStruct struct {
// 	Name     string
// 	Location string
// }

// func AddNode(c *gin.Context) {
// 	var json AddNodeStruct

// 	// Check required parameter
// 	if err := c.ShouldBindJSON(&json); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status": "error",
// 			"data":   e.ErrInvalidParams.Error(),
// 		})
// 		return
// 	}

// 	nodeService := node_service.Node{
// 		Name:     json.Name,
// 		Location: json.Location,
// 	}

// 	nodeService.Add()

// }
