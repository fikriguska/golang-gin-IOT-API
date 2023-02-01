package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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

func randomHardwareNode() models.Hardware {
	hardware := randomHardware()
	hardware.Type = []string{"single-board computer", "microcontroller unit"}[rand.Int()%2]
	return hardware
}

func randomHardwareSensor() models.Hardware {
	hardware := randomHardware()
	hardware.Type = "sensor"
	return hardware
}

func insertHardware(h models.Hardware) int {
	statement := "insert into hardware (name, type, description) values ($1, $2, $3) returning id_hardware"
	var id int
	err := db.QueryRow(context.Background(), statement, h.Name, h.Type, h.Description).Scan(&id)
	e.PanicIfNeeded(err)
	return id
}

func chekHardwareSensorBody(t *testing.T, recorder *httptest.ResponseRecorder, h models.Hardware) {
	var hardware models.HardwareSensorGet
	json.Unmarshal(recorder.Body.Bytes(), &hardware)
	// checkHardware(t, hardware, h)
	require.Equal(t, hardware.Id, h.Id)
	require.Equal(t, hardware.Name, h.Name)
	require.Equal(t, hardware.Type, h.Type)
	require.Equal(t, hardware.Description, h.Description)
}

func chekHardwareNodeBody(t *testing.T, recorder *httptest.ResponseRecorder, h models.Hardware, n models.Node) {
	var hardware models.HardwareNodeGet
	json.Unmarshal(recorder.Body.Bytes(), &hardware)
	// checkHardware(t, hardware, h)
	require.Equal(t, hardware.Id, h.Id)
	require.Equal(t, hardware.Name, h.Name)
	require.Equal(t, hardware.Type, h.Type)
	require.Equal(t, hardware.Description, h.Description)
	require.Equal(t, hardware.Node.Name, n.Name)
	require.Equal(t, hardware.Node.Location, n.Location)
}

func autoInsertSensor() models.Hardware {

	hardwareSensor := randomHardwareSensor()
	hardwareSensor.Id = insertHardware(hardwareSensor)

	return hardwareSensor
}

