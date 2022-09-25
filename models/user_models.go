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
	statement := "insert into user_person (username, email, password, status, token, is_admin) values ($1, $2, $3, $4, $5, $6)"
	_, err := db.Exec(statement, user.Username, user.Email, user.Password, user.Status, user.Token, user.Is_admin)
	e.PanicIfNeeded(err)
}

func GetUserByUsername(user User) User {
	statement := "select id_user, username, password, is_admin from user_person where username = $1"

	var u User
	err := db.QueryRow(statement, user.Username).Scan(&u.Id, &u.Username, &u.Password, &u.Is_admin)
	e.PanicIfNeeded(err)

	return u
}

func GetUserById(id int) User {
	statement := "select id_user, username, password, is_admin from user_person where id_user = $1"

	var u User
	err := db.QueryRow(statement, id).Scan(&u.Id, &u.Username, &u.Password, &u.Is_admin)
	e.PanicIfNeeded(err)

	return u
}

func IsUserUsernameExist(username string) bool {
	statement := "select username from user_person where username = $1"
	return isRowExist(statement, username)
}

func IsUserEmailExist(email string) bool {
	statement := "select email from user_person where email = $1"
	return isRowExist(statement, email)
}

func IsUserExistById(id int) bool {
	statement := "select id_user from user_person where id_user = $1"
	return isRowExist(statement, id)
}

// ******
func GetUserIdByUsername(username string) int {
	statement := "select id_user from user_person where username = $1"

	var id int
	err := db.QueryRow(statement, username).Scan(&id)
	e.PanicIfNeeded(err)

	return id
}

func IsUserTokenExist(token string) bool {
	statement := "select token from user_person where token = $1"
	return isRowExist(statement, token)
}

func IsUserActivatedCheckByToken(token string) bool {
	// ******
	statement := "select id_user from user_person where token = $1 and status = true"
	// var tmp string
	// if err := db.QueryRow(statement, token).Scan(&tmp); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return false, nil
	// 	}
	// 	return false, err
	// }
	return isRowExist(statement, token)
	// return true, nil
}

func IsUserActivatedCheckByUsername(username string) bool {
	statement := "select status from user_person where username = $1"

	var status bool
	err := db.QueryRow(statement, username).Scan(&status)
	if err == sql.ErrNoRows {
		return false
	}
	e.PanicIfNeeded(err)

	return status
}

func ActivateUser(token string) error {
	statement := "update user_person set status = true where token = $1"
	_, err := db.Exec(statement, token)
	e.PanicIfNeeded(err)
	return nil
}

func IsUsernameAndPasswordExist(username string, password string) bool {
	statement := "select id_user from user_person where username = $1 and password = $2"
	return isRowExist(statement, username, password)

}

func IsEmailAndUsernameExist(email string, username string) bool {
	statement := "select id_user from user_person where email = $1 and username = $2"
	return isRowExist(statement, email, username)
}

func UpdateUserPasswordByEmail(email string, password string) {
	statement := "update user_person set password = $1 where email = $2"
	_, err := db.Exec(statement, password, email)
	e.PanicIfNeeded(err)
}

func UpdateUserPasswordById(id int, password string) {
	statement := "update user_person set password = $1 where id_user = $2"
	_, err := db.Exec(statement, password, id)
	e.PanicIfNeeded(err)
}

func DeleteUser(id int) error {
	statement := "delete from user_person where id_user = $1"
	_, err := db.Exec(statement, id)
	// e.PanicIfNeeded(err)
	return err
}
