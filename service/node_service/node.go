package node_service

import (
	"src/models"
	"src/service/cache_service"
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
		nodes_cached, found := cache_service.Get("nodes", id_user)
		if !found {
			nodes = models.GetAllNodeByUserId(id_user)
		} else {
			nodes = nodes_cached.([]models.NodeList)
		}
	}
	return nodes
}

func (n *Node) Get() models.NodeGet {

	var node models.Node
	var user models.User
	node_cached, found := cache_service.Get("node", n.Id)
	if !found {
		node, user = models.GetNodeAndUserByNodeId(n.Id)
		var cn models.CachedNode
		cn.Node = node
		cn.User = user
		cache_service.Set("node", n.Id, cn)
	} else {
		node = node_cached.(models.CachedNode).Node
		user = node_cached.(models.CachedNode).User
	}

	hardware := models.GetHardwareByNodeId(n.Id)
	sensors := models.GetSensorByNodeId(n.Id)

	var resp models.NodeGet

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
	cache_service.Del("node", n.Id)
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
