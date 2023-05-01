package controller

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func extractJwt(c *gin.Context) (int, bool) {
	claims := jwt.ExtractClaims(c)

	return int(claims["id"].(float64)), claims["is_admin"].(bool)
}

func errorResponse(c *gin.Context, statusCode int, err error) {
	// c.IndentedJSON(statusCode, gin.H{
	// 	"status": "error",
	// 	"data":   err.Error(),
	// })
	c.String(statusCode, err.Error())
}

func successResponse(c *gin.Context, statusCode int, msg string) {
	c.IndentedJSON(statusCode, msg)
}
