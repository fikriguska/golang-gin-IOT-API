package controller

import "github.com/gin-gonic/gin"

func errorResponse(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, gin.H{
		"status": "error",
		"data":   err.Error(),
	})
}

func successResponse(c *gin.Context, statusCode int, msg string) {
	c.JSON(statusCode, gin.H{
		"status": "ok",
		"data":   msg,
	})
}
