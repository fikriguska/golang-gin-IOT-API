package middleware

import (
	"encoding/base64"
	"net/http"
	"src/models"
	"src/service/user_service"
	"strings"

	e "src/error"

	"github.com/gin-gonic/gin"
)

func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			c.String(http.StatusUnauthorized, "your auth method is not right")
			c.Abort()
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])

		pair := strings.SplitN(string(payload), ":", 2)

		userService := user_service.User{
			User: models.User{
				Username: pair[0],
				Password: pair[1],
			},
		}
		credCorrect, activated := userService.Auth()

		if len(pair) != 2 {
			c.String(http.StatusUnauthorized, "invalid authorization key")
			c.Abort()
			return
		}

		if !credCorrect {
			c.String(http.StatusUnauthorized, e.ErrUsernameOrPassIncorrect.Error())
			c.Abort()
			return
		}

		if !activated {
			c.String(http.StatusForbidden, e.ErrUserNotActive.Error())
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
