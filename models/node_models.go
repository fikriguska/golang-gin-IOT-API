package models

import (
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

type NodeGet struct {
	Id          int       `json:"id_node"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Id_user     int       `json:"id_user"`
	Id_hardware NullInt64 `json:"id_hardware"`
}

func AddNodeNoHardware(node Node) {
	// tpdo: auth
	statement := "insert into node (name, location, id_user, id_hardware) values ($1, $2, $3, $4)"
	_, err := db.Exec(statement, node.Name, node.Location, node.Id_user, nil)
	e.PanicIfNeeded(err)
}

func AddNode(node Node) {
	// tpdo: auth
	statement := "insert into node (name, location, id_user, id_hardware) values ($1, $2, $3, $4)"
	_, err := db.Exec(statement, node.Name, node.Location, node.Id_user, node.Id_hardware)
	e.PanicIfNeeded(err)
}

func GetAllNodeByUserId(id_user int) []NodeGet {
	var node NodeGet
	var nodes []NodeGet
	statement := "select id_node, name, location, id_hardware, id_user from node where id_user = $1"
	rows, err := db.Query(statement, id_user)
	e.PanicIfNeeded(err)
	for rows.Next() {
		err := rows.Scan(&node.Id, &node.Name, &node.Location, &node.Id_hardware, &node.Id_user)
		e.PanicIfNeeded(err)
		nodes = append(nodes, node)
	}
	return nodes
}

func GetAllNode() []NodeGet {
	var node NodeGet
	var nodes []NodeGet
	statement := "select * from node"
	rows, err := db.Query(statement)
	e.PanicIfNeeded(err)
	for rows.Next() {
		err := rows.Scan(&node.Id, &node.Name, &node.Location, &node.Id_hardware, &node.Id_user)
		e.PanicIfNeeded(err)
		nodes = append(nodes, node)
	}
	return nodes
}

func IsNodeExistById(id int) bool {
	statement := "select id_node from node where id_node = $1"
	return isRowExist(statement, id)
}

func DeleteNode(id int) {
	statement := "delete from node where id_node = $1"
	_, err := db.Exec(statement, id)
	e.PanicIfNeeded(err)
}

func GetUserIdByNodeId(id int) int {
	statement := "select id_user from node where id_node = $1"
	var id_user int
	err := db.QueryRow(statement, id).Scan(&id_user)
	e.PanicIfNeeded(err)
	return id_user
}
