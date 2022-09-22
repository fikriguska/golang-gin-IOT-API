package controller

import (
	"net/http/httptest"
	"src/config"
	"src/middleware"
	"src/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var router *gin.Engine
var db *pgxpool.Pool

// type response struct {
// 	Status string
// 	Data   string
// }

func TestMain(m *testing.M) {
	cfg := config.Setup()
	middleware.JwtAuth()
	db = models.Setup(cfg)
	router = SetupRouter()

	user := testUser{
		User: models.User{
			Password: "wkwk",
			Email:    "bintangf00code@gmail.com",
			Username: "fikriguska",
			Status:   true,
			Token:    "dea9d35db1b4b85bcf21ec8a3088720d0a50174193606da47a47ec0ff750f21d",
			Is_admin: true,
		},
		hashedPass: "4499c41eec361a4d8c208b5da66870e1f0ee57ef2cc6fd80d0df5fc9d81b7682",
	}
	insertUser(user)

	m.Run()
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	UserRoute(r)
	HardwareRoute(r)
	NodeRoute(r)
	SensorRoute(r)
	ChannelRoute(r)
	return r
}

func setAuth(req *http.Request, username string, password string) {
	req.SetBasicAuth(username, password)
}

func checkErrorBody(t *testing.T, recorder *httptest.ResponseRecorder, e error) {
	require.Equal(t, e.Error(), recorder.Body.String())
}
