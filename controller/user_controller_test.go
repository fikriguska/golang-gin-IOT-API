package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	// "net/http/httptest"

	"src/config"
	"src/models"
	"src/util"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var router *gin.Engine

func SetupRouter() *gin.Engine {
	r := gin.Default()
	UserRoute(r)
	return r
}

func randomUser() AddUserStruct {
	password := util.RandomString(12)
	email := util.RandomEmail()
	username := util.RandomString(6)

	return AddUserStruct{
		Password: password,
		Email:    email,
		Username: username,
	}

}

func TestAddUser(t *testing.T) {

	user1 := randomUser()
	user2 := randomUser()

	testCases := []struct {
		name          string
		body          gin.H
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			body: gin.H{
				"username": user1.Username,
				"email":    user1.Email,
				"password": user1.Password,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				fmt.Printf("code ===> %d", recorder.Code)
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "username exists",
			body: gin.H{
				"username": user1.Username,
				"email":    user2.Email,
				"password": user2.Password,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				fmt.Printf("code ===> %d", recorder.Code)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "email exists",
			body: gin.H{
				"username": user2.Username,
				"email":    user1.Email,
				"password": user2.Password,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				fmt.Printf("code ===> %d", recorder.Code)
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			data, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(data))
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}

}

func TestMain(m *testing.M) {
	cfg := config.Setup()
	models.Setup(cfg)
	router = SetupRouter()
	m.Run()
}
