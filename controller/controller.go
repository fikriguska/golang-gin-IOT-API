package controller

import (
	"reflect"
	"src/models"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func extractJwt(c *gin.Context) (int, bool) {
	jwt.ExtractClaims(c)
	user, _ := c.Get("identity")
	return user.(*models.User).Id, user.(*models.User).Is_admin
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
