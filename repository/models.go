package repository

import (
	"src/config"
	e "src/error"

	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Setup(cfg config.Configuration) *sql.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBPort)
	// db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db, _ = sql.Open("postgres", dsn)
	// fmt.Println(db)
	return db
}

func NewUserRepository(database *sql.DB) UserRepo {
	return &userRepoImpl{
		Db: database,
	}
}

func isRowExist(query string, args ...interface{}) bool {
	var exist bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := db.QueryRow(query, args...).Scan(&exist)
	if err != nil && err != sql.ErrNoRows {
		e.PanicIfNeeded(err)
	}
	return exist
}
