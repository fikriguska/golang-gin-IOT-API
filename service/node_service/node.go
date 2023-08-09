package node_service

import (
	"src/models"
	"src/service/cache_service"
	"src/service/hardware_service"
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
			cache_service.Set("nodes", id_user, nodes)
		} else {
			nodes = nodes_cached.([]models.NodeList)
		}
	}
	return nodes
}

func (n *Node) Get() models.NodeGet {

	var node models.Node
	var user models.User
	var hardware interface{}
	node_cached, found := cache_service.Get("node", n.Id)
	if !found {
		node, user = models.GetNodeAndUserByNodeId(n.Id)
		var cn models.CachedNode
		cn.Id = node.Id
		cn.Name = node.Name
		cn.Location = node.Location
		cn.Id_hardware = node.Id_hardware
		cn.Id_user = user.Id
		cn.Username = user.Username
		cache_service.Set("node", n.Id, cn)
	} else {
		cn := node_cached.(models.CachedNode)
		node.Id = cn.Id
		node.Name = cn.Name
		node.Location = cn.Location
		node.Id_hardware = cn.Id_hardware
		user.Id = cn.Id_user
		user.Username = cn.Username
	}

	hardware_cached, found := cache_service.Get("hardware", node.Id_hardware)
	if !found {
		hardwareService := hardware_service.Hardware{
			Hardware: models.Hardware{
				Id: node.Id_hardware,
			},
		}
		hardware = hardwareService.Get().(models.HardwareNodeGet)
		cache_service.Set("hardware", node.Id_hardware, hardware)
	} else {
		hardware = hardware_cached.(models.HardwareNodeGet)
	}

	sensors_cached, found := cache_service.Get("sensors-bynode", n.Id)
	var sensors []models.Sensor
	if !found {
		sensors = models.GetSensorByNodeId(n.Id)
		cache_service.Set("sensors-bynode", n.Id, sensors)
	} else {
		sensors = sensors_cached.([]models.Sensor)
	}

	var resp models.NodeGet

	resp.Id = node.Id
	resp.Name = node.Name
	resp.Location = node.Location
	resp.Id_user = user.Id
	resp.Username = user.Username

	resp.Hardware = append(resp.Hardware, models.NodeHardwareGet{})
	resp.Hardware[0].Name = hardware.(models.HardwareNodeGet).Name
	resp.Hardware[0].Type = hardware.(models.HardwareNodeGet).Type

	resp.Sensor = make([]models.NodeSensorGet, 0)
	for i, s := range sensors {
		resp.Sensor = append(resp.Sensor, models.NodeSensorGet{})
		resp.Sensor[i].Id_sensor = s.Id
		resp.Sensor[i].Name = s.Name
		resp.Sensor[i].Unit = s.Unit
	}

	return resp

}

func (n *Node) Update(node models.NodeUpdate, idUser int) {
	models.UpdateNode(node, n.Id)
	updateCache(n, idUser)

}
func (n *Node) IsExistAndOwner(id_user int) (exist bool, owner bool) {

	node_cached, found := cache_service.Get("node", n.Id)

	if !found {
		exist = models.IsNodeExistById(n.Id)
		if !exist {
			return exist, false
		} else {
			node, user := models.GetNodeAndUserByNodeId(n.Id)
			var cn models.CachedNode
			cn.Id = node.Id
			cn.Name = node.Name
			cn.Location = node.Location
			cn.Id_hardware = node.Id_hardware
			cn.Id_user = user.Id
			cn.Username = user.Username
			cache_service.Set("node", n.Id, cn)
		}
		owner = (models.GetUserIdByNodeId(n.Id) == id_user)
		return exist, owner
	} else {
		// log.Println("aaaa")
		owner = node_cached.(models.CachedNode).Id_user == id_user
		return true, owner
	}

}

func (n *Node) Delete() {
	models.DeleteNode(n.Id)
}

func updateCache(n *Node, idUser int) {
	// cache_service.Set("node", n.Id, n)
	nodes, found := cache_service.Get("nodes", idUser)
	if found {
		ns := nodes.([]models.NodeList)
		for idx, node := range ns {
			if node.Id == n.Id {
				ns[idx].Id_hardware = node.Id_hardware
				ns[idx].Location = node.Location
				ns[idx].Name = node.Name
			}
		}
		cache_service.Set("nodes", idUser, ns)
	}
}
