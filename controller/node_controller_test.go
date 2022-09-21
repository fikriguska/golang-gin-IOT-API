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

	"log"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func randomNode() models.Node {
	return models.Node{
		Name:     util.RandomString(10),
		Location: util.RandomString(7),
	}
}

func insertNode(n models.Node) int {
	statement := "insert into node (name, location, id_user, id_hardware) values ($1, $2, $3, $4) returning id_node"
	var id int
	err := db.QueryRow(statement, n.Name, n.Location, n.Id_user, n.Id_hardware).Scan(&id)
	e.PanicIfNeeded(err)
	return id
}

// create for user, hardware, node automatically, and insert it to db
func autoInsertNode(hardwareType interface{}) (testUser, models.Hardware, models.Node) {
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

	node.Id_hardware = hardwareNode.Id
	node.Id_user = user.Id
	node.Id = insertNode(node)
	return user, hardwareNode, node
}

func checkNodeSensor(node models.NodeGet, sensor models.Sensor) bool {
	containsSensor := false
	for _, v := range node.Sensor {
		if v.Id_sensor == sensor.Id {
			containsSensor = true
			break
		}
	}
	return containsSensor
}

func checkNodeHardware(node models.NodeGet, hardware models.Hardware) bool {
	containsHardware := false
	for _, v := range node.Hardware {
		if v.Name == hardware.Name && v.Type == hardware.Type {
			containsHardware = true
			break
		}
	}
	return containsHardware
}

func TestAddNode(t *testing.T) {
	node := randomNode()
	hardware := randomHardwareNode()
	id_hardware := insertHardware(hardware)

	user := randomUser()
	user.Status = true
	insertUser(user)
	testCases := []struct {
		name          string
		body          gin.H
		user          testUser
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ok no hardware",
			body: gin.H{
				"name":     node.Name,
				"location": node.Location,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "ok with hardware",
			body: gin.H{
				"name":        node.Name,
				"location":    node.Location,
				"id_hardware": id_hardware,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "not found hardware",
			body: gin.H{
				"name":        node.Name,
				"location":    node.Location,
				"id_hardware": 1337,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				checkErrorBody(t, recorder, e.ErrHardwareNotFound)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			data, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("POST", "/node/", bytes.NewBuffer(data))
			req.SetBasicAuth(tc.user.Username, tc.user.Password)
			log.Println(req.Header)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}

}

func TestGetNode(t *testing.T) {
	user, hardware, node, _, sensor := autoInsertSensor()

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
				require.Equal(t, n.Id, node.Id)

				require.Equal(t, checkNodeSensor(n, sensor), true)
				require.Equal(t, checkNodeHardware(n, hardware), true)
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
				require.Equal(t, node.Id, n.Id)
				require.Equal(t, true, checkNodeSensor(n, sensor))
				require.Equal(t, true, checkNodeHardware(n, hardware))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/node/"+strconv.Itoa(tc.id), nil)
			req.SetBasicAuth(tc.user.Username, tc.user.Password)
			log.Println(req.Header)
			router.ServeHTTP(w, req)
			// log.Println(w.Body)
			// log.Println(hardware)
			// log.Println(sensor)
			// log.Println(node)
			tc.checkResponse(w)

		})
	}
}

func TestDeleteNode(t *testing.T) {

	// create node user 1
	node := randomNode()
	hardware := randomHardwareNode()
	id_hardware := insertHardware(hardware)

	user := randomUser()
	user.Status = true

	id_user := insertUser(user)

	node.Id_hardware = id_hardware
	node.Id_user = id_user
	id_node := insertNode(node)
	log.Println(id_node)
	log.Println(node)

	// create node user 2
	node2 := randomNode()
	hardware2 := randomHardwareNode()
	id_hardware2 := insertHardware(hardware2)

	user2 := randomUser()
	user2.Status = true

	id_user2 := insertUser(user2)
	node2.Id_hardware = id_hardware2
	node2.Id_user = id_user2

	id_node2 := insertNode(node2)
	log.Println(id_node2)
	log.Println(node2)

	testCases := []struct {
		name          string
		user          testUser
		id            int
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "forbidden to delete another user's node (1)",
			user: user,
			id:   id_node2,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				checkErrorBody(t, recorder, e.ErrDeleteNodeNotPermitted)
			},
		},
		{
			name: "forbidden to delete another user's node (2)",
			user: user2,
			id:   id_node,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				checkErrorBody(t, recorder, e.ErrDeleteNodeNotPermitted)
			},
		},
		{
			name: "ok",
			user: user,
			id:   id_node,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "not found node",
			user: user,
			id:   id_node,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				checkErrorBody(t, recorder, e.ErrNodeNotFound)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", "/node/"+strconv.Itoa(tc.id), nil)
			req.SetBasicAuth(tc.user.Username, tc.user.Password)
			log.Println(req.Header)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}

}

func TestListNode(t *testing.T) {
	user, _, _ := autoInsertNode(nil)

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
			req, _ := http.NewRequest("GET", "/node/", nil)
			req.SetBasicAuth(tc.user.Username, tc.user.Password)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}
