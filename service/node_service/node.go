package node_service

import (
	"src/models"
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

func (n *Node) GetAll(id_user int, is_admin bool, limit int) []models.NodeList {
	// the default of limit is 50

	var nodes []models.NodeList
	if is_admin {
		nodes = models.GetAllNode()
	} else {
		nodes = models.GetAllNodeByUserId(id_user)
	}

	for i, n := range nodes {
		if limit > 0 {
			feeds := models.GetFeedByNodeId(n.Id, limit)
			nodes[i].Feed = feeds
		} else {
			nodes[i].Feed = nil
		}
	}
	return nodes
}

func (n *Node) Get(limit int) models.NodeGet {

	// the default of limit is 50

	node := models.GetNodeById(n.Id)
	if limit > 0 {
		feed := models.GetFeedByNodeId(n.Id, limit)
		node.Feed = feed
	} else {
		node.Feed = nil
	}

	return node

}

func (n *Node) GetNodeOnly() models.NodeGet {
	node := models.GetNodeById(n.Id)
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
