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
