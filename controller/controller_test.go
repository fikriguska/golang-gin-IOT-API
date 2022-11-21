package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

// func generate() {
// 	nUser := 10
// 	nNode := 10
// 	nSensor := 5
// 	nChannel := 10
// 	for u := 0; u < nUser; u++ {
// 		user := randomUser()
// 		if u == 0 {
// 			user = testUser{
// 				User: models.User{
// 					Password: "wkwk",
// 					Email:    "bintangf00code@gmail.com",
// 					Username: "perftest",
// 					Status:   true,
// 					Token:    "dbb68d97021afbdb7bf0f2beb87705ecd9073a5737a7ced8c9be4680ee9d3549",
// 					Is_admin: false,
// 				},
// 				hashedPass: "3c31bc6fa467cea84245bf86d594f17936880674f320b94b2cef9f73ac71e51f",
// 			}
// 		}
// 		user.Status = true
// 		user.Id = insertUser(user)
// 		fmt.Println(u, ". created user: ", user)
// 		var wgNode sync.WaitGroup
// 		wgNode.Add(nNode)
// 		for n := 0; n < nNode; n++ {
// 			go func() {
// 				hNode := randomHardwareNode()
// 				hNode.Id = insertHardware(hNode)
// 				node := randomNode()
// 				node.Id_user = user.Id
// 				node.Id_hardware = hNode.Id
// 				node.Id = insertNode(node)
// 				fmt.Println(n, ". created node: ", node)
// 				// var wgSensor sync.WaitGroup
// 				// wgSensor.Add(nSensor)
// 				for s := 0; s < nSensor; s++ {
// 					// go func() {
// 					hSensor := randomHardwareSensor()
// 					hSensor.Id = insertHardware(hSensor)
// 					sensor := randomSensor()
// 					sensor.Id_hardware = hSensor.Id
// 					sensor.Id_node = node.Id
// 					sensor.Id = insertSensor(sensor)
// 					fmt.Println(s, ". created sensor: ", sensor)
// 					// var wgChannel sync.WaitGroup
// 					// wgChannel.Add(nChannel)
// 					for c := 0; c < nChannel; c++ {
// 						// go func() {
// 						channel := randomChannel()
// 						channel.Id_sensor = sensor.Id
// 						insertChannel(channel)
// 						fmt.Println(c, ". created channel: ", channel)
// 						// wgChannel.Done()
// 						// }()
// 					}
// 					// wgChannel.Wait()
// 					// wgSensor.Done()
// 					// }()
// 				}
// 				// wgSensor.Wait()
// 				wgNode.Done()
// 			}()
// 		}
// 		wgNode.Wait()

// 	}
// }

func SetupRouter() *gin.Engine {
	r := gin.Default()
	UserRoute(r)
	HardwareRoute(r)
	NodeRoute(r)
	SensorRoute(r)
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
