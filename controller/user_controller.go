package controller

import (
	"net/http"
	e "src/error"
	"src/service/user_service"

	"github.com/gin-gonic/gin"
)

func UserRoute(r *gin.Engine) {
	r.POST("/user", AddUser)
	r.GET("/user/activation", ActivateUser)
	r.POST("/user/login", Login)
	// r.GET("/user", controller.Test)
}

type AddUserStruct struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func AddUser(c *gin.Context) {

	var json AddUserStruct

	// Check required parameter
	if err := c.ShouldBindJSON(&json); err != nil {
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

	exist := userService.IsExist()

	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"status": "error",
	// 		"data":   e.ErrAddUserFail.Error(),
	// 	})
	// 	return
	// }

	if exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrUserExist.Error(),
		})
		return
	}

	userService.Add()

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

// var ActivateUserStruct
func ActivateUser(c *gin.Context) {
	token, exist := c.GetQuery("token")

	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrInvalidParams.Error(),
		})
		return
	}

	userService := user_service.User{
		Token: token,
	}

	valid := userService.IsTokenValid()

	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrInvalidToken.Error(),
		})
		return
	}

	userService.Activate()

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   "success",
	})
}

type LoginUserStruct struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var json LoginUserStruct

	// Check required parameter
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrInvalidParams.Error(),
		})
		return
	}

	userService := user_service.User{
		Username: json.Username,
		Password: json.Password,
	}

	credCorrect, activated := userService.Auth()

	if !credCorrect {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "error",
			"data":   e.ErrUsernameOrPassIncorrect.Error(),
		})
		return
	} else if !activated {
		c.JSON(http.StatusForbidden, gin.H{
			"status": "error",
			"data":   e.ErrUserNotActive.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   "Logged in",
	})

}
