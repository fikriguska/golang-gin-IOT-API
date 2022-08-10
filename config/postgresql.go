package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func NewPostgresql(cfg Configuration) *sql.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBPort)
	// db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB, _ := sql.Open("postgres", dsn)
	fmt.Println(DB)
	return DB
}
