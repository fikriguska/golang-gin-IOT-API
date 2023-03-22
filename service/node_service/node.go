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
	if n.Id_hardware_node != -1 {
		models.AddNode(n.Node)
	} else {
		models.AddNodeNoHardware(n.Node)
	}
}

func (n *Node) GetAll(id_user int, is_admin bool, args ...int) []models.NodeList {
	// the default of limit is 50
	limit := 50
	if len(args) == 1 {
		if args[0] >= 1 {
			limit = args[0]
		}
	}
	var nodes []models.NodeList
	nodes_cached, found := cache_service.Get("nodes", id_user)
	if !found {
		if is_admin {
			nodes = models.GetAllNode()
		} else {
			nodes = models.GetAllNodeByUserId(id_user)
		}
		cache_service.Set("nodes", id_user, nodes)
	} else {
		nodes = nodes_cached.([]models.NodeList)
	}

	for i, n := range nodes {
		feeds := models.GetFeedByNodeId(n.Id, limit)
		nodes[i].Feed = feeds
	}
	return nodes
}

func (n *Node) Get(args ...int) models.NodeGet {

	// the default of limit is 50
	limit := 50
	if len(args) == 1 {
		if args[0] >= 1 {
			limit = args[0]
		}
	}

	var node models.NodeGet

	node_cached, found := cache_service.Get("node", n.Id)

	if !found {
		node = models.GetNodeById(n.Id)
		cache_service.Set("node", n.Id, node)
	} else {
		node = node_cached.(models.NodeGet)
	}

	feed := models.GetFeedByNodeId(n.Id, limit)
	node.Feed = feed

	return node

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

func (n *Node) IsPublic() (public bool) {
	return models.IsNodePublic(n.Id)
}

func (n *Node) Delete() {
	models.DeleteNode(n.Id)
}
