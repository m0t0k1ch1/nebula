package graph

type Edge interface {
	Tail() Node
	Head() Node
	Weight() float64
}

type edge struct {
	tail   Node
	head   Node
	weight float64
}

func NewEdge(tail, head Node, weight float64) Edge {
	return &edge{
		tail:   tail,
		head:   head,
		weight: weight,
	}
}

func (e *edge) Tail() Node {
	return e.tail
}

func (e *edge) Head() Node {
	return e.head
}

func (e *edge) Weight() float64 {
	return e.weight
}
