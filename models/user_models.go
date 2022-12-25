package models

import (
	"database/sql"
	e "src/error"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Status   bool   `json:"status"`
	Token    string `json:"token"`
	Is_admin bool   `json:"is_admin"`
}

type UserAdd struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserGet struct {
	Id       int    `json:"id_user"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserList struct {
	Id       int    `json:"id_user"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserForgetPassword struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type UserUpdate struct {
	OldPasswd string `json:"old password" binding:"required"`
	NewPasswd string `json:"new password" binding:"required"`
}

func AddUser(user User) {
	statement := replaceQueryParam("insert into user_person (username, email, password, status, token, isadmin) values ('%s', '%s', '%s', %t, '%s', %t)", user.Username, user.Email, user.Password, user.Status, user.Token, user.Is_admin)
	_, err := db.Exec(statement)
	e.PanicIfNeeded(err)
}

func GetUserByUsername(user User) User {
	statement := replaceQueryParam("select id_user, username, password, isadmin from user_person where username = '%s'", user.Username)

	var u User
	err := db.QueryRow(statement).Scan(&u.Id, &u.Username, &u.Password, &u.Is_admin)
	e.PanicIfNeeded(err)

	return u
}

func GetUserById(id int) User {
	statement := replaceQueryParam("select id_user, email, username, password, isadmin from user_person where id_user = %d", id)

	var u User
	err := db.QueryRow(statement).Scan(&u.Id, &u.Email, &u.Username, &u.Password, &u.Is_admin)
	e.PanicIfNeeded(err)

	return u
}

func GetAllUser() []UserList {
	var user UserList
	var users []UserList
	users = make([]UserList, 0)
	statement := "select id_user, email, username from user_person"
	rows, err := db.Query(statement)
	e.PanicIfNeeded(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Email, &user.Username)
		e.PanicIfNeeded(err)
		users = append(users, user)
	}

	return users
}

func IsUserUsernameExist(username string) bool {
	statement := replaceQueryParam("select username from user_person where username = '%s'", username)
	return isRowExist(statement)
}

func IsUserEmailExist(email string) bool {
	statement := replaceQueryParam("select email from user_person where email = '%s'", email)
	return isRowExist(statement)
}

func IsUserExistById(id int) bool {
	statement := replaceQueryParam("select id_user from user_person where id_user = %d", id)
	return isRowExist(statement)
}

// ******
func GetUserIdByUsername(username string) int {
	statement := replaceQueryParam("select id_user from user_person where username = '%s'", username)

	var id int
	err := db.QueryRow(statement).Scan(&id)
	e.PanicIfNeeded(err)

	return id
}

func IsUserTokenExist(token string) bool {
	statement := replaceQueryParam("select token from user_person where token = '%s'", token)
	return isRowExist(statement)
}

func IsUserActivatedCheckByToken(token string) bool {
	// ******
	statement := replaceQueryParam("select id_user from user_person where token = '%s' and status = true", token)
	// var tmp string
	// if err := db.QueryRow(statement, token).Scan(&tmp); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return false, nil
	// 	}
	// 	return false, err
	// }
	return isRowExist(statement)
	// return true, nil
}

func IsUserActivatedCheckByUsername(username string) bool {
	statement := replaceQueryParam("select status from user_person where username = '%s'", username)

	var status bool
	err := db.QueryRow(statement).Scan(&status)
	if err == sql.ErrNoRows {
		return false
	}
	e.PanicIfNeeded(err)

	return status
}

func ActivateUser(token string) error {
	statement := replaceQueryParam("update user_person set status = true where token = '%s'", token)
	_, err := db.Exec(statement)
	e.PanicIfNeeded(err)
	return nil
}

func IsUsernameAndPasswordExist(username string, password string) bool {
	statement := replaceQueryParam("select id_user from user_person where username = '%s' and password = '%s'", username, password)
	return isRowExist(statement)

}

func IsEmailAndUsernameExist(email string, username string) bool {
	statement := replaceQueryParam("select id_user from user_person where email = '%s' and username = '%s'", email, username)
	return isRowExist(statement)
}

func UpdateUserPasswordByEmail(email string, password string) {
	statement := replaceQueryParam("update user_person set password = '%s' where email = '%s'", password, email)
	_, err := db.Exec(statement)
	e.PanicIfNeeded(err)
}

func UpdateUserPasswordById(id int, password string) {
	statement := replaceQueryParam("update user_person set password = '%s' where id_user = %d", password, id)
	_, err := db.Exec(statement)
	e.PanicIfNeeded(err)
}

func DeleteUser(id int) error {
	statement := replaceQueryParam("delete from user_person where id_user = %d", id)
	_, err := db.Exec(statement)
	// e.PanicIfNeeded(err)
	return err
}
