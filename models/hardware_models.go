package models

import (
	"database/sql"
	e "src/error"
)

type Hardware struct {
	Id          int    `json:"id_hardware"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// GET /hardware
type HardwareList struct {
	Sensor []Hardware `json:"sensor"`
	Node   []Hardware `json:"node"`
}

// POST /hardware
type HardwareAdd struct {
	Name        string
	Type        string
	Description string
}

// GET /hardware/:id
type HardwareSensorGet struct {
	Hardware
	Sensor struct {
		Name string `json:"name"`
		Unit string `json:"unit"`
	} `json:"sensor"`
}

// GET /hardware/:id
type HardwareNodeGet struct {
	Hardware
	Node struct {
		Name     string `json:"name"`
		Location string `json:"location"`
	} `json:"node"`
}

type HardwareUpdate struct {
	Name        *string `json:"name"`
	Type        *string `json:"type"`
	Description *string `json:"description"`
}

type HardwareUpdateSQL struct {
	Name        NullString `json:"name"`
	Type        NullString `json:"type"`
	Description NullString `json:"description"`
}

func IsHardwareExistById(id int) bool {
	statement := "select id_hardware from hardware where id_hardware = $1"
	return isRowExist(statement, id)
}

func IsHardwareTypeSensorById(id int) bool {
	statement := "select type from hardware where id_hardware = $1 and (lower(type) = 'sensor')"
	return isRowExist(statement, id)
}

func IsHardwareTypeNodeById(id int) bool {
	statement := "select type from hardware where id_hardware = $1 and (lower(type) = 'single-board computer' or lower(type) = 'microcontroller unit')"
	return isRowExist(statement, id)
}

func AddHardware(h Hardware) {
	statement := "insert into hardware (name, type, description) values ($1, $2, $3)"
	_, err := db.Exec(statement, h.Name, h.Type, h.Description)
	e.PanicIfNeeded(err)
}

func GetHardwareById(id int) Hardware {
	var hardware Hardware
	statement := "select id_hardware, name, type, description from hardware where id_hardware = $1"
	err := db.QueryRow(statement, id).Scan(&hardware.Id, &hardware.Name, &hardware.Type, &hardware.Description)
	if err != nil && err != sql.ErrNoRows {
		e.PanicIfNeeded(err)
	}

	return hardware
}

func GetAllHardwareTypeSensor() []Hardware {
	var hardware Hardware
	var hardwares []Hardware
	statement := "select id_hardware, name, type, description from hardware where lower(type) = 'sensor'"
	rows, err := db.Query(statement)
	e.PanicIfNeeded(err)
	for rows.Next() {
		err := rows.Scan(&hardware.Id, &hardware.Name, &hardware.Type, &hardware.Description)
		e.PanicIfNeeded(err)
		hardwares = append(hardwares, hardware)
	}
	return hardwares
}

func GetAllHardwareTypeNode() []Hardware {
	var hardware Hardware
	var hardwares []Hardware
	statement := "select id_hardware, name, type, description from hardware where lower(type) = 'single-board computer' or lower(type) = 'microcontroller unit'"
	rows, err := db.Query(statement)
	e.PanicIfNeeded(err)
	for rows.Next() {
		err := rows.Scan(&hardware.Id, &hardware.Name, &hardware.Type, &hardware.Description)
		e.PanicIfNeeded(err)
		hardwares = append(hardwares, hardware)
	}
	return hardwares
}

func GetNodeByHardwareId(id int) Node {
	var node Node
	statement := "select name, location from node where id_hardware = $1"
	err := db.QueryRow(statement, id).Scan(&node.Name, &node.Location)
	if err != nil && err != sql.ErrNoRows {
		e.PanicIfNeeded(err)
	}
	return node
}

func GetSensorByHardwareId(id int) Sensor {
	var sensor Sensor
	statement := "select name, unit from sensor where id_hardware = $1"
	err := db.QueryRow(statement, id).Scan(&sensor.Name, &sensor.Unit)
	if err != nil && err != sql.ErrNoRows {
		e.PanicIfNeeded(err)
	}
	return sensor
}

func UpdateHardware(h HardwareUpdate, id int) {
	statement := "update hardware SET name=COALESCE($1, name), type=COALESCE($2, type), description=COALESCE($3, description) WHERE id_hardware=$4"
	_, err := db.Exec(statement, h.Name, h.Type, h.Description, id)
	e.PanicIfNeeded(err)

}

func DeleteHardware(id int) {
	statement := "delete from hardware where id_hardware = $1"
	_, err := db.Exec(statement, id)
	e.PanicIfNeeded(err)
}
