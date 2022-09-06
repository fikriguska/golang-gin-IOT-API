package controller

import (

	// "net/http/httptest"

	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"src/models"
	"src/util"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func randomSensor() models.Sensor {
	return models.Sensor{
		Name: util.RandomString(10),
		Unit: util.RandomString(7),
	}
}

func TestAddSensor(t *testing.T) {
	sensor := randomSensor()
	node := randomNode()
	hardware := randomHardware()
	hardware.Type = "sensor"
	id_hardware := insertHardware(hardware)

	user := randomUser()
	user.Status = true
	id_user := insertUser(user)

	node.Id_hardware = id_hardware
	node.Id_user = id_user
	id_node := insertNode(node)

	testCases := []struct {
		name          string
		body          gin.H
		user          testUser
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ok no hardware",
			body: gin.H{
				"name":    sensor.Name,
				"unit":    sensor.Unit,
				"id_node": id_node,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			data, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("POST", "/sensor/", bytes.NewBuffer(data))
			req.SetBasicAuth(tc.user.Username, tc.user.Password)
			log.Println(req.Header)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}
