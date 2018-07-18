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
	GetNode(ID) (Node, error)
	GetNodes() (map[ID]Node, error)
	GetTails(ID) (map[ID]Node, error)
	GetHeads(ID) (map[ID]Node, error)
	AddNode(Node) error
	RemoveNode(ID) error
	AddEdge(ID, ID, float64) error
	RemoveEdge(ID, ID) error
	GetWeight(ID, ID) (float64, error)
}

type graph struct {
	mu         sync.RWMutex
	idToNodes  map[ID]Node
	idToTails  map[ID]map[ID]float64
	idToHeads  map[ID]map[ID]float64
	isDirected bool
}

func newGraph(isDirected bool) Graph {
	return &graph{
		idToNodes:  make(map[ID]Node),
		idToTails:  make(map[ID]map[ID]float64),
		idToHeads:  make(map[ID]map[ID]float64),
		isDirected: isDirected,
	}
}

func NewDirected() Graph {
	return newGraph(true)
}

func NewUndirected() Graph {
	return newGraph(false)
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

func (g *graph) addEdge(idTail, idHead ID, weight float64) {
	if _, ok := g.idToTails[idHead]; ok {
		if _, ok := g.idToTails[idHead][idTail]; ok {
			g.idToTails[idHead][idTail] += weight
		} else {
			g.idToTails[idHead][idTail] = weight
		}
	} else {
		g.idToTails[idHead] = map[ID]float64{
			idTail: weight,
		}
	}

	if _, ok := g.idToHeads[idTail]; ok {
		if _, ok := g.idToHeads[idTail][idHead]; ok {
			g.idToHeads[idTail][idHead] += weight
		} else {
			g.idToHeads[idTail][idHead] = weight
		}
	} else {
		g.idToHeads[idTail] = map[ID]float64{
			idHead: weight,
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
	if !g.isDirected {
		g.addEdge(idHead, idTail, weight)
	}

	return nil
}

func (g *graph) removeEdge(idTail, idHead ID) {
	if _, ok := g.idToTails[idHead]; ok {
		if _, ok := g.idToTails[idHead][idTail]; ok {
			delete(g.idToTails[idHead], idTail)
		}
	}

	if _, ok := g.idToHeads[idTail]; ok {
		if _, ok := g.idToHeads[idTail][idHead]; ok {
			delete(g.idToHeads[idTail], idHead)
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
	if !g.isDirected {
		g.removeEdge(idHead, idTail)
	}

	return nil
}

func (g *graph) GetWeight(idTail, idHead ID) (float64, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.isExistNode(idTail) || !g.isExistNode(idHead) {
		return 0, ErrNodeNotExist
	}

	if _, ok := g.idToHeads[idTail]; ok {
		if _, ok := g.idToHeads[idTail][idHead]; ok {
			return g.idToHeads[idTail][idHead], nil
		}
	}

	return 0, ErrEdgeNotExist
}
