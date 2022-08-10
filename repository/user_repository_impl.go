package repository

import (
	"database/sql"
)

type userRepository struct {
	Conn *sql.DB
}

func NewUserRepository(database *sql.DB) UserRepository {
	return &userRepository{
		Conn: database,
	}
}

func (u *userRepository) Create(user User) User {
	// statement := "insert into users (username) values ($1) returning user_id"
	// stmt, err := u.Conn.Prepare(statement)
	// if err != nil {
	// 	return
	// }
	// defer stmt.Close()
	// err = stmt.QueryRow(user.Username).Scan(&user.User_id)
	// fmt.Println(user.User_id)
	return User{
		User_id:  123,
		Username: "ssss",
	}
}
