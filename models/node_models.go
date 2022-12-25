package models

import (
	"database/sql"
	e "src/error"
)

type Node struct {
	Id          int
	Name        string
	Location    string
	Id_user     int
	Id_hardware int
}

type NodeAdd struct {
	Name        string `json:"name" binding:"required"`
	Location    string `json:"location" binding:"required"`
	Id_hardware *int   `json:"id_hardware"`
}

type NodeList struct {
	Id          int       `json:"id_node"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Id_user     int       `json:"id_user"`
	Id_hardware NullInt64 `json:"id_hardware"`
}

type NodeSensorGet struct {
	Id_sensor int    `json:"id_sensor"`
	Name      string `json:"name"`
	Unit      string `json:"unit"`
}
type NodeHardwareGet struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
type NodeGet struct {
	Id       int               `json:"id_node"`
	Name     string            `json:"name"`
	Location string            `json:"location"`
	Id_user  int               `json:"id_user"`
	Username string            `json:"username"`
	Hardware []NodeHardwareGet `json:"hardware"`
	Sensor   []NodeSensorGet   `json:"sensor"`
}

type NodeUpdate struct {
	Name     *string `json:"name"`
	Location *string `json:"location"`
}

func AddNodeNoHardware(node Node) {
	statement := replaceQueryParam("insert into node (name, location, id_user, id_hardware) values ('%s', '%s', %d, NULL)", node.Name, node.Location, node.Id_user)
	_, err := db.Exec(statement)
	e.PanicIfNeeded(err)
}

func AddNode(node Node) {
	statement := replaceQueryParam("insert into node (name, location, id_user, id_hardware) values ('%s', '%s', %d, %d)", node.Name, node.Location, node.Id_user, node.Id_hardware)
	_, err := db.Exec(statement)
	e.PanicIfNeeded(err)
}

func GetAllNodeByUserId(id_user int) []NodeList {
	var node NodeList
	var nodes []NodeList
	nodes = make([]NodeList, 0)
	statement := replaceQueryParam("select id_node, name, location, id_hardware, id_user from node where id_user = %d", id_user)
	rows, err := db.Query(statement)
	e.PanicIfNeeded(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&node.Id, &node.Name, &node.Location, &node.Id_hardware, &node.Id_user)
		e.PanicIfNeeded(err)
		nodes = append(nodes, node)
	}
	return nodes
}

func GetAllNode() []NodeList {
	var node NodeList
	var nodes []NodeList
	nodes = make([]NodeList, 0)
	statement := "select id_node, name, location, id_hardware, id_user from node"
	rows, err := db.Query(statement)
	e.PanicIfNeeded(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&node.Id, &node.Name, &node.Location, &node.Id_hardware, &node.Id_user)
		e.PanicIfNeeded(err)
		nodes = append(nodes, node)
	}
	return nodes
}

func GetNodeAndUserByNodeId(id int) (Node, User) {
	statement := replaceQueryParam("select node.id_node, node.name, node.location, user_person.id_user, user_person.username from node left join user_person on node.id_user = user_person.id_user where node.id_node = %d", id)
	var node Node
	var user User
	err := db.QueryRow(statement).Scan(&node.Id, &node.Name, &node.Location, &user.Id, &user.Username)
	e.PanicIfNeeded(err)
	return node, user
}

func GetHardwareByNodeId(id int) Hardware {
	statement := replaceQueryParam("select hardware.name, hardware.type from hardware left join node on hardware.id_hardware = node.id_hardware where id_node = %d", id)
	var hardware Hardware
	err := db.QueryRow(statement).Scan(&hardware.Name, &hardware.Type)
	if err != nil && err != sql.ErrNoRows {
		e.PanicIfNeeded(err)
	}
	return hardware
}

func GetSensorByNodeId(id int) []Sensor {
	var sensors []Sensor
	var sensor Sensor
	sensors = make([]Sensor, 0)
	statement := replaceQueryParam("select sensor.id_sensor, sensor.name, sensor.unit from sensor left join node on sensor.id_node = node.id_node where sensor.id_node = %d", id)
	rows, err := db.Query(statement)
	e.PanicIfNeeded(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&sensor.Id, &sensor.Name, &sensor.Unit)
		e.PanicIfNeeded(err)
		sensors = append(sensors, sensor)
	}
	return sensors
}

func IsNodeExistById(id int) bool {
	statement := replaceQueryParam("select id_node from node where id_node = %d", id)
	return isRowExist(statement)
}

func UpdateNode(n NodeUpdate, id int) {
	fillByNullIfNeeded(&n.Name, &n.Location)
	statement := replaceQueryParam("update node SET name=COALESCE(%s, name), location=COALESCE(%s, location) where id_node = %d", *n.Name, *n.Location, id)
	_, err := db.Exec(statement)
	e.PanicIfNeeded(err)
}

func DeleteNode(id int) {
	statement := replaceQueryParam("delete from node where id_node = %d", id)
	_, err := db.Exec(statement)
	e.PanicIfNeeded(err)
}

func GetUserIdByNodeId(id int) int {
	statement := replaceQueryParam("select id_user from node where id_node = %d", id)
	var id_user int
	err := db.QueryRow(statement).Scan(&id_user)
	e.PanicIfNeeded(err)
	return id_user
}
