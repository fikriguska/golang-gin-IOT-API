package controller

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"strconv"

// 	// "net/http/httptest"

// 	// "net/http/httptest"

// 	e "src/error"
// 	"src/models"
// 	"src/util"

// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/require"
// )

// type testUser struct {
// 	models.User
// 	hashedPass string
// }

// func randomUser() testUser {
// 	password := util.RandomString(12)
// 	email := util.RandomEmail()
// 	username := util.RandomString(6)

// 	hashedToken := util.Sha256String(username + email + password)
// 	hashedPass := util.Sha256String(password)

// 	return testUser{
// 		User: models.User{
// 			Password: password,
// 			Email:    email,
// 			Username: username,
// 			Status:   false,
// 			Token:    hashedToken,
// 			Is_admin: false,
// 		},
// 		hashedPass: hashedPass,
// 	}
// }

// func insertUser(u testUser) int {
// 	statement := "insert into user_person (username, email, password, status, token, isadmin) values ($1, $2, $3, $4, $5, $6) returning id_user"
// 	var id int
// 	err := db.QueryRow(context.Background(), statement, u.Username, u.Email, u.hashedPass, u.Status, u.Token, u.Is_admin).Scan(&id)
// 	e.PanicIfNeeded(err)
// 	return id
// }

// func TestAddUser(t *testing.T) {

// 	user1 := randomUser()
// 	user2 := randomUser()

// 	testCases := []struct {
// 		name          string
// 		body          gin.H
// 		checkResponse func(recoder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "ok",
// 			body: gin.H{
// 				"username": user1.Username,
// 				"email":    user1.Email,
// 				"password": user1.Password,
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusCreated, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "username exists",
// 			body: gin.H{
// 				"username": user1.Username,
// 				"email":    user2.Email,
// 				"password": user2.Password,
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 				checkErrorBody(t, recorder, e.ErrEmailUsernameAlreadyUsed)
// 			},
// 		},
// 		{
// 			name: "email exists",
// 			body: gin.H{
// 				"username": user2.Username,
// 				"email":    user1.Email,
// 				"password": user2.Password,
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 				checkErrorBody(t, recorder, e.ErrEmailUsernameAlreadyUsed)
// 			},
// 		},
// 		{
// 			name: "invalid email format",
// 			body: gin.H{
// 				"username": user1.Username,
// 				"email":    "thisisnotan.email",
// 				"password": user1.Password,
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 				checkErrorBody(t, recorder, e.ErrInvalidEmail)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			w := httptest.NewRecorder()
// 			data, _ := json.Marshal(tc.body)
// 			req, _ := http.NewRequest("POST", "/user/signup", bytes.NewBuffer(data))
// 			router.ServeHTTP(w, req)
// 			tc.checkResponse(w)
// 		})
// 	}

// }

// func TestActivateUser(t *testing.T) {
// 	user := randomUser()

// 	insertUser(user)

// 	testCases := []struct {
// 		name          string
// 		token         string
// 		checkResponse func(recoder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name:  "ok",
// 			token: user.Token,
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name:  "user is activated",
// 			token: user.Token,
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 				checkErrorBody(t, recorder, e.ErrUserAlreadyActive)
// 			},
// 		},
// 		{
// 			name:  "token doensn't exists",
// 			token: "invalidtoken1337133713371337133713371337133713371337133713371337",
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusNotFound, recorder.Code)
// 				checkErrorBody(t, recorder, e.ErrTokenNotFound)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			w := httptest.NewRecorder()
// 			req, _ := http.NewRequest("GET", "/user/activation?token="+tc.token, nil)
// 			router.ServeHTTP(w, req)
// 			tc.checkResponse(w)
// 		})
// 	}
// }

// func TestUserLogin(t *testing.T) {

// 	user := randomUser()
// 	user2 := randomUser()

// 	user.Status = true
// 	insertUser(user)
// 	insertUser(user2)

