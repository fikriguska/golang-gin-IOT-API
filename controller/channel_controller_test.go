package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	e "src/error"
	"src/models"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func randomChannel() models.Channel {
	return models.Channel{
		Value: rand.Float64(),
	}
}

func insertChannel(c models.Channel) {
	c.Time = time.Now()
	statement := "insert into channel (time, value, id_sensor) values (($1), $2, $3)"
	_, err := db.Exec(context.Background(), statement, c.Time, c.Value, c.Id_sensor)
	e.PanicIfNeeded(err)
}

func TestAddChannel(t *testing.T) {

	channel := randomChannel()
	user, _, _, _, sensor := autoInsertSensor()
	channel.Id_sensor = sensor.Id

	user2, _, _, _, sensor2 := autoInsertSensor()

	testCases := []struct {
		name          string
		body          gin.H
		user          testUser
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "no sensor",
			body: gin.H{
				"value":     channel.Value,
				"id_sensor": 1337,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				checkErrorBody(t, recorder, e.ErrSensorIdNotFound)
			},
		},
		{
			name: "forbidden to delete another user's sensor (1)",
			body: gin.H{
				"value":     channel.Value,
				"id_sensor": sensor2.Id,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				checkErrorBody(t, recorder, e.ErrUseSensorNotPermitted)
			},
		},
		{
			name: "forbidden to delete another user's sensor (2)",
			body: gin.H{
				"value":     channel.Value,
				"id_sensor": sensor.Id,
			},
			user: user2,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				checkErrorBody(t, recorder, e.ErrUseSensorNotPermitted)
			},
		},
		{
			name: "ok",
			body: gin.H{
				"value":     channel.Value,
				"id_sensor": sensor.Id,
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
			req, _ := http.NewRequest("POST", "/channel", bytes.NewBuffer(data))
			req.SetBasicAuth(tc.user.Username, tc.user.Password)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}

func TestListChannel(t *testing.T) {
	user := randomUser()
	user.Status = true
	user.Id = insertUser(user)
	// todo testing to check listed node
	testCases := []struct {
		name          string
		user          testUser
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
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
			req, _ := http.NewRequest("GET", "/node", nil)
			req.SetBasicAuth(tc.user.Username, tc.user.Password)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}
