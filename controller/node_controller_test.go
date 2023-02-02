package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"math/rand"
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

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"math/rand"
// 	"net/http"
// 	"net/http/httptest"
// 	"strconv"

// 	// "net/http/httptest"

// 	e "src/error"
// 	"src/models"
// 	"src/util"

// 	"log"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/require"
// )

func randomNode() models.Node {
	return models.Node{
		Name:     util.RandomString(10),
		Location: util.RandomString(7),
	}
}

func insertNode(n models.Node) int {
	statement := "insert into node (name, location, id_user, id_hardware_node, id_hardware_sensor, field_sensor, is_public) values ($1, $2, $3, $4, $5, $6, $7) returning id_node"
	var id int
	err := db.QueryRow(context.Background(), statement, n.Name, n.Location, n.Id_user, n.Id_hardware_node, n.Id_hardware_sensor, n.Field_sensor, n.Is_public).Scan(&id)
	e.PanicIfNeeded(err)
	return id
}

// // create for user, hardware, node automatically, and insert it to db
func autoInsertNode(hardwareType interface{}, isPublic bool) (testUser, models.Hardware, models.Node) {
	node := randomNode()

	hardwareNode := randomHardware()
	if hardwareType != nil {
		hardwareNode.Type = hardwareType.(string)
	} else {
		hardwareNode.Type = []string{"single-board computer", "microcontroller unit"}[rand.Int()%2]
	}

	hardwareNode.Id = insertHardware(hardwareNode)

	user := randomUser()
	user.Status = true
	user.Id = insertUser(user)

	node.Id_hardware_node = hardwareNode.Id
	node.Id_user = user.Id

	sensor := autoInsertSensor()
	sensor2 := autoInsertSensor()

	node.Id_hardware_sensor = []*int{&sensor.Id, &sensor2.Id, nil, nil, nil, nil, nil, nil, nil, nil}
	field1 := "c"
	field2 := "m"
	node.Field_sensor = []*string{&field1, &field2, nil, nil, nil, nil, nil, nil, nil, nil}

	node.Is_public = isPublic
	node.Id = insertNode(node)
	return user, hardwareNode, node
}

func checkNodeBody(t *testing.T, recorder *httptest.ResponseRecorder, n models.Node) {
	var node models.NodeGet
	json.Unmarshal(recorder.Body.Bytes(), &node)
	require.Equal(t, node.Id, n.Id)
	require.Equal(t, node.Name, n.Name)
	require.Equal(t, node.Location, n.Location)

}

// func checkNodeSensor(node models.NodeGet, sensor models.Sensor) bool {
// 	containsSensor := false
// 	for _, v := range node.Sensor {
// 		if v.Id_sensor == sensor.Id {
// 			containsSensor = true
// 			break
// 		}
// 	}
// 	return containsSensor
// }

// func checkNodeHardware(node models.NodeGet, hardware models.Hardware) bool {
// 	containsHardware := false
// 	for _, v := range node.Hardware {
// 		if v.Name == hardware.Name && v.Type == hardware.Type {
// 			containsHardware = true
// 			break
// 		}
// 	}
// 	return containsHardware
// }

