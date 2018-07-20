package graph

type Edge interface {
	IsDirected() bool
	Tail() Node
	Head() Node
	Weight() float64
	AddWeight(weight float64) error
}

type edge struct {
	isDirected bool
	tail       Node
	head       Node
	weight     float64
}

func newEdge(isDirected bool, nTail, nHead Node, weight float64) Edge {
	return &edge{
		isDirected: isDirected,
		tail:       nTail,
		head:       nHead,
		weight:     weight,
	}
}

func NewDirectedEdge(nTail, nHead Node, weight float64) Edge {
	return newEdge(true, nTail, nHead, weight)
}

func NewUndirectedEdge(nTail, nHead Node, weight float64) Edge {
	return newEdge(false, nTail, nHead, weight)
}

func (e *edge) IsDirected() bool {
	return e.isDirected
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

func (e *edge) AddWeight(weight float64) error {
	e.weight += weight
	return nil
}