func TestAddHardware(t *testing.T) {
	user := randomUser()
	user.Is_admin = true
	user.Status = true
	insertUser(user)

	user2 := randomUser()
	user2.Status = true
	insertUser(user2)

	hardware := randomHardware()
	testCases := []struct {
		name          string
		body          gin.H
		user          testUser
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			body: gin.H{
				"name":        hardware.Name,
				"type":        hardware.Type,
				"description": hardware.Description,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "user is not admin",
			body: gin.H{
				"name":        hardware.Name,
				"type":        hardware.Type,
				"description": hardware.Description,
			},
			user: user2,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				checkErrorBody(t, recorder, e.ErrNotAdministrator)
			},
		},
		{
			name: "invalid hardware type",
			body: gin.H{
				"name":        hardware.Name,
				"type":        "namikaze satellite",
				"description": hardware.Description,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkErrorBody(t, recorder, e.ErrInvalidHardwareType)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			data, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("POST", "/hardware", bytes.NewBuffer(data))
			fmt.Println(tc.user.Username, tc.user.Password)
			setAuth(req, tc.user.Username, tc.user.Password)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}

func TestDeleteHardware(t *testing.T) {
	hardware := randomHardware()

	id := insertHardware(hardware)

	_, hardware2, _ := autoInsertNode(nil)

	user := randomUser()
	user.Status = true
	insertUser(user)

	useradmin := randomUser()
	useradmin.Status = true
	useradmin.Is_admin = true
	insertUser(useradmin)

	testCases := []struct {
		name          string
		id            int
		user          testUser
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			id:   id,
			user: useradmin,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "user is not admin",
			id:   id,
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				checkErrorBody(t, recorder, e.ErrNotAdministrator)

			},
		},
		{
			name: "not found hardware",
			id:   1337,
			user: useradmin,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				checkErrorBody(t, recorder, e.ErrHardwareIdNotFound)
			},
		},
		{
			name: "hardware is still used",
			id:   hardware2.Id,
			user: useradmin,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkErrorBody(t, recorder, e.ErrHardwareStillUsed)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/hardware/"+strconv.Itoa(tc.id), nil)
			setAuth(req, tc.user.Username, tc.user.Password)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}

func TestGetHardwareSensor(t *testing.T) {
	hardware := autoInsertSensor()
	user := randomUser()
	user.Status = true
	insertUser(user)

	user2 := randomUser()
	user2.Status = true
	insertUser(user2)

	useradmin := randomUser()
	useradmin.Status = true
	useradmin.Is_admin = true
	insertUser(useradmin)

	testCases := []struct {
		name          string
		id            int
		hardware      models.Hardware
		user          testUser
		checkResponse func(recoder *httptest.ResponseRecorder, h models.Hardware)
	}{
		{
			name:     "ok",
			id:       hardware.Id,
			hardware: hardware,
			user:     user,
			checkResponse: func(recorder *httptest.ResponseRecorder, h models.Hardware) {
				require.Equal(t, http.StatusOK, recorder.Code)
				chekHardwareSensorBody(t, recorder, h)
			},
		},
		{
			name:     "ok using another user",
			id:       hardware.Id,
			hardware: hardware,
			user:     user2,
			checkResponse: func(recorder *httptest.ResponseRecorder, h models.Hardware) {
				require.Equal(t, http.StatusOK, recorder.Code)
				chekHardwareSensorBody(t, recorder, h)
			},
		},
		{
			name:     "ok using admin",
			id:       hardware.Id,
			hardware: hardware,
			user:     useradmin,
			checkResponse: func(recorder *httptest.ResponseRecorder, h models.Hardware) {
				require.Equal(t, http.StatusOK, recorder.Code)
				chekHardwareSensorBody(t, recorder, h)
			},
		},
		{
			name:     "hardware not found",
			id:       1337,
			hardware: hardware,
			user:     user,
			checkResponse: func(recorder *httptest.ResponseRecorder, h models.Hardware) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				checkErrorBody(t, recorder, e.ErrHardwareIdNotFound)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/hardware/"+strconv.Itoa(tc.id), nil)
			setAuth(req, tc.user.Username, tc.user.Password)
			router.ServeHTTP(w, req)
			tc.checkResponse(w, tc.hardware)

		})
	}

}

func TestGetHardwareNode(t *testing.T) {
	user, hardware, node := autoInsertNode(nil)
	hardware2 := randomHardwareNode()
	hardware2.Id = insertHardware(hardware2)

	user2 := randomUser()
	user2.Status = true
	insertUser(user2)

	useradmin := randomUser()
	useradmin.Status = true
	useradmin.Is_admin = true
	insertUser(useradmin)

	testCases := []struct {
		name          string
		id            int
		hardware      models.Hardware
		node          models.Node
		user          testUser
		checkResponse func(recoder *httptest.ResponseRecorder, h models.Hardware, n models.Node)
	}{
		{
			name:     "ok",
			id:       hardware.Id,
			hardware: hardware,
			node:     node,
			user:     user,
			checkResponse: func(recorder *httptest.ResponseRecorder, h models.Hardware, n models.Node) {
				require.Equal(t, http.StatusOK, recorder.Code)
				chekHardwareNodeBody(t, recorder, h, n)
			},
		},
		{
			name:     "ok using another user",
			id:       hardware.Id,
			hardware: hardware,
			node:     node,
			user:     user2,
			checkResponse: func(recorder *httptest.ResponseRecorder, h models.Hardware, n models.Node) {
				require.Equal(t, http.StatusOK, recorder.Code)
				chekHardwareNodeBody(t, recorder, h, n)
			},
		},
		{
			name:     "ok using admin",
			id:       hardware.Id,
			hardware: hardware,
			node:     node,
			user:     useradmin,
			checkResponse: func(recorder *httptest.ResponseRecorder, h models.Hardware, n models.Node) {
				require.Equal(t, http.StatusOK, recorder.Code)
				chekHardwareNodeBody(t, recorder, h, n)
			},
		},
		{
			name:     "hardware not found",
			id:       1337,
			hardware: hardware,
			node:     node,
			user:     user,
			checkResponse: func(recorder *httptest.ResponseRecorder, h models.Hardware, n models.Node) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				// chekHardwareNodeBody(t, recorder, h, n)
				checkErrorBody(t, recorder, e.ErrHardwareIdNotFound)
			},
		},
		{
			name:     "ok but node is empty",
			id:       hardware2.Id,
			hardware: hardware2,
			node:     models.Node{},
			user:     user,
			checkResponse: func(recorder *httptest.ResponseRecorder, h models.Hardware, n models.Node) {
				require.Equal(t, http.StatusOK, recorder.Code)
				chekHardwareNodeBody(t, recorder, h, n)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/hardware/"+strconv.Itoa(tc.id), nil)
			setAuth(req, tc.user.Username, tc.user.Password)
			router.ServeHTTP(w, req)
			tc.checkResponse(w, tc.hardware, tc.node)

		})
	}

}

