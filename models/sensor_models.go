package models

import e "src/error"

type Sensor struct {
	Id          int
	Name        string
	Unit        string
	Id_node     int
	Id_hardware int
}

func AddSensorNoHardware(s Sensor) {
	statement := "insert into sensor (name, unit, id_node, id_hardware) values ($1, $2, $3, $4)"
	_, err := db.Exec(statement, s.Name, s.Unit, s.Id_node, nil)
	e.PanicIfNeeded(err)
}

func AddSensor(s Sensor) {
	statement := "insert into sensor (name, unit, id_node, id_hardware) values ($1, $2, $3, $4)"
	_, err := db.Exec(statement, s.Name, s.Unit, s.Id_node, s.Id_hardware)
	e.PanicIfNeeded(err)
}
