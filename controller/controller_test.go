package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"src/config"
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
	db = models.Setup(cfg)
	router = SetupRouter()

	// user := testUser{
	// 	User: models.User{
	// 		Password: "wkwk",
	// 		Email:    "bintangf00code@gmail.com",
	// 		Username: "perftest",
	// 		Status:   true,
	// 		Token:    "dbb68d97021afbdb7bf0f2beb87705ecd9073a5737a7ced8c9be4680ee9d3549",
	// 		Is_admin: false,
	// 	},
	// 	hashedPass: "3c31bc6fa467cea84245bf86d594f17936880674f320b94b2cef9f73ac71e51f",
	// }
	// user.Id = insertUser(user)
	// hardwareNode := randomHardwareNode()
	// hardwareNode.Id = insertHardware(hardwareNode)

	// node := randomNode()
	// node.Id_hardware = hardwareNode.Id
	// node.Id_user = user.Id
	// node.Id = insertNode(node)

	// hardwareSensor := randomHardwareSensor()
	// hardwareSensor.Id = insertHardware(hardwareSensor)

	// sensor := randomSensor()
	// sensor.Id_hardware = hardwareSensor.Id
	// sensor.Id_node = node.Id

	// sensor.Id = insertSensor(sensor)
	m.Run()
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	UserRoute(r)
	HardwareRoute(r)
	NodeRoute(r)
	ChannelRoute(r)
	return r
}

var jwtToken = make(map[string]string)

type TokenResponse struct {
	Token string
}

func login(username string, password string) string {
	body := gin.H{
		"username": username,
		"password": password,
	}
	w := httptest.NewRecorder()
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/user/login", bytes.NewBuffer(data))
	router.ServeHTTP(w, req)

	var tokenResp TokenResponse
	json.Unmarshal(w.Body.Bytes(), &tokenResp)
	log.Println(tokenResp.Token)

	return tokenResp.Token
}

func setAuth(req *http.Request, username string, password string) {
	var token string
	if _, ok := jwtToken[username]; !ok {
		token = login(username, password)
		jwtToken[username] = token
	}
	log.Println(jwtToken[username])
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", jwtToken[username]))
}

func checkErrorBody(t *testing.T, recorder *httptest.ResponseRecorder, e error) {
	require.Equal(t, e.Error(), recorder.Body.String())
}
