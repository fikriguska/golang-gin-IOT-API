package controller

import (
	"net/http"
	"src/repository"
	"src/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(userService *service.UserService) UserController {
	return UserController{
		UserService: *userService,
	}
}

func (controller *UserController) Route(r *gin.Engine) {
	r.GET("/user", controller.Create)
}

func (controller *UserController) Create(ctx *gin.Context) {
	// ctx.JSON(http.StatusOK, gin.H{
	// 	"message": "pong",
	// })

	res := controller.UserService.Create(repository.User{})
	ctx.JSON(http.StatusOK, gin.H{
		"Status": "OK",
		"data":   res,
	})
}
