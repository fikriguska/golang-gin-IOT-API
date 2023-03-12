package models

import (
	"encoding/json"
	"reflect"
	"src/config"
	e "src/error"

	"database/sql"
	"database/sql/driver"
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
	db.SetMaxOpenConns(8)
	fmt.Println(db)
	return db
}

func isRowExist(query string) bool {
	var exist bool
	query = fmt.Sprintf("SELECT exists (%s)", query)
	err := db.QueryRow(query).Scan(&exist)
	if err != nil && err != sql.ErrNoRows {
		e.PanicIfNeeded(err)
	}
	return exist
}

// used to prevent unnecessary statement preparation that slow down the performance
func replaceQueryParam(query string, args ...interface{}) string {
	return fmt.Sprintf(query, args...)
}

// to support COALESCE
func fillByNullIfNeeded(args ...interface{}) {
	for _, arg := range args {
		var str string
		switch v := arg.(type) {
		case **string:
			if *v == nil {
				str = "NULL"
				*arg.(**string) = &str
			} else {
				str = fmt.Sprintf("'%s'", **v)
				*arg.(**string) = &str
			}
		}
	}
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
