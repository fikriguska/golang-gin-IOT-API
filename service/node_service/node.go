package node_service

import (
	"src/models"
)

type Node struct {
	models.Node
}

func (n *Node) Add() {

	// check if there is a hardware
	if n.Id_hardware != -1 {
		models.AddNode(n.Node)
	} else {
		models.AddNodeNoHardware(n.Node)
	}
}

func (n *Node) GetAll(id_user int, is_admin bool) []models.NodeList {
	var nodes []models.NodeList
	if is_admin {
		nodes = models.GetAllNode()
	} else {
		// key := fmt.Sprintf("%d-node", id_user)
		// nodes_byte, err := cache_service.Cache.Get(key)
		// 	if err != nil {
		nodes = models.GetAllNodeByUserId(id_user)
		// 		nodesJson, _ := json.Marshal(nodes)
		// 		cache_service.Cache.Set(key, nodesJson)
		// 	} else {
		// 		// fmt.Println(string(nodes_byte))
		// 		_ = json.Unmarshal(nodes_byte, &nodes)
		// 		// if err != nil {
		// 		// 	fmt.Println(err)
		// 		// }
		// 		// fmt.Println(nodes)
		// 	}
	}
	return nodes
}

func (n *Node) Get() models.NodeGet {

	var user models.User
	var node models.Node
	var hardware models.Hardware
	var sensors []models.Sensor
	var resp models.NodeGet

	node, user = models.GetNodeAndUserByNodeId(n.Id)
	hardware = models.GetHardwareByNodeId(n.Id)
	sensors = models.GetSensorByNodeId(n.Id)

	resp.Id = node.Id
	resp.Name = node.Name
	resp.Location = node.Location
	resp.Id_user = user.Id
	resp.Username = user.Username

	resp.Hardware = append(resp.Hardware, models.NodeHardwareGet{})
	resp.Hardware[0].Name = hardware.Name
	resp.Hardware[0].Type = hardware.Type

	resp.Sensor = make([]models.NodeSensorGet, 0)
	for i, s := range sensors {
		resp.Sensor = append(resp.Sensor, models.NodeSensorGet{})
		resp.Sensor[i].Id_sensor = s.Id
		resp.Sensor[i].Name = s.Name
		resp.Sensor[i].Unit = s.Unit
	}

	return resp

}

func (n *Node) Update(node models.NodeUpdate) {
	models.UpdateNode(node, n.Id)
}

func (n *Node) IsExistAndOwner(id_user int) (exist bool, owner bool) {
	exist = models.IsNodeExistById(n.Id)
	if !exist {
		return exist, false
	}
	owner = (models.GetUserIdByNodeId(n.Id) == id_user)
	return exist, owner
}

func (n *Node) Delete() {
	models.DeleteNode(n.Id)
}
