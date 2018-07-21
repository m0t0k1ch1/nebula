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
	GetHeads(ID) (map[ID]Node, error)
	GetTails(ID) (map[ID]Node, error)
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
	nodes      map[ID]Node
	heads      map[ID]map[ID]Node
	tails      map[ID]map[ID]Node
	edges      map[ID]map[ID]Edge
}

func newGraph(isDirected bool) Graph {
	return &graph{
		isDirected: isDirected,
		nodes:      map[ID]Node{},
		heads:      map[ID]map[ID]Node{},
		tails:      map[ID]map[ID]Node{},
		edges:      map[ID]map[ID]Edge{},
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
	_, exists = g.nodes[id]
	return
}

func (g *graph) GetNode(id ID) (Node, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.isExistNode(id) {
		return nil, ErrNodeNotExist
	}

	return g.nodes[id], nil
}

func (g *graph) GetNodes() (map[ID]Node, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.nodes, nil
}

func (g *graph) GetHeads(idTail ID) (map[ID]Node, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.isExistNode(idTail) {
		return nil, ErrNodeNotExist
	}

	if _, ok := g.heads[idTail]; !ok {
		return map[ID]Node{}, nil
	}

	return g.heads[idTail], nil
}

func (g *graph) GetTails(idHead ID) (map[ID]Node, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.isExistNode(idHead) {
		return nil, ErrNodeNotExist
	}

	if _, ok := g.tails[idHead]; !ok {
		return map[ID]Node{}, nil
	}

	return g.tails[idHead], nil
}

func (g *graph) AddNode(n Node) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.isExistNode(n.ID()) {
		return nil
	}

	g.nodes[n.ID()] = n

	return nil
}

func (g *graph) RemoveNode(id ID) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.isExistNode(id) {
		return nil
	}

	delete(g.nodes, id)

	delete(g.heads, id)
	for _, headNodes := range g.heads {
		delete(headNodes, id)
	}

	delete(g.tails, id)
	for _, tailNodes := range g.tails {
		delete(tailNodes, id)
	}

	delete(g.edges, id)
	for _, nodeEdges := range g.edges {
		delete(nodeEdges, id)
	}

	return nil
}

func (g *graph) isExistEdge(idTail, idHead ID) bool {
	if _, ok := g.edges[idTail]; ok {
		if _, ok := g.edges[idTail][idHead]; ok {
			return true
		}
	}
	return false
}

func (g *graph) newEdge(idTail, idHead ID, weight float64) Edge {
	return newEdge(g.isDirected, g.nodes[idTail], g.nodes[idHead], weight)
}

func (g *graph) GetEdge(idTail, idHead ID) (Edge, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.isExistNode(idTail) || !g.isExistNode(idHead) {
		return nil, ErrNodeNotExist
	}

	if !g.isExistEdge(idTail, idHead) {
		return nil, ErrEdgeNotExist
	}

	return g.edges[idTail][idHead], nil
}

func (g *graph) GetEdges() (map[ID]map[ID]Edge, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.edges, nil
}

func (g *graph) addEdge(idTail, idHead ID, weight float64) {
	if _, ok := g.edges[idTail]; ok {
		if _, ok := g.edges[idTail][idHead]; ok {
			e := g.edges[idTail][idHead]
			e.SetWeight(e.Weight() + weight)
		} else {
			g.edges[idTail][idHead] = g.newEdge(idTail, idHead, weight)
		}
	} else {
		g.edges[idTail] = map[ID]Edge{
			idHead: g.newEdge(idTail, idHead, weight),
		}
	}
}

func (g *graph) addRelation(idTail, idHead ID) {
	if _, ok := g.heads[idTail]; ok {
		if _, ok := g.heads[idTail][idHead]; !ok {
			g.heads[idTail][idHead] = g.nodes[idHead]
		}
	} else {
		g.heads[idTail] = map[ID]Node{
			idHead: g.nodes[idHead],
		}
	}

	if _, ok := g.tails[idHead]; ok {
		if _, ok := g.tails[idHead][idTail]; !ok {
			g.tails[idHead][idTail] = g.nodes[idTail]
		}
	} else {
		g.tails[idHead] = map[ID]Node{
			idTail: g.nodes[idTail],
		}
	}
}

func (g *graph) AddEdge(idTail, idHead ID, weight float64) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.isExistNode(idTail) || !g.isExistNode(idHead) {
		return ErrNodeNotExist
	}

	g.addEdge(idTail, idHead, weight)
	g.addRelation(idTail, idHead)

	if !g.isDirected {
		g.addEdge(idHead, idTail, weight)
		g.addRelation(idHead, idTail)
	}

	return nil
}

func (g *graph) removeEdge(idTail, idHead ID) {
	if _, ok := g.edges[idTail]; ok {
		if _, ok := g.edges[idTail][idHead]; ok {
			delete(g.edges[idTail], idHead)
			if len(g.edges[idTail]) == 0 {
				delete(g.edges, idTail)
			}
		}
	}
}

func (g *graph) removeRelation(idTail, idHead ID) {
	if _, ok := g.tails[idHead]; ok {
		if _, ok := g.tails[idHead][idTail]; ok {
			delete(g.tails[idHead], idTail)
			if len(g.tails[idHead]) == 0 {
				delete(g.tails, idHead)
			}
		}
	}

	if _, ok := g.heads[idTail]; ok {
		if _, ok := g.heads[idTail][idHead]; ok {
			delete(g.heads[idTail], idHead)
			if len(g.heads[idTail]) == 0 {
				delete(g.heads, idTail)
			}
		}
	}
}

func (g *graph) RemoveEdge(idTail, idHead ID) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.isExistNode(idTail) || !g.isExistNode(idHead) {
		return ErrNodeNotExist
	}

	g.removeEdge(idTail, idHead)
	g.removeRelation(idTail, idHead)

	if !g.isDirected {
		g.removeEdge(idHead, idTail)
		g.removeRelation(idHead, idTail)
	}

	return nil
}
