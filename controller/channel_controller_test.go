package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	e "src/error"
	"src/models"
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
// 	e "src/error"
// 	"src/models"
// 	"testing"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/require"
// )

func randomChannel() models.Channel {
	var c models.Channel
	c.Value = make([]*float64, 10)
	tmp := make([]float64, 10)
	for i := 0; i < 10; i++ {
		tmp[i] = rand.Float64()
		c.Value[i] = &tmp[i]
	}
	return c
}

// func insertChannel(c models.Channel) {
// 	c.Time = time.Now()
// 	statement := "insert into channel (time, value, id_sensor) values (($1), $2, $3)"
// 	_, err := db.Exec(context.Background(), statement, c.Time, c.Value, c.Id_node)
// 	e.PanicIfNeeded(err)
// }

func TestAddChannel(t *testing.T) {

	channel := randomChannel()
	user, _, node := autoInsertNode(nil, false)
	fmt.Println(node.Id)

	user2 := randomUser()
	user2.Status = true
	user2.Id = insertUser(user2)

	useradmin := randomUser()
	useradmin.Status = true
	useradmin.Is_admin = true
	useradmin.Id = insertUser(useradmin)

	testCases := []struct {
		name          string
		body          gin.H
		user          testUser
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			body: gin.H{
				"value":   []*float64{channel.Value[0], channel.Value[1], nil, nil, nil, nil, nil, nil, nil, nil},
				"id_node": node.Id,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "field is empty",
			body: gin.H{
				"value":   []*float64{channel.Value[0], channel.Value[1], channel.Value[2], nil, nil, nil, nil, nil, nil, nil},
				"id_node": node.Id,
			},
			user: user,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkErrorBody(t, recorder, e.ErrFieldIsEmpty)
			},
		},
		{
			name: "using another user",
			body: gin.H{
				"value":   []*float64{channel.Value[0], channel.Value[1], nil, nil, nil, nil, nil, nil, nil, nil},
				"id_node": node.Id,
			},
			user: user2,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
				checkErrorBody(t, recorder, e.ErrUseNodeNotPermitted)
			},
		},
		{
			name: "ok using admin",
			body: gin.H{
				"value":   []*float64{channel.Value[0], channel.Value[1], nil, nil, nil, nil, nil, nil, nil, nil},
				"id_node": node.Id,
			},
			user: useradmin,
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
			setAuth(req, tc.user.Username, tc.user.Password)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}
