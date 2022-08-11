package repository

import (
	"database/sql"
	"fmt"
)

type userRepository struct {
	Conn *sql.DB
}

func NewUserRepository(database *sql.DB) UserRepository {
	return &userRepository{
		Conn: database,
	}
}

func (u *userRepository) Create(user User) {
	// statement := "insert into users (username) values ($1) returning user_id"
	// stmt, err := u.Conn.Prepare(statement)
	// if err != nil {
	// 	return
	// }
	// defer stmt.Close()
	// err = stmt.QueryRow(user.Username).Scan(&user.User_id)
	// fmt.Println(user.User_id)
	statement := "insert into user_person (username, email, password, token) values ($1, $2, $3, $4)"
	stmt, err := u.Conn.Prepare(statement)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()
	stmt.QueryRow(user.Username, user.Email, user.Password, user.Token)
	// return User{}
}

func (u *userRepository) GetByUsername(userName string) User {
	statement := "select * from user_person where username='$1'"
	stmt, err := u.Conn.Prepare(statement)
	if err != nil {
		fmt.Println(err)
		return User{}
	}
	defer stmt.Close()
	var res User
	err = stmt.QueryRow(userName).Scan(&res)
	fmt.Println(err)
	// fmt.Println(res)
	return res
}

func (u *userRepository) GetPassword(hashedPass string) string {
	return "ss"
}

func (u *userRepository) GetStatusByUsername(userName string) bool {
	return true
}
