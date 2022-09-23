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

func autoInsertSensor() (testUser, models.Hardware, models.Node, models.Hardware, models.Sensor) {
	sensor := randomSensor()

	user, hardwareNode, node := autoInsertNode(nil)

	hardwareSensor := randomHardwareSensor()
	hardwareSensor.Id = insertHardware(hardwareSensor)

	sensor.Id_hardware = hardwareSensor.Id
	sensor.Id_node = node.Id

	sensor.Id = insertSensor(sensor)
	return user, hardwareNode, node, hardwareSensor, sensor
}

func checkSensorChannel(sensor models.SensorGet, channel models.Channel) bool {
	containsChannel := false
	for _, v := range sensor.Channel {
		if v.Value == channel.Value {
			containsChannel = true
			break
		}
	}
	return containsChannel
}

func TestAddSensor(t *testing.T) {
	sensor := randomSensor()
	hardwareSensor := randomHardwareSensor()
	hardwareSensor.Id = insertHardware(hardwareSensor)

	user, _, node := autoInsertNode(nil)

	_, _, node2 := autoInsertNode(nil)

	// another hardware typed not a sensor
	hardwareNotSensor := randomHardware()
	hardwareNotSensor.Type = "single-board computer"
	hardwareNotSensor.Id = insertHardware(hardwareNotSensor)

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
				"id_hardware": hardwareSensor.Id,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				log.Println(recorder.Body)
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "node doesnt exist",
			body: gin.H{
				"name":        sensor.Name,
				"unit":        sensor.Unit,
				"id_node":     1337,
				"id_hardware": hardwareSensor.Id,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				checkErrorBody(t, recorder, e.ErrNodeIdNotFound)
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
				checkErrorBody(t, recorder, e.ErrHardwareIdNotFound)
			},
		},
		{
			name: "hardware is not a sensor",
			body: gin.H{
				"name":        sensor.Name,
				"unit":        sensor.Unit,
				"id_node":     node.Id,
				"id_hardware": hardwareNotSensor.Id,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkErrorBody(t, recorder, e.ErrHardwareMustbeSensor)
			},
		},
		{
			name: "using another user's node",
			body: gin.H{
				"name":        sensor.Name,
				"unit":        sensor.Unit,
				"id_node":     node2.Id,
				"id_hardware": hardwareSensor.Id,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				checkErrorBody(t, recorder, e.ErrUseNodeNotPermitted)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			data, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("POST", "/sensor/", bytes.NewBuffer(data))
			setAuth(req, tc.user.Username, tc.user.Password)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}

func TestGetChannel(t *testing.T) {
	channel := randomChannel()
	user, _, _, _, sensor := autoInsertSensor()
	channel.Id_sensor = sensor.Id
	insertChannel(channel)

	testCases := []struct {
		name          string
		id            int
		user          testUser
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			id:   sensor.Id,
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var s models.SensorGet
				json.Unmarshal(recorder.Body.Bytes(), &s)
				require.Equal(t, s.Id, sensor.Id)
				require.Equal(t, true, checkSensorChannel(s, channel))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/sensor/"+strconv.Itoa(tc.id), nil)
			setAuth(req, tc.user.Username, tc.user.Password)
			log.Println(req.Header)
			router.ServeHTTP(w, req)
			// log.Println(w.Body)
			// log.Println(sensor)
			// log.Println(channel)
			// log.Println(node)
			tc.checkResponse(w)

		})
	}
}

func TestDeleteSensor(t *testing.T) {

	user, _, _, _, sensor := autoInsertSensor()
	user2, _, _, _, sensor2 := autoInsertSensor()

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
				checkErrorBody(t, recorder, e.ErrSensorIdNotFound)
			},
		},
		{
			name: "forbidden to delete another user's sensor (1)",
			id:   sensor.Id,
			user: user2,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				checkErrorBody(t, recorder, e.ErrDeleteSensorNotPermitted)
			},
		},
		{
			name: "forbidden to delete another user's sensor (2)",
			id:   sensor2.Id,
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				checkErrorBody(t, recorder, e.ErrDeleteSensorNotPermitted)
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
			setAuth(req, tc.user.Username, tc.user.Password)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}