// 	fmt.Println(user)
// 	testCases := []struct {
// 		name          string
// 		body          gin.H
// 		checkResponse func(recoder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "ok",
// 			body: gin.H{
// 				"username": user.Username,
// 				"password": user.Password,
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "wrong password",
// 			body: gin.H{
// 				"username": user.Username,
// 				"password": "wrong_password",
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 				checkErrorBody(t, recorder, e.ErrUsernameOrPassIncorrect)
// 			},
// 		},
// 		{
// 			name: "wrong username",
// 			body: gin.H{
// 				"username": "wrong_username",
// 				"password": user.Password,
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 				checkErrorBody(t, recorder, e.ErrUsernameOrPassIncorrect)
// 			},
// 		},
// 		{
// 			name: "user is not activated",
// 			body: gin.H{
// 				"username": user2.Username,
// 				"password": user2.Password,
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 				checkErrorBody(t, recorder, e.ErrUserNotActive)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			w := httptest.NewRecorder()
// 			data, _ := json.Marshal(tc.body)
// 			req, _ := http.NewRequest("POST", "/user/login", bytes.NewBuffer(data))
// 			router.ServeHTTP(w, req)
// 			tc.checkResponse(w)
// 		})
// 	}
// }

// func TestUpdateUser(t *testing.T) {
// 	user := randomUser()
// 	user.Status = true
// 	user.Id = insertUser(user)

// 	user_newpass := "maguire"
// 	user_hashednewpass := util.Sha256String(user_newpass)

// 	user2 := randomUser()
// 	user2.Status = true
// 	user2.Id = insertUser(user2)

// 	testCases := []struct {
// 		name          string
// 		id            int
// 		user          testUser
// 		body          gin.H
// 		checkResponse func(recorder *httptest.ResponseRecorder)
// 		checkInDB     func(id int)
// 	}{
// 		{
// 			name: "wrong old pass",
// 			id:   user.Id,
// 			user: user,
// 			body: gin.H{
// 				"old password": "wr0ngP4ss",
// 				"new password": user_newpass,
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 				checkErrorBody(t, recorder, e.ErrOldPasswordIncorrect)
// 			},
// 			checkInDB: func(id int) {},
// 		},
// 		{
// 			name: "using different user",
// 			id:   user.Id,
// 			user: user2,
// 			body: gin.H{
// 				"old password": user.Password,
// 				"new password": user_newpass,
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusForbidden, recorder.Code)
// 				checkErrorBody(t, recorder, e.ErrEditUserNotPermitted)
// 			},
// 			checkInDB: func(id int) {},
// 		},
// 		{
// 			name: "user not found",
// 			id:   31337,
// 			user: user,
// 			body: gin.H{
// 				"old password": user.Password,
// 				"new password": user_newpass,
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusNotFound, recorder.Code)
// 				checkErrorBody(t, recorder, e.ErrUserIdNotFound)
// 			},
// 			checkInDB: func(id int) {},
// 		},
// 		{
// 			name: "ok",
// 			id:   user.Id,
// 			user: user,
// 			body: gin.H{
// 				"old password": user.Password,
// 				"new password": user_newpass,
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 			checkInDB: func(id int) {
// 				u := models.GetUserById(id)
// 				require.Equal(t, user_hashednewpass, u.Password)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			w := httptest.NewRecorder()
// 			data, _ := json.Marshal(tc.body)
// 			req, _ := http.NewRequest("PUT", "/user/"+strconv.Itoa(tc.id), bytes.NewBuffer(data))
// 			setAuth(req, tc.user.Username, tc.user.Password)
// 			router.ServeHTTP(w, req)
// 			tc.checkResponse(w)
// 			tc.checkInDB(tc.id)
// 		})
// 	}
// }

// func TestForgetPassword(t *testing.T) {
// 	user := randomUser()
// 	user.Status = true
// 	user.Id = insertUser(user)

// 	user2 := randomUser()
// 	user2.Id = insertUser(user2)

// 	testCases := []struct {
// 		name          string
// 		body          gin.H
// 		checkResponse func(recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "ok",
// 			body: gin.H{
// 				"email":    user.Email,
// 				"username": user.Username,
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},

