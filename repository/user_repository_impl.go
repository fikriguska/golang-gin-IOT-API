package repository

import (
	"database/sql"
	e "src/error"
	"src/model"
)

type userRepoImpl struct {
	Db *sql.DB
}

func (repository *userRepoImpl) Add(data map[string]interface{}) {
	user := model.User{
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

func (repository *userRepoImpl) IsUsernameExist(username string) bool {
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
func (repository *userRepoImpl) GetIdByUsername(username string) int {
	statement := "select id_user from user_person where username = $1"

	var id int
	err := db.QueryRow(statement, username).Scan(&id)
	e.PanicIfNeeded(err)

	return id
}

func (repository *userRepoImpl) IsTokenExist(token string) bool {
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

func (repository *userRepoImpl) IsActivatedCheckByToken(token string) bool {
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

func (repository *userRepoImpl) IsActivatedCheckByUsername(username string) bool {
	statement := "select status from user_person where username = $1"

	var status bool
	err := db.QueryRow(statement, username).Scan(&status)
	e.PanicIfNeeded(err)

	return status
}

func (repository *userRepoImpl) Activate(token string) {
	statement := "update user_person set status = true where token = $1"
	_, err := db.Exec(statement, token)
	e.PanicIfNeeded(err)
}

func (repository *userRepoImpl) Auth(username string, password string) bool {
	statement := "select token from user_person where username = $1 and password = $2"
	// return isRowExist(statement, username, password)
	// var token string
	// err := db.QueryRow(statement, username, password).Scan(&token)
	// return token
	return isRowExist(statement, username, password)

}
