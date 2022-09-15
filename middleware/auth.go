package middleware

import (
	"encoding/base64"
	"log"
	"net/http"
	"src/models"
	"src/service/user_service"
	"strings"

	"github.com/gin-gonic/gin"
)

func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"data":   "invalid authorization key",
			})
			c.Abort()
			return
		}
		log.Println(auth)

		payload, _ := base64.StdEncoding.DecodeString(auth[1])

		pair := strings.SplitN(string(payload), ":", 2)

		userService := user_service.User{
			User: models.User{
				Username: pair[0],
				Password: pair[1],
			},
		}
		credCorrect, activated := userService.Auth()
		if len(pair) != 2 || !(credCorrect && activated) {

			c.JSON(http.StatusUnauthorized, gin.H{
				"status": "error",
				"data":   "invalid authorization key",
			})
			c.Abort()
			return
		}
		id, _, pass, isAdmin := userService.Get()

		c.Set("id_user", id)
		c.Set("username", pair[0])
		c.Set("password", pass)
		c.Set("is_admin", isAdmin)

		c.Next()
	}
}
