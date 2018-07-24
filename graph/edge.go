package graph

type Edge struct {
	isDirected bool
	tail       *Node
	head       *Node
	weight     float64
}

func newEdge(isDirected bool, nTail, nHead *Node, weight float64) *Edge {
	return &Edge{
		isDirected: isDirected,
		tail:       nTail,
		head:       nHead,
		weight:     weight,
	}
}

func NewDirectedEdge(nTail, nHead *Node, weight float64) *Edge {
	return newEdge(true, nTail, nHead, weight)
}

func NewUndirectedEdge(nTail, nHead *Node, weight float64) *Edge {
	return newEdge(false, nTail, nHead, weight)
}

func (e *Edge) IsDirected() bool {
	return e.isDirected
}

func (e *Edge) Tail() *Node {
	return e.tail
}

func (e *Edge) Head() *Node {
	return e.head
}

func (e *Edge) Weight() float64 {
	return e.weight
}

func (e *Edge) SetWeight(weight float64) {
	e.weight = weight
}