func TestUpdateHardware(t *testing.T) {
	hardware := randomHardware()
	hardware.Id = insertHardware(hardware)

	hardware2 := randomHardware()
	hardware3 := randomHardware()

	user := randomUser()
	user.Status = true
	insertUser(user)

	useradmin := randomUser()
	useradmin.Status = true
	useradmin.Is_admin = true
	insertUser(useradmin)

	testCases := []struct {
		name          string
		id            int
		body          gin.H
		user          testUser
		checkResponse func(recorder *httptest.ResponseRecorder)
		checkInDB     func(id int)
	}{
		{
			name: "ok update name",
			id:   hardware.Id,
			body: gin.H{
				"name": hardware2.Name,
			},
			user: useradmin,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
			checkInDB: func(id int) {
				h := models.GetHardwareById(id)
				require.Equal(t, hardware2.Name, h.Name)
			},
		},
		{
			name: "user is not admin",
			id:   hardware.Id,
			body: gin.H{
				"name": hardware2.Name,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
				checkErrorBody(t, recorder, e.ErrNotAdministrator)
			},
			checkInDB: func(id int) {
			},
		},
		{
			name: "hardware not found",
			id:   1337,
			body: gin.H{},
			user: useradmin,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				checkErrorBody(t, recorder, e.ErrHardwareIdNotFound)
			},
			checkInDB: func(id int) {},
		},

		{
			name: "ok update type",
			id:   hardware.Id,
			body: gin.H{
				"type": hardware2.Type,
			},
			user: useradmin,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
			checkInDB: func(id int) {
				h := models.GetHardwareById(id)
				require.Equal(t, hardware2.Type, h.Type)
			},
		},

		{
			name: "update type not valid",
			id:   hardware.Id,
			body: gin.H{
				"type": "bakugan",
			},
			user: useradmin,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkErrorBody(t, recorder, e.ErrInvalidHardwareType)
			},
			checkInDB: func(id int) {
			},
		},

		{
			name: "ok update description",
			id:   hardware.Id,
			body: gin.H{
				"description": hardware2.Description,
			},
			user: useradmin,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
			checkInDB: func(id int) {
				h := models.GetHardwareById(id)
				require.Equal(t, hardware2.Description, h.Description)
			},
		},

		{
			name: "ok update all fields",
			id:   hardware.Id,
			body: gin.H{
				"name":        hardware3.Name,
				"type":        hardware3.Type,
				"description": hardware3.Description,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
			user: useradmin,
			checkInDB: func(id int) {
				h := models.GetHardwareById(id)
				require.Equal(t, hardware3.Name, h.Name)
				require.Equal(t, hardware3.Type, h.Type)
				require.Equal(t, hardware3.Description, h.Description)

			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			data, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("PUT", "/hardware/"+strconv.Itoa(tc.id), bytes.NewBuffer(data))
			setAuth(req, tc.user.Username, tc.user.Password)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
			tc.checkInDB(tc.id)
		})
	}

}

func TestListHardware(t *testing.T) {
	hardware := randomHardware()

	hardware.Id = insertHardware(hardware)

	user := randomUser()
	user.Status = true
	insertUser(user)

	useradmin := randomUser()
	useradmin.Status = true
	useradmin.Is_admin = true
	insertUser(useradmin)

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
		{
			name: "ok using admin",
			user: useradmin,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/hardware", nil)
			setAuth(req, tc.user.Username, tc.user.Password)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}
