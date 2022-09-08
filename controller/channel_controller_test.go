package controller

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	e "src/error"
	"src/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func randomChannel() models.Channel {
	return models.Channel{
		Value: rand.Float64(),
	}
}

func TestAddChannel(t *testing.T) {

	channel := randomChannel()
	user, _, _, sensor := autoInsertSensor()
	channel.Id_sensor = sensor.Id

	user2, _, _, sensor2 := autoInsertSensor()

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
				checkBody(t, recorder, e.ErrSensorNotFound)
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
				checkBody(t, recorder, e.ErrUseSensorNotPermitted)
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
				checkBody(t, recorder, e.ErrUseSensorNotPermitted)
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
			req, _ := http.NewRequest("POST", "/channel/", bytes.NewBuffer(data))
			req.SetBasicAuth(tc.user.Username, tc.user.Password)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}
