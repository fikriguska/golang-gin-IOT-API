package controller

import (
	"log"
	"net/http"
	e "src/error"
	"src/middleware"
	"src/models"
	"src/service/user_service"
	"src/util"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserRoute(r *gin.Engine) {
	r.POST("/user", AddUser)
	r.GET("/user/activation", ActivateUser)
	r.POST("/user/login", Login)
	r.POST("/user/forget-password", ForgetPassword)
	// r.GET("/user/:id", middleware.BasicAuth(), GetUser)

	authorized := r.Group("/user/:id", middleware.BasicAuth())
	authorized.DELETE("", DeleteUser)
	authorized.PUT("", UpdateUser)
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
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
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
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidEmail)
		return
	}

	exist := userService.IsExist()

	if exist {
		errorResponse(c, http.StatusBadRequest, e.ErrUserExist)
		return
	}

	userService.Add()

	successResponse(c, http.StatusCreated, "success sign up, check email for verification")

}

// var ActivateUserStruct
func ActivateUser(c *gin.Context) {
	token, exist := c.GetQuery("token")

	if !exist {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
		return
	}

	userService := user_service.User{
		User: models.User{
			Token: token,
		},
	}

	valid := userService.IsTokenValid()

	if !valid {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidToken)
		return
	}

	userService.Activate()

	successResponse(c, http.StatusOK, "user is activated")

}

type LoginUserStruct struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var json LoginUserStruct

	// Check required parameter
	if err := c.ShouldBindJSON(&json); err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
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
		errorResponse(c, http.StatusBadRequest, e.ErrUsernameOrPassIncorrect)
		return
	} else if !activated {
		errorResponse(c, http.StatusBadRequest, e.ErrUserNotActive)
		return
	}

	successResponse(c, http.StatusOK, "logged in")

}

func ForgetPassword(c *gin.Context) {
	var json models.UserForgetPassword

	// Check required parameter
	if err := c.ShouldBindJSON(&json); err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
		return
	}

	userService := user_service.User{
		User: models.User{
			Email:    json.Email,
			Username: json.Username,
		},
	}

	// Check email format
	if !userService.IsEmailValid() {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidEmail)
		return
	}

	match, userActivated := userService.IsEmailAndUsernameMatched()

	if !match {
		errorResponse(c, http.StatusBadRequest, e.ErrUsernameOrEmailIncorrect)
		return
	}

	if !userActivated {
		errorResponse(c, http.StatusBadRequest, e.ErrUserNotActive)
		return
	}

	userService.SetRandomPassword()

}

func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
		return
	}

	userService := user_service.User{
		User: models.User{
			Id: id,
		},
	}

	exist := userService.IsExist()

	if !exist {
		errorResponse(c, http.StatusNotFound, e.ErrUserNotFound)
		return
	}

	isAdmin, _ := c.Get("is_admin")

	if !isAdmin.(bool) {
		errorResponse(c, http.StatusForbidden, e.ErrDeleteUserNotPermitted)
		return
	}

	isUsingNode := userService.IsUsingNode()
	log.Println(isUsingNode)
	if isUsingNode {
		errorResponse(c, http.StatusBadRequest, e.ErrUserStillUsingNode)
		return
	}

	userService.Delete()

}

func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
		return
	}

	var json models.UserUpdate

	// Check required parameter
	if err := c.ShouldBindJSON(&json); err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
		return
	}

	userService := user_service.User{
		User: models.User{
			Id: id,
		},
	}

	exist := userService.IsExist()

	if !exist {
		errorResponse(c, http.StatusNotFound, e.ErrUserNotFound)
		return
	}

	idUser, _ := c.Get("id_user")

	if idUser.(int) != id {
		errorResponse(c, http.StatusForbidden, e.ErrEditUserNotPermitted)
		return
	}

	oldPasswdHash := util.Sha256String(json.OldPasswd)
	RealOldPasswdHash, _ := c.Get("password")
	if oldPasswdHash != RealOldPasswdHash {
		errorResponse(c, http.StatusBadRequest, e.ErrOldPasswordIncorrect)
		return
	}

	userService.Password = json.NewPasswd
	userService.SetPassword()

}

// func GetUser(c *gin.Context) {
// 	// c.JSON(200, "wllwlw")
// 	id, err := strconv.Atoi(c.Param("id"))

// 	if err != nil {
// 		errorResponse(c, http.StatusBadRequest, e.ErrUserExist)
// 		return
// 	}

// 	is_admin, _ := c.Get("is_admin")
// 	if !is_admin.(bool) {
// 		errorResponse(c, http.StatusForbidden, e.ErrNotAdministrator)
// 	}

// }
