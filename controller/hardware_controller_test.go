package controller

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"

	// "net/http/httptest"
	e "src/error"
	"src/models"
	"src/util"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func randomHardware() models.Hardware {

	hardwareType := []string{"single-board computer", "microcontroller unit", "sensor"}
	return models.Hardware{
		Name:        util.RandomString(10),
		Type:        hardwareType[rand.Intn(3)],
		Description: util.RandomString(20),
	}
}

func insertHardware(h models.Hardware) int {
	statement := "insert into hardware (name, type, description) values ($1, $2, $3) returning id_hardware"
	var id int
	err := db.QueryRow(statement, h.Name, h.Type, h.Description).Scan(&id)
	e.PanicIfNeeded(err)
	return id
}

func TestAddHardware(t *testing.T) {
	hardware := randomHardware()
	testCases := []struct {
		name          string
		body          gin.H
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			body: gin.H{
				"name":        hardware.Name,
				"type":        hardware.Type,
				"description": hardware.Description,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "invalid hardware type",
			body: gin.H{
				"name":        hardware.Name,
				"type":        "namikaze satellite",
				"description": hardware.Description,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkBody(t, recorder, e.ErrInvalidHardwareType)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			data, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("POST", "/hardware", bytes.NewBuffer(data))
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}

func TestDeleteHardware(t *testing.T) {
	hardware := randomHardware()

	id := insertHardware(hardware)

	testCases := []struct {
		name          string
		id            int
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			id:   id,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "not found hardware",
			id:   1337,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				checkBody(t, recorder, e.ErrHardwareNotFound)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/hardware/"+strconv.Itoa(tc.id), nil)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}
