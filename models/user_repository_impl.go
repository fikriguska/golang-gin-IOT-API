package models

import (
	e "src/error"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Status   bool   `json:"status"`
	Token    string `json:"token"`
}

func AddUser(data map[string]interface{}) {
	user := User{
		Email:    data["email"].(string),
		Username: data["username"].(string),
		Password: data["password"].(string),
		Status:   data["status"].(bool),
		Token:    data["token"].(string),
	}

	statement := "insert into user_person (username, email, password, status, token) values ($1, $2, $3, $4, $5)"
	_, err := db.Exec(statement, user.Username, user.Email, user.Password, user.Status, user.Token)
	e.PanicIfNeeded(err)
}

func IsUserUsernameExist(username string) bool {
	statement := "select username from user_person where username = $1"

	// **to do** gimana cara gak pakai var tmp
	// var tmp string
	// if err := db.QueryRow(statement, username).Scan(&tmp); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return false, nil
	// 	}
	// 	return false, err
	// }
	// return true, nil
	return isRowExist(statement, username)

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

	// var tmp string
	// if err := db.QueryRow(statement, token).Scan(&tmp); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return false, nil
	// 	}
	// 	return false, err
	// }
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
	e.PanicIfNeeded(err)

	return status
}

func ActivateUser(token string) error {
	statement := "update user_person set status = true where token = $1"
	_, err := db.Exec(statement, token)
	e.PanicIfNeeded(err)
	return nil
}

func AuthUser(username string, password string) bool {
	statement := "select token from user_person where username = $1 and password = $2"
	// return isRowExist(statement, username, password)
	// var token string
	// err := db.QueryRow(statement, username, password).Scan(&token)
	// return token
	return isRowExist(statement, username, password)

}
