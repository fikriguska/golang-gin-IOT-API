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
	Realm:       "test zone",
	Key:         []byte("s3cr3tz_k3y"),
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
