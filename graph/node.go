package graph

type Node interface {
	ID() ID
}

type node struct {
	id StringID
}

func NewNode(id string) Node {
	return &node{
		id: StringID(id),
	}
}

func (n *node) ID() ID {
	return n.id
}
