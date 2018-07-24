package graph

type Node struct {
	id ID
}

func NewNode(id string) *Node {
	return &Node{
		id: ID(id),
	}
}

func (n *Node) ID() ID {
	return n.id
}
