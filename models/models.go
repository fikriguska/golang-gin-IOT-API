package models

import (
	"context"
	"encoding/json"
	"reflect"
	"src/config"
	e "src/error"

	"database/sql"
	"database/sql/driver"
	"fmt"

	// _ "github.com/lib/pq"
	// _ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// var db *sql.DB
var db *pgxpool.Pool

func Setup(cfg config.Configuration) *pgxpool.Pool {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable pool_max_conns=8",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBPort)
	// db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// db, _ = sql.Open("pgx", dsn)
	db, _ = pgxpool.New(context.Background(), dsn)
	fmt.Println(db)
	return db
}

func isRowExist(query string, args ...interface{}) bool {
	var exist bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := db.QueryRow(cb(), query, args...).Scan(&exist)
	if err != nil && err != pgx.ErrNoRows {
		e.PanicIfNeeded(err)
	}
	return exist
}

func cb() context.Context {
	return context.Background()
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

type NullString sql.NullString

func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// func arrIntToPSQL(val []int) string {
// 	for i, v := range val {

// 	}
// }
