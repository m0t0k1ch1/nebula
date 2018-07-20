package graph

import (
	"errors"
	"sync"
)

var (
	ErrNodeNotExist = errors.New("graph: the node does not exist in the graph")
	ErrEdgeNotExist = errors.New("graph: the edge does not exist in the graph")
)

type Graph interface {
	IsDirected() bool
	GetNode(ID) (Node, error)
	GetNodes() (map[ID]Node, error)
	GetTails(ID) (map[ID]Node, error)
	GetHeads(ID) (map[ID]Node, error)
	AddNode(Node) error
	RemoveNode(ID) error
	GetEdge(ID, ID) (Edge, error)
	GetEdges() (map[ID]map[ID]Edge, error)
	AddEdge(ID, ID, float64) error
	RemoveEdge(ID, ID) error
}

type graph struct {
	mu         sync.RWMutex
	isDirected bool
	idToNodes  map[ID]Node
	idToTails  map[ID]map[ID]Edge
	idToHeads  map[ID]map[ID]Edge
}

func newGraph(isDirected bool) Graph {
	return &graph{
		isDirected: isDirected,
		idToNodes:  make(map[ID]Node),
		idToTails:  make(map[ID]map[ID]Edge),
		idToHeads:  make(map[ID]map[ID]Edge),
	}
}

func NewDirected() Graph {
	return newGraph(true)
}

func NewUndirected() Graph {
	return newGraph(false)
}

func (g *graph) IsDirected() bool {
	return g.isDirected
}

func (g *graph) isExistNode(id ID) (exists bool) {
	_, exists = g.idToNodes[id]
	return
}

func (g *graph) GetNode(id ID) (Node, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.isExistNode(id) {
		return nil, ErrNodeNotExist
	}

	return g.idToNodes[id], nil
}

func (g *graph) GetNodes() (map[ID]Node, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.idToNodes, nil
}

func (g *graph) GetTails(idHead ID) (map[ID]Node, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.isExistNode(idHead) {
		return nil, ErrNodeNotExist
	}

	tails := make(map[ID]Node)
	if _, ok := g.idToTails[idHead]; ok {
		for id := range g.idToTails[idHead] {
			tails[id] = g.idToNodes[id]
		}
	}

	return tails, nil
}

func (g *graph) GetHeads(idTail ID) (map[ID]Node, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.isExistNode(idTail) {
		return nil, ErrNodeNotExist
	}

	heads := make(map[ID]Node)
	if _, ok := g.idToHeads[idTail]; ok {
		for id := range g.idToHeads[idTail] {
			heads[id] = g.idToNodes[id]
		}
	}

	return heads, nil
}

func (g *graph) AddNode(n Node) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.isExistNode(n.ID()) {
		return nil
	}

	g.idToNodes[n.ID()] = n

	return nil
}

func (g *graph) RemoveNode(id ID) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.isExistNode(id) {
		return nil
	}

	delete(g.idToNodes, id)

	delete(g.idToTails, id)
	for _, tails := range g.idToTails {
		delete(tails, id)
	}

	delete(g.idToHeads, id)
	for _, heads := range g.idToHeads {
		delete(heads, id)
	}

	return nil
}

func (g *graph) newEdge(idTail, idHead ID, weight float64) Edge {
	return newEdge(g.IsDirected(), g.idToNodes[idTail], g.idToNodes[idHead], weight)
}

func (g *graph) GetEdge(idTail, idHead ID) (Edge, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.isExistNode(idTail) || !g.isExistNode(idHead) {
		return nil, ErrNodeNotExist
	}

	if _, ok := g.idToHeads[idTail]; ok {
		if _, ok := g.idToHeads[idTail][idHead]; ok {
			return g.idToHeads[idTail][idHead], nil
		}
	}

	return nil, ErrEdgeNotExist
}

func (g *graph) GetEdges() (map[ID]map[ID]Edge, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.idToHeads, nil
}

func (g *graph) addEdge(idTail, idHead ID, weight float64) {
	e := g.newEdge(idTail, idHead, weight)

	if _, ok := g.idToTails[idHead]; ok {
		if _, ok := g.idToTails[idHead][idTail]; ok {
			g.idToTails[idHead][idTail].AddWeight(e.Weight())
		} else {
			g.idToTails[idHead][idTail] = e
		}
	} else {
		g.idToTails[idHead] = map[ID]Edge{
			idTail: e,
		}
	}

	if _, ok := g.idToHeads[idTail]; ok {
		if _, ok := g.idToHeads[idTail][idHead]; ok {
			g.idToHeads[idTail][idHead].AddWeight(e.Weight())
		} else {
			g.idToHeads[idTail][idHead] = e
		}
	} else {
		g.idToHeads[idTail] = map[ID]Edge{
			idHead: e,
		}
	}

	return
}

func (g *graph) AddEdge(idTail, idHead ID, weight float64) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.isExistNode(idTail) || !g.isExistNode(idHead) {
		return ErrNodeNotExist
	}

	g.addEdge(idTail, idHead, weight)
	if !g.IsDirected() {
		g.addEdge(idHead, idTail, weight)
	}

	return nil
}

func (g *graph) removeEdge(idTail, idHead ID) {
	if _, ok := g.idToTails[idHead]; ok {
		if _, ok := g.idToTails[idHead][idTail]; ok {
			delete(g.idToTails[idHead], idTail)
			if len(g.idToTails[idHead]) == 0 {
				delete(g.idToTails, idHead)
			}
		}
	}

	if _, ok := g.idToHeads[idTail]; ok {
		if _, ok := g.idToHeads[idTail][idHead]; ok {
			delete(g.idToHeads[idTail], idHead)
			if len(g.idToHeads[idTail]) == 0 {
				delete(g.idToHeads, idTail)
			}
		}
	}

	return
}

func (g *graph) RemoveEdge(idTail, idHead ID) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.isExistNode(idTail) || !g.isExistNode(idHead) {
		return ErrNodeNotExist
	}

	g.removeEdge(idTail, idHead)
	if !g.IsDirected() {
		g.removeEdge(idHead, idTail)
	}

	return nil
}
