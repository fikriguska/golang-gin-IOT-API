package controller

import (
	"database/sql"
	"encoding/json"
	"net/http/httptest"
	"src/config"
	"src/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var router *gin.Engine
var db *sql.DB

type response struct {
	Status string
	Data   string
}

func TestMain(m *testing.M) {
	cfg := config.Setup()
	db = models.Setup(cfg)
	router = SetupRouter()
	m.Run()
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	UserRoute(r)
	HardwareRoute(r)
	NodeRoute(r)
	return r
}

func checkBody(t *testing.T, recorder *httptest.ResponseRecorder, e error) {
	var resp response
	json.Unmarshal(recorder.Body.Bytes(), &resp)
	require.Equal(t, e.Error(), resp.Data)
}
