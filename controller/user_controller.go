package controller

import (
	"net/http"
	e "src/error"
	"src/model"
	"src/service"

	"github.com/gin-gonic/gin"
)

func (controller *UserController) Route(r *gin.Engine) {
	r.POST("/user", controller.Add)
	r.GET("/user/activation", controller.Activate)
	r.POST("/user/login", controller.Login)
}

type UserController struct {
	UserService service.UserService
}

func NewUserController(userService *service.UserService) UserController {
	return UserController{
		UserService: *userService,
	}
}

func (controller *UserController) Add(c *gin.Context) {

	var request model.AddUserRequest

	// Check required parameter
	if err := c.ShouldBindJSON(&request); err != nil {
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

	user := service.User{
		Email:    request.Email,
		Username: request.Username,
		Password: request.Password,
	}

	exist := controller.UserService.IsExist(user)

	if exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrUserExist.Error(),
		})
		return
	}

	controller.UserService.Add(user)

	c.JSON(http.StatusCreated, gin.H{
		"status": "ok",
		"data":   "Success sign up, check email for verification",
	})

}

// var ActivateUserStruct
func (controller *UserController) Activate(c *gin.Context) {
	token, exist := c.GetQuery("token")

	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrInvalidParams.Error(),
		})
		return
	}

	user := service.User{
		Token: token,
	}

	valid := controller.UserService.IsTokenValid(user)

	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrInvalidToken.Error(),
		})
		return
	}

	controller.UserService.Activate(user)

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   "success",
	})
}

func (controller *UserController) Login(c *gin.Context) {
	var json model.LoginUserRequest

	// Check required parameter
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"data":   e.ErrInvalidParams.Error(),
		})
		return
	}

	user := service.User{
		Username: json.Username,
		Password: json.Password,
	}

	credCorrect, activated := controller.UserService.Auth(user)

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
