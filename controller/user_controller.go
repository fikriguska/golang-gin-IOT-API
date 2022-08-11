package controller

import (
	"fmt"
	"net/http"
	"net/mail"
	e "src/error"
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
	r.POST("/user", controller.Create)
	r.GET("/user", controller.Test)
}

func (controller *UserController) Create(c *gin.Context) {
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "pong",
	// })
	var user repository.User

	// Check required parameter
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrInvalidParams.Error(),
		})
		return
	}

	// Check email format
	if !isEmailValid(user.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrInvalidEmail.Error(),
		})
		return
	}

	controller.UserService.IsExist(user)

	// res := controller.UserService.Create(user)
	// c.JSON(http.StatusOK, gin.H{
	// 	"status": "OK",
	// 	"data":   res,
	// })
}

// type Person struct {
// Name     string    `form:"name"`
// }

func (controller *UserController) Test(c *gin.Context) {
	// var p Person
	// fmt.Println(c.QParam("email"))
	// c.ShouldBind()
	fmt.Println("xxxx")
}

func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
