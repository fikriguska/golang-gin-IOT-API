package models

import (
	"src/config"

	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Setup(cfg config.Configuration) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBPort)
	// db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db, _ = sql.Open("postgres", dsn)
	fmt.Println(db)
}
