package controller

import (
	"reflect"

	"github.com/gin-gonic/gin"
)

func errorResponse(c *gin.Context, statusCode int, err error) {
	c.String(statusCode, err.Error())
}

func successResponse(c *gin.Context, statusCode int, msg interface{}) {
	if reflect.ValueOf(msg).Kind() == reflect.String {
		c.String(statusCode, msg.(string))
	} else {
		c.IndentedJSON(statusCode, msg)
	}
}
