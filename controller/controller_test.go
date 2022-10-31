package controller

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"src/config"
	"src/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var router *gin.Engine
var db *sql.DB

// type response struct {
// 	Status string
// 	Data   string
// }

func TestMain(m *testing.M) {
	cfg := config.Setup()
	db = models.Setup(cfg)
	router = SetupRouter()

	user := testUser{
		User: models.User{
			Password: "wkwk",
			Email:    "bintangf00code@gmail.com",
			Username: "perftest",
			Status:   true,
			Token:    "dbb68d97021afbdb7bf0f2beb87705ecd9073a5737a7ced8c9be4680ee9d3549",
			Is_admin: false,
		},
		hashedPass: "3c31bc6fa467cea84245bf86d594f17936880674f320b94b2cef9f73ac71e51f",
	}
	user.Id = insertUser(user)
	hardwareNode := randomHardwareNode()
	hardwareNode.Id = insertHardware(hardwareNode)

	node := randomNode()
	node.Id_hardware = hardwareNode.Id
	node.Id_user = user.Id
	node.Id = insertNode(node)

	hardwareSensor := randomHardwareSensor()
	hardwareSensor.Id = insertHardware(hardwareSensor)

	sensor := randomSensor()
	sensor.Id_hardware = hardwareSensor.Id
	sensor.Id_node = node.Id

	sensor.Id = insertSensor(sensor)
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
