package utils

import (
	"github.com/awalterschulze/gographviz"
	"github.com/m0t0k1ch1/nebula/graph"
)

var (
	defaultGraphAttrs = map[string]string{
		string(gographviz.Layout): "fdp",
	}
	defaultNodeAttrs = map[string]string{
		string(gographviz.Style):       "\"filled,solid\"",
		string(gographviz.Shape):       "circle",
		string(gographviz.ColorScheme): "svg",
		string(gographviz.Color):       "darkslategray",
		string(gographviz.FillColor):   "darkslategray",
		string(gographviz.FontColor):   "white",
	}
	defaultEdgeAttrs = map[string]string{
		string(gographviz.Style): "solid",
		string(gographviz.Color): "black",
	}
)

func NewDOTGraph(g *graph.Graph) (*gographviz.Graph, error) {
	gv := gographviz.NewGraph()

	if err := gv.SetDir(g.IsDirected()); err != nil {
		return nil, err
	}
	for k, v := range defaultGraphAttrs {
		if err := gv.AddAttr(gv.Name, k, v); err != nil {
			return nil, err
		}
	}

	// add nodes
	nodes := g.GetNodes()
	for _, n := range nodes {
		if err := gv.AddNode(gv.Name, n.ID().String(), defaultNodeAttrs); err != nil {
			return nil, err
		}
	}

	// add edges
	edges := g.GetEdges()
	for _, nodeEdges := range edges {
		for _, e := range nodeEdges {
			if err := gv.AddEdge(
				e.Tail().ID().String(),
				e.Head().ID().String(),
				e.IsDirected(),
				defaultEdgeAttrs,
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
