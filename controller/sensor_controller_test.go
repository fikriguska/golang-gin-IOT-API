package controller

import (

	// "net/http/httptest"

	"bytes"
	"encoding/json"
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

func autoInsertSensor() (testUser, models.Hardware, models.Node, models.Sensor) {
	sensor := randomSensor()

	user, hardware, node := autoInsertNode("sensor")
	sensor.Id_hardware = hardware.Id
	sensor.Id_node = node.Id

	sensor.Id = insertSensor(sensor)
	return user, hardware, node, sensor
}

func TestAddSensor(t *testing.T) {
	sensor := randomSensor()

	user, hardware, node := autoInsertNode("sensor")

	_, _, node2 := autoInsertNode("sensor")

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
				"id_node": node.Id,
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
				"id_node":     node.Id,
				"id_hardware": hardware.Id,
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
				"id_hardware": hardware.Id,
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
				"id_node":     node.Id,
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
				"id_node":     node.Id,
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
				"id_node":     node2.Id,
				"id_hardware": hardware.Id,
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
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}

func TestDeleteSensor(t *testing.T) {

	user, _, _, sensor := autoInsertSensor()
	user2, _, _, sensor2 := autoInsertSensor()

	// user2 := autoInsertUser()

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
			name: "forbidden to delete another user's sensor (1)",
			id:   sensor.Id,
			user: user2,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				checkBody(t, recorder, e.ErrDeleteSensorNotPermitted)
			},
		},
		{
			name: "forbidden to delete another user's sensor (2)",
			id:   sensor2.Id,
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				checkBody(t, recorder, e.ErrDeleteSensorNotPermitted)
			},
		},
		{
			name: "ok",
			id:   sensor.Id,
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
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}
