package middleware

import (
	"encoding/base64"
	"log"
	"net/http"
	e "src/error"
	"src/models"
	"src/service/user_service"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
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

var JwtMiddleware *jwt.GinJWTMiddleware

func JwtAuth() {
	JwtMiddleware, _ = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: "identity",
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					"id":       v.Id,
					"is_admin": v.Is_admin,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			log.Println(claims)
			log.Println(claims["id"])
			return &models.User{
				Id:       int(claims["id"].(float64)),
				Is_admin: claims["is_admin"].(bool),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			// var loginVals login
			// if err := c.ShouldBind(&loginVals); err != nil {
			// 	return "", jwt.ErrMissingLoginValues
			// }
			// userID := loginVals.Username
			// password := loginVals.Password

			// if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
			// 	return &User{
			// 		UserName:  userID,
			// 		LastName:  "Bo-Yi",
			// 		FirstName: "Wu",
			// 	}, nil
			// }

			// return nil, jwt.ErrFailedAuthentication

			var json models.UserLogin

			// Check required parameter
			if err := c.ShouldBindJSON(&json); err != nil {
				// errorResponse(c, http.StatusBadRequest, e.ErrInvalidParams)
				return nil, e.ErrInvalidParams
			}

			userService := user_service.User{
				User: models.User{
					Username: json.Username,
					Password: json.Password,
				},
			}

			credCorrect, activated := userService.Auth()

			if !credCorrect {
				// errorResponse(c, http.StatusBadRequest, e.ErrUsernameOrPassIncorrect)
				return nil, e.ErrUsernameOrPassIncorrect
			} else if !activated {
				// errorResponse(c, http.StatusBadRequest, e.ErrUserNotActive)
				return nil, e.ErrUserNotActive
			}

			id, _, _, isAdmin := userService.Get()
			return &models.User{
				Id:       id,
				Is_admin: isAdmin,
			}, nil

			// successResponse(c, http.StatusOK, "logged in")
		},
		// Authorizator: func(data interface{}, c *gin.Context) bool {
		// 	if v, ok := data.(*models.User); ok && v.Username == "admin" {
		// 		return true
		// 	}

		// 	return false
		// },
		// Unauthorized: func(c *gin.Context, code int, message string) {
		// 	c.JSON(code, gin.H{
		// 		"code":    code,
		// 		"message": message,
		// 	})
		// },
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

}