// 		{
// 			name: "email and username is not matched (1)",
// 			body: gin.H{
// 				"email":    user.Email,
// 				"username": "princess_mononoke",
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 				checkErrorBody(t, recorder, e.ErrUsernameOrEmailIncorrect)
// 			},
// 		},

// 		{
// 			name: "email and username is not matched (2)",
// 			body: gin.H{
// 				"email":    "mononoke.princess@ghibli.com",
// 				"username": user.Username,
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 				checkErrorBody(t, recorder, e.ErrUsernameOrEmailIncorrect)
// 			},
// 		},

// 		{
// 			name: "invalid email format",
// 			body: gin.H{
// 				"email":    "this.is.not.an.email",
// 				"username": user.Username,
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 				checkErrorBody(t, recorder, e.ErrInvalidEmail)
// 			},
// 		},

// 		{
// 			name: "user is not activated",
// 			body: gin.H{
// 				"email":    user2.Email,
// 				"username": user2.Username,
// 			},
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 				checkErrorBody(t, recorder, e.ErrUserNotActive)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			w := httptest.NewRecorder()
// 			data, _ := json.Marshal(tc.body)
// 			req, _ := http.NewRequest("POST", "/user/forget-password", bytes.NewBuffer(data))
// 			router.ServeHTTP(w, req)
// 			tc.checkResponse(w)
// 		})
// 	}
// }

// func TestDeleteUser(t *testing.T) {
// 	user := randomUser()
// 	user.Status = true
// 	user.Id = insertUser(user)

// 	user2, _, _ := autoInsertNode(nil)

// 	user3 := randomUser()
// 	user3.Status = true
// 	user3.Id = insertUser(user3)

// 	admin := randomUser()
// 	admin.Status = true
// 	admin.Is_admin = true
// 	admin.Id = insertUser(admin)

// 	testCases := []struct {
// 		name          string
// 		id            int
// 		user          testUser
// 		checkResponse func(recorder *httptest.ResponseRecorder)
// 	}{
// 		{
// 			name: "using another user",
// 			id:   user.Id,
// 			user: user2,
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusForbidden, recorder.Code)
// 				checkErrorBody(t, recorder, e.ErrDeleteUserNotPermitted)
// 			},
// 		},
// 		{
// 			name: "user is not exists",
// 			id:   31337,
// 			user: admin,
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusNotFound, recorder.Code)
// 				checkErrorBody(t, recorder, e.ErrUserIdNotFound)
// 			},
// 		},
// 		{
// 			name: "user is still using node",
// 			id:   user2.Id,
// 			user: user2,
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusBadRequest, recorder.Code)
// 				checkErrorBody(t, recorder, e.ErrUserStillUsingNode)
// 			},
// 		},
// 		{
// 			name: "ok",
// 			id:   user.Id,
// 			user: user,
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 		{
// 			name: "ok using admin",
// 			id:   user3.Id,
// 			user: admin,
// 			checkResponse: func(recorder *httptest.ResponseRecorder) {
// 				require.Equal(t, http.StatusOK, recorder.Code)
// 			},
// 		},
// 	}

// 	for i := range testCases {
// 		tc := testCases[i]

// 		t.Run(tc.name, func(t *testing.T) {
// 			w := httptest.NewRecorder()
// 			req, _ := http.NewRequest("DELETE", "/user/"+strconv.Itoa(tc.id), nil)
// 			setAuth(req, tc.user.Username, tc.user.Password)
// 			router.ServeHTTP(w, req)
// 			tc.checkResponse(w)
// 		})
// 	}
// }
