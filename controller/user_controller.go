package controller

import (
	"fmt"
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
	r.POST("/user/signup", AddUser)
	r.GET("/user/activation", ActivateUser)
	r.POST("/user/login", Login)
	r.POST("/user/forget-password", ForgetPassword)
	r.GET("/user", middleware.BasicAuth(), ListUser)

	authorized := r.Group("/user/:id", middleware.BasicAuth())
	authorized.GET("", GetUser)
	authorized.DELETE("", DeleteUser)
	authorized.PUT("", UpdateUser)
}

func AddUser(c *gin.Context) {

	var json models.UserAdd

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
		errorResponse(c, http.StatusBadRequest, e.ErrEmailUsernameAlreadyUsed)
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

	exist, activated := userService.TokenValidation()

	if !exist {
		errorResponse(c, http.StatusNotFound, e.ErrTokenNotFound)
		return
	}

	if activated {
		errorResponse(c, http.StatusBadRequest, e.ErrUserAlreadyActive)
		return
	}

	userService.Activate()

	successResponse(c, http.StatusOK, "your account has been activated")

}

func Login(c *gin.Context) {
	var json models.UserLogin

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
	log.Println(credCorrect, activated)

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

	successResponse(c, http.StatusOK, "forget password request sent. Check email for new password")

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
		errorResponse(c, http.StatusNotFound, e.ErrUserIdNotFound)
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

	successResponse(c, http.StatusOK, "success change password, check your email")

}

func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
		return
	}

	idUser, _ := c.Get("id_user")
	isAdmin, _ := c.Get("is_admin")

	if idUser != id && !isAdmin.(bool) {
		errorResponse(c, http.StatusForbidden, e.ErrDeleteUserNotPermitted)
		return
	}

	userService := user_service.User{
		User: models.User{
			Id: id,
		},
	}

	exist := userService.IsExist()

	if !exist {
		errorResponse(c, http.StatusNotFound, e.ErrUserIdNotFound)
		return
	}

	err = userService.Delete()

	if err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrUserStillUsingNode)
		return
	}

	successResponse(c, http.StatusOK, fmt.Sprintf("delete user, id: %d", id))

}

func ListUser(c *gin.Context) {
	userService := user_service.User{}
	users := userService.GetAll()
	c.IndentedJSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
		return
	}

	is_admin, _ := c.Get("is_admin")
	if !is_admin.(bool) {
		errorResponse(c, http.StatusUnauthorized, e.ErrNotAdministrator)
		return
	}

	userService := user_service.User{
		User: models.User{
			Id: id,
		},
	}
	isExist := userService.IsExist()
	if !isExist {
		errorResponse(c, http.StatusOK, e.ErrUserIdNotFound)
		return
	}
	user := userService.Get()
	c.IndentedJSON(http.StatusOK, user)
}
