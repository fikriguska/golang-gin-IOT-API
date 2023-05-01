package controller

import (
	"reflect"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func extractJwt(c *gin.Context) (int, bool) {
	claims := jwt.ExtractClaims(c)

	return int(claims["id"].(float64)), claims["is_admin"].(bool)
}

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
