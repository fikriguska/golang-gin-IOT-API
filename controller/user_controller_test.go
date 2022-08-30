package controller

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	// "net/http/httptest"

	"src/config"
	e "src/error"
	"src/models"
	"src/util"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var router *gin.Engine
var db *sql.DB

type testUser struct {
	models.User
	hashedPass string
}

func SetupRouter() *gin.Engine {
	r := gin.Default()
	UserRoute(r)
	return r
}

func randomUser() testUser {
	password := util.RandomString(12)
	email := util.RandomEmail()
	username := util.RandomString(6)

	hashedTokenByte := sha256.Sum256([]byte(username + email + password))
	hashedToken := hex.EncodeToString(hashedTokenByte[:])
	hashedPassByte := sha256.Sum256([]byte(password))
	hashedPass := hex.EncodeToString(hashedPassByte[:])

	return testUser{
		User: models.User{
			Password: password,
			Email:    email,
			Username: username,
			Status:   false,
			Token:    hashedToken,
			Is_admin: false,
		},
		hashedPass: hashedPass,
	}
}

func insertUser(u testUser) {
	_, err := db.Exec("insert into user_person (username, email, password, status, token, is_admin) values ($1, $2, $3, $4, $5, $6)", u.Username, u.Email, u.hashedPass, u.Status, u.Token, u.Is_admin)
	e.PanicIfNeeded(err)
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
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkBody(t, recorder, e.ErrUserExist)
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
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkBody(t, recorder, e.ErrUserExist)
			},
		},
		{
			name: "invalid email format",
			body: gin.H{
				"username": user1.Username,
				"email":    "thisisnotan.email",
				"password": user1.Password,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkBody(t, recorder, e.ErrInvalidEmail)
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

func TestActivateUser(t *testing.T) {
	user := randomUser()

	insertUser(user)

	testCases := []struct {
		name          string
		token         string
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name:  "ok",
			token: user.Token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "user is activated",
			token: user.Token,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkBody(t, recorder, e.ErrInvalidToken)
			},
		},
		{
			name:  "token doensn't exists",
			token: "invalidtoken1337133713371337133713371337133713371337133713371337",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkBody(t, recorder, e.ErrInvalidToken)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/user/activation?token="+tc.token, nil)
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}

func TestUserLogin(t *testing.T) {

	user := randomUser()
	user2 := randomUser()

	user.Status = true
	insertUser(user)
	insertUser(user2)

	fmt.Println(user)
	testCases := []struct {
		name          string
		body          gin.H
		checkResponse func(recoder *httptest.ResponseRecorder)
	}{
		{
			name: "ok",
			body: gin.H{
				"username": user.Username,
				"password": user.Password,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "wrong password",
			body: gin.H{
				"username": user.Username,
				"password": "wrong_password",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkBody(t, recorder, e.ErrUsernameOrPassIncorrect)
			},
		},
		{
			name: "wrong username",
			body: gin.H{
				"username": "wrong_username",
				"password": user.Password,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkBody(t, recorder, e.ErrUsernameOrPassIncorrect)
			},
		},
		{
			name: "user is not activated",
			body: gin.H{
				"username": user2.Username,
				"password": user2.Password,
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				checkBody(t, recorder, e.ErrUserNotActive)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			data, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("POST", "/user/login", bytes.NewBuffer(data))
			router.ServeHTTP(w, req)
			tc.checkResponse(w)
		})
	}
}

func TestMain(m *testing.M) {
	cfg := config.Setup()
	db = models.Setup(cfg)
	router = SetupRouter()
	m.Run()
}
