package models

import (
	"database/sql"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Status   bool   `json:"status"`
	Token    string `json:"token"`
}

func AddUser(data map[string]interface{}) error {
	user := User{
		Email:    data["email"].(string),
		Username: data["username"].(string),
		Password: data["password"].(string),
		Status:   data["status"].(bool),
		Token:    data["token"].(string),
	}

	statement := "insert into user_person (username, email, password, token) values ($1, $2, $3, $4)"
	_, err := db.Exec(statement, user.Username, user.Email, user.Password, user.Token)
	return err
}

func IsUserUsernameExist(username string) (bool, error) {
	statement := "select username from user_person where username = $1"

	// **to do** gimana cara gak pakai var tmp
	var tmp string
	if err := db.QueryRow(statement, username).Scan(&tmp); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil

}

// func GetByUsername(userName string) User {
// 	statement := "select * from user_person where username='$1'"
// 	stmt, err := u.Conn.Prepare(statement)
// 	if err != nil {
// 		fmt.Println(err)
// 		return User{}
// 	}
// 	defer stmt.Close()
// 	var res User
// 	err = stmt.QueryRow(userName).Scan(&res)
// 	fmt.Println(err)
// 	// fmt.Println(res)
// 	return res
// }

// func GetPassword(hashedPass string) string {
// 	return "ss"
// }

// func GetStatusByUsername(userName string) bool {
// 	return true
// }
