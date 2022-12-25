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
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Description string `json:"description" binding:"required"`
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
	statement := replaceQueryParam("select id_hardware from hardware where id_hardware = %d", id)
	return isRowExist(statement)
}

func IsHardwareTypedSensorById(id int) bool {
	statement := replaceQueryParam("select type from hardware where id_hardware = %d and (lower(type) = 'sensor')", id)
	return isRowExist(statement)
}

func IsHardwareTypedNodeById(id int) bool {
	statement := replaceQueryParam("select type from hardware where id_hardware = %d and (lower(type) = 'single-board computer' or lower(type) = 'microcontroller unit')", id)
	return isRowExist(statement)
}

func AddHardware(h Hardware) {
	statement := replaceQueryParam("insert into hardware (name, type, description) values ('%s', '%s', '%s')", h.Name, h.Type, h.Description)
	_, err := db.Exec(statement)
	e.PanicIfNeeded(err)
}

func GetHardwareById(id int) Hardware {
	var hardware Hardware
	statement := replaceQueryParam("select id_hardware, name, type, description from hardware where id_hardware = %d", id)
	err := db.QueryRow(statement).Scan(&hardware.Id, &hardware.Name, &hardware.Type, &hardware.Description)
	if err != nil && err != sql.ErrNoRows {
		e.PanicIfNeeded(err)
	}

	return hardware
}

func GetAllHardwareTypedSensor() []Hardware {
	var hardware Hardware
	var hardwares []Hardware
	hardwares = make([]Hardware, 0)
	statement := "select id_hardware, name, type, description from hardware where lower(type) = 'sensor'"
	rows, err := db.Query(statement)
	e.PanicIfNeeded(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&hardware.Id, &hardware.Name, &hardware.Type, &hardware.Description)
		e.PanicIfNeeded(err)
		hardwares = append(hardwares, hardware)
	}
	return hardwares
}

func GetAllHardwareTypedNode() []Hardware {
	var hardware Hardware
	var hardwares []Hardware
	hardwares = make([]Hardware, 0)
	statement := "select id_hardware, name, type, description from hardware where lower(type) = 'single-board computer' or lower(type) = 'microcontroller unit'"
	rows, err := db.Query(statement)
	e.PanicIfNeeded(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&hardware.Id, &hardware.Name, &hardware.Type, &hardware.Description)
		e.PanicIfNeeded(err)
		hardwares = append(hardwares, hardware)
	}
	return hardwares
}

func GetNodeByHardwareId(id int) Node {
	var node Node
	statement := replaceQueryParam("select name, location from node where id_hardware = %d", id)
	err := db.QueryRow(statement).Scan(&node.Name, &node.Location)
	if err != nil && err != sql.ErrNoRows {
		e.PanicIfNeeded(err)
	}
	return node
}

func GetSensorByHardwareId(id int) Sensor {
	var sensor Sensor
	statement := replaceQueryParam("select name, unit from sensor where id_hardware = %d", id)
	err := db.QueryRow(statement).Scan(&sensor.Name, &sensor.Unit)
	if err != nil && err != sql.ErrNoRows {
		e.PanicIfNeeded(err)
	}
	return sensor
}

func UpdateHardware(h HardwareUpdate, id int) {
	fillByNullIfNeeded(&h.Name, &h.Type, &h.Description)
	statement := replaceQueryParam("update hardware SET name=COALESCE(%s, name), type=COALESCE(%s, type), description=COALESCE(%s, description) where id_hardware=%d", *h.Name, *h.Type, *h.Description, id)
	_, err := db.Exec(statement)
	e.PanicIfNeeded(err)
}

func DeleteHardware(id int) error {
	statement := replaceQueryParam("delete from hardware where id_hardware = %d", id)
	_, err := db.Exec(statement)
	// e.PanicIfNeeded(err)
	return err
}
