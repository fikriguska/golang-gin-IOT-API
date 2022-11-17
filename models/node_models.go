package models

import (
	"fmt"
	e "src/error"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Node struct {
	Id                 int
	Name               string
	Location           string
	Id_user            int
	Id_hardware_node   int
	Id_hardware_sensor []int `pg:"type:uint[10],array,pk"`
	Field_sensor       []string
	Is_public          bool
}

type NodeAdd struct {
	Name               string   `json:"name" binding:"required"`
	Location           string   `json:"location" binding:"required"`
	Id_hardware_node   *int     `json:"id_hardware_node"`
	Id_hardware_sensor []int    `json:"id_hardware_sensor" binding:"required"`
	Field_sensor       []string `json:"field_sensor" binding:"required"`
	Is_public          *bool    `json:"is_public"`
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
	id_sensors := "{" + strings.Replace(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(node.Id_hardware_sensor)), ","), "[]"), "-1", "NULL", -1) + "}"
	field_sensors := "{" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(node.Field_sensor)), ","), "[]") + "}"
	statement := "insert into node (name, location, id_user, id_hardware_node, id_hardware_sensor, field_sensor, is_public) values ($1, $2, $3, $4, $5, $6, $7)"
	_, err := db.Exec(cb(), statement, node.Name, node.Location, node.Id_user, nil, id_sensors, field_sensors, node.Is_public)
	e.PanicIfNeeded(err)
}

func AddNode(node Node) {
	id_sensors := "{" + strings.Replace(strings.Trim(strings.Join(strings.Fields(fmt.Sprint(node.Id_hardware_sensor)), ","), "[]"), "-1", "NULL", -1) + "}"
	field_sensors := "{" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(node.Field_sensor)), ","), "[]") + "}"
	statement := "insert into node (name, location, id_user, id_hardware_node, id_hardware_sensor, field_sensor, is_public) values ($1, $2, $3, $4, $5, $6, $7)"
	_, err := db.Exec(cb(), statement, node.Name, node.Location, node.Id_user, node.Id_hardware_node, id_sensors, field_sensors, node.Is_public)
	e.PanicIfNeeded(err)
}

func GetAllNodeByUserId(id_user int) []NodeList {
	var node NodeList
	var nodes []NodeList
	nodes = make([]NodeList, 0)
	statement := "select id_node, name, location, id_hardware, id_user from node where id_user = $1"
	rows, err := db.Query(cb(), statement, id_user)
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
	rows, err := db.Query(cb(), statement)
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
	statement := "select node.id_node, node.name, node.location, user_person.id_user, user_person.username from node left join user_person on node.id_user = user_person.id_user where node.id_node = $1"
	var node Node
	var user User
	err := db.QueryRow(cb(), statement, id).Scan(&node.Id, &node.Name, &node.Location, &user.Id, &user.Username)
	e.PanicIfNeeded(err)
	return node, user
}

func GetHardwareByNodeId(id int) Hardware {
	statement := "select hardware.name, hardware.type from hardware left join node on hardware.id_hardware = node.id_hardware where id_node = $1"
	var hardware Hardware
	err := db.QueryRow(cb(), statement, id).Scan(&hardware.Name, &hardware.Type)
	if err != nil && err != pgx.ErrNoRows {
		e.PanicIfNeeded(err)
	}
	return hardware
}

func GetSensorByNodeId(id int) []Sensor {
	var sensors []Sensor
	var sensor Sensor
	sensors = make([]Sensor, 0)
	statement := "select sensor.id_sensor, sensor.name, sensor.unit from sensor left join node on sensor.id_node = node.id_node where sensor.id_node = $1"
	rows, err := db.Query(cb(), statement, id)
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
	statement := "select id_node from node where id_node = $1"
	return isRowExist(statement, id)
}

func UpdateNode(n NodeUpdate, id int) {
	statement := "update node SET name=COALESCE($1, name), location=COALESCE($2, location) where id_node=$3"
	_, err := db.Exec(cb(), statement, n.Name, n.Location, id)
	e.PanicIfNeeded(err)
}

func DeleteNode(id int) {
	statement := "delete from node where id_node = $1"
	_, err := db.Exec(cb(), statement, id)
	e.PanicIfNeeded(err)
}

func GetUserIdByNodeId(id int) int {
	statement := "select id_user from node where id_node = $1"
	var id_user int
	err := db.QueryRow(cb(), statement, id).Scan(&id_user)
	e.PanicIfNeeded(err)
	return id_user
}
