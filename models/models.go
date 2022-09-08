package models

import (
	"encoding/json"
	"reflect"
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
	fmt.Println(db)
	return db
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

// func getRows(query string, model interface{}, args ...interface{}) {

// }
type NullInt64 sql.NullInt64

func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

func (ni *NullInt64) Scan(value interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ni = NullInt64{i.Int64, false}
	} else {
		*ni = NullInt64{i.Int64, true}
	}
	return nil
}
