package controller

import (
	"net/http"
	e "src/error"
	"src/models"
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

	userService := user_service.User{
		User: models.User{
			Email:    json.Email,
			Username: json.Username,
			Password: json.Password,
		},
	}

	// Check email format
	if !userService.IsEmailValid() {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrInvalidEmail.Error(),
		})
		return
	}

	exist := userService.IsExist()

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
		User: models.User{
			Token: token,
		},
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
		User: models.User{
			Username: json.Username,
			Password: json.Password,
		},
	}

	credCorrect, activated := userService.Auth()

	if !credCorrect {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrUsernameOrPassIncorrect.Error(),
		})
		return
	} else if !activated {
		c.JSON(http.StatusBadRequest, gin.H{
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
