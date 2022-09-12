package models

import (
	e "src/error"
	"time"
)

type Sensor struct {
	Id          int
	Name        string
	Unit        string
	Id_node     int
	Id_hardware int
}

type SensorList struct {
	Id          int
	Name        string
	Unit        string
	Id_node     int
	Id_hardware NullInt64
}

type SensorGet struct {
	Id      int
	Name    string
	Unit    string
	Channel []SensorChannelGet
}

type SensorChannelGet struct {
	Time  time.Time
	Value float64
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

func GetAllSensor() []SensorList {
	var sensor SensorList
	var sensors []SensorList
	statement := "select id_sensor, name, unit, id_hardware, id_node from sensor"
	rows, err := db.Query(statement)
	e.PanicIfNeeded(err)
	for rows.Next() {
		err := rows.Scan(&sensor.Id, &sensor.Name, &sensor.Unit, &sensor.Id_hardware, &sensor.Id_node)
		e.PanicIfNeeded(err)
		sensors = append(sensors, sensor)
	}

	return sensors
}

func GetAllSensorByUserId(id_user int) []SensorList {
	var sensor SensorList
	var sensors []SensorList
	statement := "select sensor.id_sensor, sensor.name, sensor.unit, sensor.id_hardware, sensor.id_node from sensor left join node on sensor.id_node = node.id_node where node.id_user = $1"
	rows, err := db.Query(statement, id_user)
	e.PanicIfNeeded(err)
	for rows.Next() {
		err := rows.Scan(&sensor.Id, &sensor.Name, &sensor.Unit, &sensor.Id_hardware, &sensor.Id_node)
		e.PanicIfNeeded(err)
		sensors = append(sensors, sensor)
	}
	return sensors
}

func GetUserIdBySensorId(id int) int {
	statement := "select node.id_user from node left join sensor on sensor.id_node = node.id_node where id_sensor = $1"
	var id_user int
	err := db.QueryRow(statement, id).Scan(&id_user)
	e.PanicIfNeeded(err)
	return id_user
}

func GetSensorById(id int) Sensor {
	statement := "select id_sensor, name, unit from sensor where id_sensor = $1"
	var sensor Sensor
	err := db.QueryRow(statement, id).Scan(&sensor.Id, &sensor.Name, &sensor.Unit)
	e.PanicIfNeeded(err)
	return sensor
}

func GetChannelBySensorId(id int) []Channel {
	var channels []Channel
	var channel Channel
	statement := "select time, value from channel where id_sensor = $1"
	rows, err := db.Query(statement, id)
	e.PanicIfNeeded(err)
	for rows.Next() {
		err := rows.Scan(&channel.Time, &channel.Value)
		e.PanicIfNeeded(err)
		channels = append(channels, channel)
	}
	return channels
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
