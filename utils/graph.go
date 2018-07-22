package utils

import (
	"github.com/awalterschulze/gographviz"
	"github.com/m0t0k1ch1/nebula/graph"
)

func NewDOTGraph(g graph.Graph) (*gographviz.Graph, error) {
	gv := gographviz.NewGraph()

	if err := gv.SetDir(g.IsDirected()); err != nil {
		return nil, err
	}

	// add nodes
	nodes, err := g.GetNodes()
	if err != nil {
		return nil, err
	}
	for _, n := range nodes {
		if err := gv.AddNode(gv.Name, n.ID().String(), nil); err != nil {
			return nil, err
		}
	}

	// add edges
	edges, err := g.GetEdges()
	if err != nil {
		return nil, err
	}
	for _, nodeEdges := range edges {
		for _, e := range nodeEdges {
			if err := gv.AddEdge(
				e.Tail().ID().String(),
				e.Head().ID().String(),
				e.IsDirected(),
				nil,
			); err != nil {
				return nil, err
			}
			if !g.IsDirected() {
				// remove reversed edge
				if _, ok := edges[e.Head().ID()]; ok {
					if _, ok := edges[e.Head().ID()][e.Tail().ID()]; ok {
						delete(edges[e.Head().ID()], e.Tail().ID())
					}
				}
			}
		}
	}

	return gv, nil
}
