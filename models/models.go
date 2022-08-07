package models

import (
	"fmt"
	// "gorm.io/driver/postgres"
	// "gorm.io/gorm"
	"database/sql"
	"src/cfg"

	_ "github.com/lib/pq"
)

type User struct {
	User_id  int
	Username string
}

var DB *sql.DB

func Setup() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.App.DBHost,
		cfg.App.DBUser,
		cfg.App.DBPass,
		cfg.App.DBName,
		cfg.App.DBPort)
	fmt.Println(cfg.App)
	// db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB, _ = sql.Open("postgres", dsn)
	fmt.Println(DB)
}

func Cx() {
	user := User{Username: "Bintang"}
	fmt.Println(user.Create())
}

func (user User) Create() (err error) {
	statement := "insert into users (username) values ($1) returning user_id"
	fmt.Println(DB)
	stmt, err := DB.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(user.Username).Scan(&user.User_id)
	fmt.Println(user.User_id)
	return
}
