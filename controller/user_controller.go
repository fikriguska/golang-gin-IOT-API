package controller

import (
	"fmt"
	"net/http"
	e "src/error"
	"src/service/user_service"

	"github.com/gin-gonic/gin"
)

func UserRoute(r *gin.Engine) {
	r.POST("/user", AddUser)
	// r.GET("/user", controller.Test)
}

type AddUserStruct struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func AddUser(c *gin.Context) {
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "pong",
	// })
	fmt.Println("XXXX")
	var json AddUserStruct

	// Check required parameter
	err := c.ShouldBindJSON(&json)
	fmt.Println("YYY")

	if err != nil {
		fmt.Println("ZZZ")

		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrInvalidParams.Error(),
		})
		return
	}

	// Check email format
	// if !isEmailValid(user.Email) {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"status": "error",
	// 		"data":   e.ErrInvalidEmail.Error(),
	// 	})
	// 	return
	// }

	userService := user_service.User{
		Email:    json.Email,
		Username: json.Username,
		Password: json.Password,
	}

	exist, err := userService.IsExist()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"data":   e.ErrAddUserFail.Error(),
		})
		return
	}

	if exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrUserExist.Error(),
		})
		return
	}

	if err := userService.Add(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"data":   e.ErrAddUserFail.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
		"data":   "Success sign up, check email for verification",
	})
	// res := controller.UserService.Create(user)
	// c.JSON(http.StatusOK, gin.H{
	// 	"status": "OK",
	// 	"data":   res,
	// })
}

// type Person struct {
// Name     string    `form:"name"`
// }

// func (controller *UserController) Test(c *gin.Context) {
// 	// var p Person
// 	// fmt.Println(c.QParam("email"))
// 	// c.ShouldBind()
// 	fmt.Println("xxxx")
// }

// func isEmailValid(email string) bool {
// 	_, err := mail.ParseAddress(email)
// 	return err == nil
// }
