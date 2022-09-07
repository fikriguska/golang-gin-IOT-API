package controller

import (

	// "net/http/httptest"

	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	e "src/error"
	"src/models"
	"src/util"
	"strconv"

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

func insertSensor(s models.Sensor) int {
	statement := "insert into sensor (name, unit, id_node, id_hardware) values ($1, $2, $3, $4) returning id_sensor"
	var id int
	err := db.QueryRow(statement, s.Name, s.Unit, s.Id_node, s.Id_hardware).Scan(&id)
	e.PanicIfNeeded(err)
	return id
}

// func autoInsertSensor() {
// }

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

	// another user
	// sensor2 := randomSensor()

	node2 := randomNode()
	hardware2 := randomHardware()
	hardware2.Type = "sensor"
	id_hardware2 := insertHardware(hardware2)
	user2 := randomUser()
	user2.Status = true
	id_user2 := insertUser(user2)
	node2.Id_hardware = id_hardware2
	node2.Id_user = id_user2
	id_node2 := insertNode(node2)

	// another hardware typed not a sensor
	hardware3 := randomHardware()
	hardware3.Type = "single-board computer"
	id_hardware3 := insertHardware(hardware3)

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
		{
			name: "ok with hardware",
			body: gin.H{
				"name":        sensor.Name,
				"unit":        sensor.Unit,
				"id_node":     id_node,
				"id_hardware": id_hardware,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "node doesnt exist",
			body: gin.H{
				"name":        sensor.Name,
				"unit":        sensor.Unit,
				"id_node":     1337,
				"id_hardware": id_hardware,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				checkBody(t, recorder, e.ErrNodeNotFound)
			},
		},
		{
			name: "hardware doesnt exist",
			body: gin.H{
				"name":        sensor.Name,
				"unit":        sensor.Unit,
				"id_node":     id_node,
				"id_hardware": 1337,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				checkBody(t, recorder, e.ErrHardwareNotFound)
			},
		},
		{
			name: "hardware is not a sensor",
			body: gin.H{
				"name":        sensor.Name,
				"unit":        sensor.Unit,
				"id_node":     id_node,
				"id_hardware": id_hardware3,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkBody(t, recorder, e.ErrHardwareMustbeSensor)
			},
		},
		{
			name: "using another user's node",
			body: gin.H{
				"name":        sensor.Name,
				"unit":        sensor.Unit,
				"id_node":     id_node2,
				"id_hardware": id_hardware,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				checkBody(t, recorder, e.ErrUseNodeNotPermitted)
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

func TestDeleteSensor(t *testing.T) {
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
	sensor.Id_hardware = id_hardware
	sensor.Id_node = id_node

	id_sensor := insertSensor(sensor)

	testCases := []struct {
		name          string
		id            int
		user          testUser
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "not exist",
			id:   1337,
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				checkBody(t, recorder, e.ErrSensorNotFound)
			},
		},
		{
			name: "ok",
			id:   id_sensor,
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/sensor/"+strconv.Itoa(tc.id), nil)
			req.SetBasicAuth(tc.user.Username, tc.user.Password)
			log.Println(req.Header)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}
