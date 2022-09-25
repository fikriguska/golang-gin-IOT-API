package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	// "net/http/httptest"

	e "src/error"
	"src/models"
	"src/util"

	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type testUser struct {
	models.User
	hashedPass string
}

func randomUser() testUser {
	password := util.RandomString(12)
	email := util.RandomEmail()
	username := util.RandomString(6)

	hashedToken := util.Sha256String(username + email + password)
	hashedPass := util.Sha256String(password)

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

func insertUser(u testUser) int {
	statement := "insert into user_person (username, email, password, status, token, is_admin) values ($1, $2, $3, $4, $5, $6) returning id_user"
	var id int
	err := db.QueryRow(statement, u.Username, u.Email, u.hashedPass, u.Status, u.Token, u.Is_admin).Scan(&id)
	e.PanicIfNeeded(err)
	return id
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
				checkErrorBody(t, recorder, e.ErrEmailUsernameAlreadyUsed)
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
				checkErrorBody(t, recorder, e.ErrEmailUsernameAlreadyUsed)
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
				checkErrorBody(t, recorder, e.ErrInvalidEmail)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			data, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("POST", "/user/signup", bytes.NewBuffer(data))
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
				checkErrorBody(t, recorder, e.ErrUserAlreadyActive)
			},
		},
		{
			name:  "token doensn't exists",
			token: "invalidtoken1337133713371337133713371337133713371337133713371337",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				checkErrorBody(t, recorder, e.ErrTokenNotFound)
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
				checkErrorBody(t, recorder, e.ErrUsernameOrPassIncorrect)
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
				checkErrorBody(t, recorder, e.ErrUsernameOrPassIncorrect)
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
				checkErrorBody(t, recorder, e.ErrUserNotActive)
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
