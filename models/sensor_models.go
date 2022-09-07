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

func GetUserIdBySensorId(id int) int {
	statement := "select node.id_user from node left join sensor on sensor.id_node = node.id_node where id_sensor = $1"
	var id_user int
	err := db.QueryRow(statement, id).Scan(&id_user)
	e.PanicIfNeeded(err)
	return id_user
}

func IsSensorExistById(id int) bool {
	statement := "select id_sensor from sensor where id_sensor = $1"
	return isRowExist(statement, id)
}

func DeleteSensor(id int) {
	statement := "delete from sensor where id_sensor = $1"
	_, err := db.Exec(statement, id)
	e.PanicIfNeeded(err)
}
