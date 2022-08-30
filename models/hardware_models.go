package models

import (
	e "src/error"
)

type Hardware struct {
	Id          int
	Name        string
	Type        string
	Description string
}

func AddHardware(h Hardware) {
	statement := "insert into hardware (name, type, description) values ($1, $2, $3)"
	_, err := db.Exec(statement, h.Name, h.Type, h.Description)
	e.PanicIfNeeded(err)
}

func DeleteHardware(id int) {
	statement := "delete from hardware where id_hardware = $1"
	_, err := db.Exec(statement, id)
	e.PanicIfNeeded(err)
}