func TestAddNode(t *testing.T) {
	node := randomNode()
	hardware := randomHardwareNode()
	hardware.Id = insertHardware(hardware)

	hardware2 := randomHardware()
	hardware2.Type = "sensor"
	hardware2.Id = insertHardware(hardware2)

	user := randomUser()
	user.Status = true
	insertUser(user)

	useradmin := randomUser()
	useradmin.Status = true
	insertUser(useradmin)

	field := "mm"
	testCases := []struct {
		name          string
		body          gin.H
		user          testUser
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ok no hardware",
			body: gin.H{
				"name":               node.Name,
				"location":           node.Location,
				"field_sensor":       []*string{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
				"id_hardware_sensor": []*int{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "ok with hardware",
			body: gin.H{
				"name":               node.Name,
				"location":           node.Location,
				"id_hardware_node":   hardware.Id,
				"field_sensor":       []*string{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
				"id_hardware_sensor": []*int{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "ok using admin",
			body: gin.H{
				"name":               node.Name,
				"location":           node.Location,
				"id_hardware_node":   hardware.Id,
				"field_sensor":       []*string{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
				"id_hardware_sensor": []*int{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
			},
			user: useradmin,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "field is empty",
			body: gin.H{
				"name":               node.Name,
				"location":           node.Location,
				"id_hardware_node":   hardware.Id,
				"field_sensor":       []*string{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
				"id_hardware_sensor": []*int{&hardware2.Id, nil, nil, nil, nil, nil, nil, nil, nil, nil},
			},
			user: useradmin,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkErrorBody(t, recorder, e.ErrFieldIsEmpty)
			},
		},
		{
			name: "sensor is empty",
			body: gin.H{
				"name":               node.Name,
				"location":           node.Location,
				"id_hardware_node":   hardware.Id,
				"field_sensor":       []*string{&field, nil, nil, nil, nil, nil, nil, nil, nil, nil},
				"id_hardware_sensor": []*int{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
			},
			user: useradmin,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkErrorBody(t, recorder, e.ErrIdSensorIsEmpty)
			},
		},
		{
			name: "wrong type hardware",
			body: gin.H{
				"name":               node.Name,
				"location":           node.Location,
				"id_hardware_node":   hardware2.Id,
				"field_sensor":       []*string{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
				"id_hardware_sensor": []*int{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkErrorBody(t, recorder, e.ErrHardwareMustbeNode)
			},
		},
		{
			name: "not found hardware",
			body: gin.H{
				"name":               node.Name,
				"location":           node.Location,
				"id_hardware_node":   1337,
				"field_sensor":       []*string{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
				"id_hardware_sensor": []*int{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				checkErrorBody(t, recorder, e.ErrHardwareIdNotFound)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			data, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("POST", "/node", bytes.NewBuffer(data))
			setAuth(req, tc.user.Username, tc.user.Password)
			log.Println(req.Header)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}

}

func TestGetNode(t *testing.T) {
	user, _, node := autoInsertNode(nil, false)
	_, _, node2 := autoInsertNode(nil, true)

	user2 := randomUser()
	user2.Status = true
	user2.Id = insertUser(user2)

	user3 := randomUser()
	user3.Status = true
	user3.Is_admin = true
	user3.Id = insertUser(user3)

	testCases := []struct {
		name          string
		id            int
		user          testUser
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			id:   node.Id,
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var n models.NodeGet
				json.Unmarshal(recorder.Body.Bytes(), &n)
				checkNodeBody(t, recorder, node)
			},
		},
		{
			name: "forbidden to access another user's node",
			id:   node.Id,
			user: user2,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				checkErrorBody(t, recorder, e.ErrSeeNodeNotPermitted)
			},
		},
		{
			name: "ok using admin",
			id:   node.Id,
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var n models.NodeGet
				json.Unmarshal(recorder.Body.Bytes(), &n)
				checkNodeBody(t, recorder, node)
			},
		},
		{
			name: "ok view public node",
			id:   node2.Id,
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				var n models.NodeGet
				json.Unmarshal(recorder.Body.Bytes(), &n)
				checkNodeBody(t, recorder, node2)
			},
		},
		{
			name: "node not found",
			id:   1337,
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/node/"+strconv.Itoa(tc.id), nil)
			setAuth(req, tc.user.Username, tc.user.Password)
			log.Println(req.Header)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)

		})
	}
}

func TestDeleteNode(t *testing.T) {

	// create node user 1
	// node := randomNode()
	// hardware := randomHardwareNode()
	// id_hardware := insertHardware(hardware)

	// user := randomUser()
	// user.Status = true

	// id_user := insertUser(user)

	// node.Id_hardware_node = id_hardware
	// node.Id_user = id_user
	// id_node := insertNode(node)
	user, _, node := autoInsertNode(nil, false)

	// create node user 2
	// node2 := randomNode()
	// hardware2 := randomHardwareNode()
	// id_hardware2 := insertHardware(hardware2)

	// user2 := randomUser()
	// user2.Status = true

	// id_user2 := insertUser(user2)
	// node2.Id_hardware_node = id_hardware2
	// node2.Id_user = id_user2

	// id_node2 := insertNode(node2)
	user2, _, node2 := autoInsertNode(nil, false)

	testCases := []struct {
		name          string
		user          testUser
		id            int
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "forbidden to delete another user's node (1)",
			user: user,
			id:   node2.Id,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				checkErrorBody(t, recorder, e.ErrDeleteNodeNotPermitted)
			},
		},
		{
			name: "forbidden to delete another user's node (2)",
			user: user2,
			id:   node.Id,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				checkErrorBody(t, recorder, e.ErrDeleteNodeNotPermitted)
			},
		},
		{
			name: "ok",
			user: user,
			id:   node.Id,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "not found node",
			user: user,
			id:   node.Id,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				checkErrorBody(t, recorder, e.ErrNodeIdNotFound)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/node/"+strconv.Itoa(tc.id), nil)
			setAuth(req, tc.user.Username, tc.user.Password)
			log.Println(req.Header)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}

}

func TestListNode(t *testing.T) {
	user, _, _ := autoInsertNode(nil, false)

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
			setAuth(req, tc.user.Username, tc.user.Password)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}

func TestUpdateNode(t *testing.T) {
	user, hardware, node := autoInsertNode(nil, false)
	node2 := randomNode()
	user2 := randomUser()
	user2.Status = true
	user2.Id = insertUser(user2)

	admin := randomUser()
	admin.Status = true
	admin.Is_admin = true
	admin.Id = insertUser(admin)

	sensor := autoInsertSensor()
	field := "mm"

	testCases := []struct {
		name          string
		id            int
		user          testUser
		body          gin.H
		checkResponse func(recorder *httptest.ResponseRecorder)
		checkInDB     func(id int)
	}{
		{
			name: "ok update name",
			id:   node.Id,
			user: user,
			body: gin.H{
				"name": node2.Name,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
			checkInDB: func(id int) {
				n := models.GetNodeByHardwareId(hardware.Id)
				require.Equal(t, node2.Name, n.Name)
			},
		},
		{
			name: "ok update location",
			id:   node.Id,
			user: user,
			body: gin.H{
				"location": node2.Location,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
			checkInDB: func(id int) {
				n := models.GetNodeByHardwareId(hardware.Id)
				require.Equal(t, node2.Location, n.Location)
			},
		},
		{
			name: "ok update field and sensor",
			id:   node.Id,
			user: user,
			body: gin.H{
				"field_sensor":       []*string{nil, nil, &field, nil, nil, nil, nil, nil, nil, nil},
				"id_hardware_sensor": []*int{nil, nil, &sensor.Id, nil, nil, nil, nil, nil, nil, nil},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
			checkInDB: func(id int) {
				n := models.GetNodeByHardwareId(hardware.Id)
				require.Equal(t, node2.Location, n.Location)
			},
		},
		{
			name: "ok all nil",
			id:   node.Id,
			user: user,
			body: gin.H{
				"field_sensor":       []*string{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
				"id_hardware_sensor": []*int{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
			checkInDB: func(id int) {
				n := models.GetNodeByHardwareId(hardware.Id)
				require.Equal(t, node2.Location, n.Location)
			},
		},
		{
			name: "field is empty",
			id:   node.Id,
			user: user,
			body: gin.H{
				"field_sensor":       []*string{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
				"id_hardware_sensor": []*int{nil, nil, &sensor.Id, nil, nil, nil, nil, nil, nil, nil},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkErrorBody(t, recorder, e.ErrFieldIsEmpty)
			},
			checkInDB: func(id int) {

			},
		},
		{
			name: "field is empty 2",
			id:   node.Id,
			user: user,
			body: gin.H{
				"id_hardware_sensor": []*int{&sensor.Id, nil, nil, nil, nil, nil, nil, nil, nil, nil},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkErrorBody(t, recorder, e.ErrFieldIsEmpty)
			},
			checkInDB: func(id int) {

			},
		},
		{
			name: "sensor is empty",
			id:   node.Id,
			user: user,
			body: gin.H{
				"field_sensor":       []*string{nil, nil, nil, &field, nil, nil, nil, nil, nil, nil},
				"id_hardware_sensor": []*int{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkErrorBody(t, recorder, e.ErrIdSensorIsEmpty)
			},
			checkInDB: func(id int) {

			},
		},
		{
			name: "sensor is empty 2",
			id:   node.Id,
			user: user,
			body: gin.H{
				"field_sensor": []*string{nil, nil, nil, nil, nil, nil, nil, nil, &field, nil},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkErrorBody(t, recorder, e.ErrIdSensorIsEmpty)
			},
			checkInDB: func(id int) {

			},
		},
		{
			name: "hardware is not sensor",
			id:   node.Id,
			user: user,
			body: gin.H{
				"field_sensor":       []*string{&field, nil, nil, nil, nil, nil, nil, nil, nil, nil},
				"id_hardware_sensor": []*int{&hardware.Id, nil, nil, nil, nil, nil, nil, nil, nil, nil},
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkErrorBody(t, recorder, e.ErrHardwareMustbeSensor)
			},
			checkInDB: func(id int) {

			},
		},
		{
			name: "using another user",
			id:   node.Id,
			user: user2,
			body: gin.H{
				"location": node2.Location,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				checkErrorBody(t, recorder, e.ErrEditNodeNotPermitted)
			},
			checkInDB: func(id int) {},
		},
		{
			name: "node is not exists",
			id:   1337,
			user: user,
			body: gin.H{
				"location": node2.Location,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				checkErrorBody(t, recorder, e.ErrNodeIdNotFound)
			},
			checkInDB: func(id int) {},
		},
		{
			name: "using admin",
			id:   node.Id,
			user: admin,
			body: gin.H{
				"location": node.Location,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
			checkInDB: func(id int) {
				n := models.GetNodeByHardwareId(hardware.Id)
				require.Equal(t, node.Location, n.Location)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			data, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("PUT", "/node/"+strconv.Itoa(tc.id), bytes.NewBuffer(data))
			setAuth(req, tc.user.Username, tc.user.Password)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
			tc.checkInDB(tc.id)
		})
	}

}
