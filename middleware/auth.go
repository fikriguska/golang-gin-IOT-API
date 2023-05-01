package middleware

import (
	"net/http"
	e "src/error"
	"src/models"
	"src/service/user_service"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var JwtMiddleware, _ = jwt.New(&jwt.GinJWTMiddleware{
	SigningAlgorithm: "HS256",
	Key:              []byte("s3cr3tz_k3y"),
	PayloadFunc: func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*models.User); ok {
			return jwt.MapClaims{
				"id":       v.Id,
				"is_admin": v.Is_admin,
			}
		}
		return jwt.MapClaims{}
	},

	Authenticator: func(c *gin.Context) (interface{}, error) {

		var json models.UserLogin

		// Check required parameter
		if err := c.ShouldBindJSON(&json); err != nil {
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
			return nil, e.ErrUsernameOrPassIncorrect
		} else if !activated {
			return nil, e.ErrUserNotActive
		}

		id, _, _, isAdmin := userService.GetForAuth()
		return &models.User{
			Id:       id,
			Is_admin: isAdmin,
		}, nil

	},
	Unauthorized: func(c *gin.Context, _ int, message string) {
		c.String(http.StatusBadRequest, message)
	},

	LoginResponse: func(c *gin.Context, _ int, token string, expire time.Time) {
		c.JSON(http.StatusOK, gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		})
	},
})
