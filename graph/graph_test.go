package graph

import (
	"sort"
	"testing"
)

func testInitialized(t *testing.T, g *Graph) {
	if len(g.nodes) > 0 {
		t.Errorf("expected: %d, actual: %d", 0, len(g.nodes))
	}
	if len(g.heads) > 0 {
		t.Errorf("expected: %d, actual: %d", 0, len(g.heads))
	}
	if len(g.tails) > 0 {
		t.Errorf("expected: %d, actual: %d", 0, len(g.tails))
	}
	if len(g.edges) > 0 {
		t.Errorf("expected: %d, actual: %d", 0, len(g.edges))
	}
}

func testGraphEquality(t *testing.T, expected, actual *Graph) {
	if actual.isDirected != expected.isDirected {
		t.Errorf("expected: %t, actual: %t", expected.isDirected, actual.isDirected)
	}
	testNodesEquality(t, expected.nodes, actual.nodes)
	testEndsEquality(t, expected.heads, actual.heads)
	testEndsEquality(t, expected.tails, actual.tails)
	testEdgesEquality(t, expected.edges, actual.edges)
}

func testNodesEquality(t *testing.T, expected, actual map[ID]*Node) {
	if len(actual) != len(expected) {
		t.Errorf("expected: %d, actual: %d", len(expected), len(actual))
	}

	for id, nExpected := range expected {
		if _, ok := actual[id]; !ok {
			t.Errorf("expected: %t, actual: %t", true, ok)
			continue
		}

		nActual := actual[id]
		testNodeEquality(t, nExpected, nActual)
	}
}

func testEndsEquality(t *testing.T, expected, actual map[ID]map[ID]*Node) {
	if len(actual) != len(expected) {
		t.Errorf("expected: %d, actual: %d", len(expected), len(actual))
	}

	for id1, endsExpected := range expected {
		if _, ok := actual[id1]; !ok {
			t.Errorf("expected: %t, actual: %t", true, ok)
			continue
		}

		endsActual := actual[id1]

		if len(endsActual) != len(endsExpected) {
			t.Errorf("expected: %d, actual: %d", len(endsExpected), len(endsActual))
			continue
		}

		for id2, nExpected := range endsExpected {
			if _, ok := endsActual[id2]; !ok {
				t.Errorf("expected: %t, actual: %t", true, ok)
				continue
			}

			nActual := endsActual[id2]
			testNodeEquality(t, nExpected, nActual)
		}
	}
}

func testEdgesEquality(t *testing.T, expected, actual map[ID]map[ID]*Edge) {
	if len(actual) != len(expected) {
		t.Errorf("expected: %d, actual: %d", len(expected), len(actual))
	}

	for id1, endsExpected := range expected {
		if _, ok := actual[id1]; !ok {
			t.Errorf("expected: %t, actual: %t", true, ok)
			continue
		}

		endsActual := actual[id1]

		if len(endsActual) != len(endsExpected) {
			t.Errorf("expected: %d, actual: %d", len(endsExpected), len(endsActual))
			continue
		}

		for id2, eExpected := range endsExpected {
			if _, ok := endsActual[id2]; !ok {
				t.Errorf("expected: %t, actual: %t", true, ok)
				continue
			}

			eActual := endsActual[id2]
			testEdgeEquality(t, eExpected, eActual)
		}
	}
}

func TestNewDirected(t *testing.T) {
	g := NewDirected()
	if !g.isDirected {
		t.Errorf("expected: %t, actual: %t", true, g.isDirected)
	}
	testInitialized(t, g)
}

func TestNewUndirecred(t *testing.T) {
	g := NewUndirected()
	if g.isDirected {
		t.Errorf("expected: %t, actual: %t", false, g.isDirected)
	}
	testInitialized(t, g)
}

func TestGraph_IsDirected(t *testing.T) {
	g := &Graph{
		isDirected: true,
	}

	t.Run("directed", func(t *testing.T) {
		if !g.IsDirected() {
			t.Errorf("expected: %t, actual: %t", true, g.IsDirected())
		}
	})

	g.isDirected = false

	t.Run("undirected", func(t *testing.T) {
		if g.IsDirected() {
			t.Errorf("expected: %t, actual: %t", false, g.IsDirected())
		}
	})
}

func TestGraph_GetNode(t *testing.T) {
	type input struct {
		id ID
	}
	type output struct {
		node *Node
		err  error
	}

	n1 := newTestNode("1")
	n2 := newTestNode("2")

	testCases := []struct {
		name  string
		graph *Graph
		in    input
		out   output
	}{
		{
			"success",
			&Graph{nodes: map[ID]*Node{n1.id: n1}},
			input{n1.id},
			output{n1, nil},
		},
		{
			"failure: non-existent node",
			&Graph{nodes: map[ID]*Node{n1.id: n1}},
			input{n2.id},
			output{nil, ErrNodeNotExist},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, in, out := tc.graph, tc.in, tc.out

			n, err := g.GetNode(in.id)
			if err != out.err {
				t.Errorf("expected: %v, actual: %v", out.err, err)
				return
			}
			if out.node == nil {
				if n != nil {
					t.Errorf("expected: nil, actual: non-nil")
				}
			} else {
				if n == nil {
					t.Errorf("expected: non-nil, actual: nil")
					return
				}
				testNodeEquality(t, out.node, n)
			}
		})
	}
}

func TestGraph_GetNodes(t *testing.T) {
	n1 := newTestNode("1")
	n2 := newTestNode("2")
	expected := map[ID]*Node{n1.id: n1, n2.id: n2}
	g := &Graph{
		nodes: expected,
	}

	actual := g.GetNodes()
	testNodesEquality(t, expected, actual)
}

func TestGraph_GetTails(t *testing.T) {
	type input struct {
		id ID
	}
	type output struct {
		nodes map[ID]*Node
		err   error
	}

	n1 := newTestNode("1")
	n2 := newTestNode("2")
	n3 := newTestNode("3")

	testCases := []struct {
		name  string
		graph *Graph
		in    input
		out   output
	}{
		{
			"success: empty",
			&Graph{
				nodes: map[ID]*Node{n1.id: n1},
				tails: map[ID]map[ID]*Node{},
			},
			input{n1.id},
			output{map[ID]*Node{}, nil},
		},
		{
			"success",
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				tails: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
				},
			},
			input{n1.id},
			output{map[ID]*Node{n2.id: n2, n3.id: n3}, nil},
		},
		{
			"failure: non-existent node",
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2},
				tails: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
				},
			},
			input{n3.id},
			output{nil, ErrNodeNotExist},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, in, out := tc.graph, tc.in, tc.out

			nodes, err := g.GetTails(in.id)
			if err != out.err {
				t.Errorf("expected: %v, actual: %v", out.err, err)
			}
			testNodesEquality(t, out.nodes, nodes)
		})
	}
}

func TestGraph_GetHeads(t *testing.T) {
	type input struct {
		id ID
	}
	type output struct {
		nodes map[ID]*Node
		err   error
	}

	n1 := newTestNode("1")
	n2 := newTestNode("2")
	n3 := newTestNode("3")

	testCases := []struct {
		name  string
		graph *Graph
		in    input
		out   output
	}{
		{
			"success: empty",
			&Graph{
				nodes: map[ID]*Node{n1.id: n1},
				heads: map[ID]map[ID]*Node{},
			},
			input{n1.id},
			output{map[ID]*Node{}, nil},
		},
		{
			"success",
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
				},
			},
			input{n1.id},
			output{map[ID]*Node{n2.id: n2, n3.id: n3}, nil},
		},
		{
			"failure: non-existent node",
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
				},
			},
			input{n3.id},
			output{nil, ErrNodeNotExist},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, in, out := tc.graph, tc.in, tc.out

			nodes, err := g.GetHeads(in.id)
			if err != out.err {
				t.Errorf("expected: %v, actual: %v", out.err, err)
			}
			testNodesEquality(t, out.nodes, nodes)
		})
	}
}

func TestGraph_AddNode(t *testing.T) {
	type input struct {
		node *Node
	}
	type output struct {
		err error
	}

	n1 := newTestNode("1")

	testCases := []struct {
		name     string
		actual   *Graph
		expected *Graph
		in       input
		out      output
	}{
		{
			"success",
			&Graph{nodes: map[ID]*Node{}},
			&Graph{nodes: map[ID]*Node{n1.id: n1}},
			input{n1},
			output{nil},
		},
		{
			"success: existent node",
			&Graph{nodes: map[ID]*Node{n1.id: n1}},
			&Graph{nodes: map[ID]*Node{n1.id: n1}},
			input{n1},
			output{nil},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, expected, in, out := tc.actual, tc.expected, tc.in, tc.out

			if err := actual.AddNode(in.node); err != out.err {
				t.Errorf("expected: %v, actual: %v", out.err, err)
			}
			testGraphEquality(t, expected, actual)
		})
	}
}

func TestGraph_RemoveNode(t *testing.T) {
	type input struct {
		id ID
	}
	type output struct {
		err error
	}

	n1 := newTestNode("1")
	n2 := newTestNode("2")
	n3 := newTestNode("3")
	e12 := newTestEdgeGenerator("1", "2")
	e13 := newTestEdgeGenerator("1", "3")
	e21 := newTestEdgeGenerator("2", "1")
	e23 := newTestEdgeGenerator("2", "3")
	e31 := newTestEdgeGenerator("3", "1")
	e32 := newTestEdgeGenerator("3", "2")

	testCases := []struct {
		name     string
		actual   *Graph
		expected *Graph
		in       input
		out      output
	}{
		{
			"success",
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
					n2.id: {n1.id: n1, n3.id: n3},
					n3.id: {n1.id: n1, n2.id: n2},
				},
				tails: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
					n2.id: {n1.id: n1, n3.id: n3},
					n3.id: {n1.id: n1, n2.id: n2},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1), n3.id: e13(true, 1)},
					n2.id: {n1.id: e21(true, 1), n3.id: e23(true, 1)},
					n3.id: {n1.id: e31(true, 1), n2.id: e32(true, 1)},
				},
			},
			&Graph{
				nodes: map[ID]*Node{n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n2.id: {n3.id: n3},
					n3.id: {n2.id: n2},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n3.id: n3},
					n3.id: {n2.id: n2},
				},
				edges: map[ID]map[ID]*Edge{
					n2.id: {n3.id: e23(true, 1)},
					n3.id: {n2.id: e32(true, 1)},
				},
			},
			input{n1.id},
			output{nil},
		},
		{
			"success: non-existent node",
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1)},
				},
			},
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1)},
				},
			},
			input{n3.id},
			output{nil},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, expected, in, out := tc.actual, tc.expected, tc.in, tc.out

			if err := tc.actual.RemoveNode(in.id); err != out.err {
				t.Errorf("expected: %v, actual: %v", out.err, err)
			}
			testGraphEquality(t, expected, actual)
		})
	}
}

func TestGraph_GetEdge(t *testing.T) {
	type input struct {
		idTail ID
		idHead ID
	}
	type output struct {
		edge *Edge
		err  error
	}

	n1 := newTestNode("1")
	n2 := newTestNode("2")
	n3 := newTestNode("3")
	e12 := newTestEdgeGenerator("1", "2")
	e13 := newTestEdgeGenerator("1", "3")

	testCases := []struct {
		name  string
		graph *Graph
		in    input
		out   output
	}{
		{
			"success",
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
					n3.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1), n3.id: e13(true, 1)},
				},
			},
			input{n1.id, n2.id},
			output{e12(true, 1), nil},
		},
		{
			"failure: non-existent tail node",
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
					n3.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1), n3.id: e13(true, 1)},
				},
			},
			input{n3.id, n2.id},
			output{nil, ErrNodeNotExist},
		},
		{
			"failure: non-existent head node",
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
					n3.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1), n3.id: e13(true, 1)},
				},
			},
			input{n1.id, n3.id},
			output{nil, ErrNodeNotExist},
		},
		{
			"failure: non-existent edge",
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
					n3.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1), n3.id: e13(true, 1)},
				},
			},
			input{n2.id, n3.id},
			output{nil, ErrEdgeNotExist},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, in, out := tc.graph, tc.in, tc.out

			e, err := g.GetEdge(in.idTail, in.idHead)
			if err != out.err {
				t.Errorf("expected: %v, actual: %v", out.err, err)
				return
			}
			if out.edge == nil {
				if e != nil {
					t.Errorf("expected: nil, actual: non-nil")
				}
			} else {
				if e == nil {
					t.Errorf("expected: non-nil, actual: nil")
					return
				}
				testEdgeEquality(t, out.edge, e)
			}
		})
	}
}

func TestGraph_GetEdges(t *testing.T) {
	n1 := newTestNode("1")
	n2 := newTestNode("2")
	n3 := newTestNode("3")
	expected := map[ID]map[ID]*Edge{
		n1.id: {n2.id: newTestEdge(true, "1", "2", 1.2)},
		n2.id: {n3.id: newTestEdge(true, "2", "3", 2.3)},
	}
	g := &Graph{
		edges: expected,
	}

	actual := g.GetEdges()
	testEdgesEquality(t, expected, actual)
}

func TestGraph_AddEdge(t *testing.T) {
	type input struct {
		idTail ID
		idHead ID
		weight float64
	}
	type output struct {
		err error
	}

	n1 := newTestNode("1")
	n2 := newTestNode("2")
	n3 := newTestNode("3")
	e12 := newTestEdgeGenerator("1", "2")
	e13 := newTestEdgeGenerator("1", "3")
	e21 := newTestEdgeGenerator("2", "1")
	e23 := newTestEdgeGenerator("2", "3")
	e31 := newTestEdgeGenerator("3", "1")
	e32 := newTestEdgeGenerator("3", "2")

	testCases := []struct {
		name     string
		actual   *Graph
		expected *Graph
		in       input
		out      output
	}{
		{
			"success: directed, first edge",
			&Graph{
				isDirected: true,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads:      map[ID]map[ID]*Node{},
				tails:      map[ID]map[ID]*Node{},
				edges:      map[ID]map[ID]*Edge{},
			},
			&Graph{
				isDirected: true,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1)},
				},
			},
			input{n1.id, n2.id, 1.0},
			output{nil},
		},
		{
			"success: directed, second edge",
			&Graph{
				isDirected: true,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
					n2.id: {n3.id: n3},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
					n3.id: {n2.id: n2},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1)},
					n2.id: {n3.id: e23(true, 1)},
				},
			},
			&Graph{
				isDirected: true,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
					n2.id: {n3.id: n3},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
					n3.id: {n1.id: n1, n2.id: n2},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1), n3.id: e13(true, 1)},
					n2.id: {n3.id: e23(true, 1)},
				},
			},
			input{n1.id, n3.id, 1.0},
			output{nil},
		},
		{
			"success: directed, reversed edge",
			&Graph{
				isDirected: true,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
					n2.id: {n3.id: n3},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
					n3.id: {n2.id: n2},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1)},
					n2.id: {n3.id: e23(true, 1)},
				},
			},
			&Graph{
				isDirected: true,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
					n2.id: {n1.id: n1, n3.id: n3},
				},
				tails: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
					n2.id: {n1.id: n1},
					n3.id: {n2.id: n2},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1)},
					n2.id: {n1.id: e21(true, 1), n3.id: e23(true, 1)},
				},
			},
			input{n2.id, n1.id, 1.0},
			output{nil},
		},
		{
			"success: directed, existent edge",
			&Graph{
				isDirected: true,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
					n3.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1), n3.id: e13(true, 1)},
				},
			},
			&Graph{
				isDirected: true,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
					n3.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 2), n3.id: e13(true, 1)},
				},
			},
			input{n1.id, n2.id, 1.0},
			output{nil},
		},
		{
			"success: undirected, first edge",
			&Graph{
				isDirected: false,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads:      map[ID]map[ID]*Node{},
				tails:      map[ID]map[ID]*Node{},
				edges:      map[ID]map[ID]*Edge{},
			},
			&Graph{
				isDirected: false,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
					n2.id: {n1.id: n1},
				},
				tails: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
					n2.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(false, 1)},
					n2.id: {n1.id: e21(false, 1)},
				},
			},
			input{n1.id, n2.id, 1.0},
			output{nil},
		},
		{
			"success: undirected, second edge",
			&Graph{
				isDirected: false,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
					n2.id: {n1.id: n1, n3.id: n3},
					n3.id: {n2.id: n2},
				},
				tails: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
					n2.id: {n1.id: n1, n3.id: n3},
					n3.id: {n2.id: n2},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(false, 1)},
					n2.id: {n1.id: e21(false, 1), n3.id: e23(false, 1)},
					n3.id: {n2.id: e32(false, 1)},
				},
			},
			&Graph{
				isDirected: false,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
					n2.id: {n1.id: n1, n3.id: n3},
					n3.id: {n1.id: n1, n2.id: n2},
				},
				tails: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
					n2.id: {n1.id: n1, n3.id: n3},
					n3.id: {n1.id: n1, n2.id: n2},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(false, 1), n3.id: e13(false, 1)},
					n2.id: {n1.id: e21(false, 1), n3.id: e23(false, 1)},
					n3.id: {n1.id: e31(false, 1), n2.id: e32(false, 1)},
				},
			},
			input{n1.id, n3.id, 1.0},
			output{nil},
		},
		{
			"success: undirected, existent edge",
			&Graph{
				isDirected: false,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
					n2.id: {n1.id: n1},
					n3.id: {n1.id: n1},
				},
				tails: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
					n2.id: {n1.id: n1},
					n3.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(false, 1), n3.id: e13(false, 1)},
					n2.id: {n1.id: e21(false, 1)},
					n3.id: {n1.id: e31(false, 1)},
				},
			},
			&Graph{
				isDirected: false,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
					n2.id: {n1.id: n1},
					n3.id: {n1.id: n1},
				},
				tails: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
					n2.id: {n1.id: n1},
					n3.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(false, 2), n3.id: e13(false, 1)},
					n2.id: {n1.id: e21(false, 2)},
					n3.id: {n1.id: e31(false, 1)},
				},
			},
			input{n1.id, n2.id, 1.0},
			output{nil},
		},
		{
			"failure: looped edge",
			&Graph{
				nodes: map[ID]*Node{n1.id: n1},
				heads: map[ID]map[ID]*Node{},
				tails: map[ID]map[ID]*Node{},
			},
			&Graph{
				nodes: map[ID]*Node{n1.id: n1},
				heads: map[ID]map[ID]*Node{},
				tails: map[ID]map[ID]*Node{},
			},
			input{n1.id, n1.id, 1.0},
			output{ErrEdgeLooped},
		},
		{
			"failure: non-existent tail node",
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2},
				heads: map[ID]map[ID]*Node{},
				tails: map[ID]map[ID]*Node{},
			},
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2},
				heads: map[ID]map[ID]*Node{},
				tails: map[ID]map[ID]*Node{},
			},
			input{n3.id, n2.id, 1.0},
			output{ErrNodeNotExist},
		},
		{
			"failure: non-existent head node",
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2},
				heads: map[ID]map[ID]*Node{},
				tails: map[ID]map[ID]*Node{},
			},
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2},
				heads: map[ID]map[ID]*Node{},
				tails: map[ID]map[ID]*Node{},
			},
			input{n1.id, n3.id, 1.0},
			output{ErrNodeNotExist},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, expected, in, out := tc.actual, tc.expected, tc.in, tc.out

			if err := actual.AddEdge(in.idTail, in.idHead, in.weight); err != out.err {
				t.Errorf("expected: %v, actual: %v", out.err, err)
			}
			testGraphEquality(t, expected, actual)
		})
	}
}

func TestGraph_RemoveEdge(t *testing.T) {
	type input struct {
		idTail ID
		idHead ID
	}
	type output struct {
		err error
	}

	n1 := newTestNode("1")
	n2 := newTestNode("2")
	n3 := newTestNode("3")
	e12 := newTestEdgeGenerator("1", "2")
	e21 := newTestEdgeGenerator("2", "1")
	e13 := newTestEdgeGenerator("1", "3")
	e31 := newTestEdgeGenerator("3", "1")

	testCases := []struct {
		name     string
		actual   *Graph
		expected *Graph
		in       input
		out      output
	}{
		{
			"success: directed",
			&Graph{
				isDirected: true,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
					n3.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1), n3.id: e13(true, 1)},
				},
			},
			&Graph{
				isDirected: true,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n3.id: n3},
				},
				tails: map[ID]map[ID]*Node{
					n3.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n3.id: e13(true, 1)},
				},
			},
			input{n1.id, n2.id},
			output{nil},
		},
		{
			"success: directed, bidirectional edge",
			&Graph{
				isDirected: true,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
					n2.id: {n1.id: n1},
				},
				tails: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
					n2.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1)},
					n2.id: {n1.id: e21(true, 1)},
				},
			},
			&Graph{
				isDirected: true,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2},
				heads: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
				},
				tails: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
				},
				edges: map[ID]map[ID]*Edge{
					n2.id: {n1.id: e21(true, 1)},
				},
			},
			input{n1.id, n2.id},
			output{nil},
		},
		{
			"success: directed, non-existent edge",
			&Graph{
				isDirected: true,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1)},
				},
			},
			&Graph{
				isDirected: true,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1)},
				},
			},
			input{n1.id, n3.id},
			output{nil},
		},
		{
			"success: undirected",
			&Graph{
				isDirected: false,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
					n2.id: {n1.id: n1},
					n3.id: {n1.id: n1},
				},
				tails: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2, n3.id: n3},
					n2.id: {n1.id: n1},
					n3.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(false, 1), n3.id: e13(false, 1)},
					n2.id: {n1.id: e21(false, 1)},
					n3.id: {n1.id: e31(false, 1)},
				},
			},
			&Graph{
				isDirected: false,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n3.id: n3},
					n3.id: {n1.id: n1},
				},
				tails: map[ID]map[ID]*Node{
					n1.id: {n3.id: n3},
					n3.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n3.id: e13(false, 1)},
					n3.id: {n1.id: e31(false, 1)},
				},
			},
			input{n1.id, n2.id},
			output{nil},
		},
		{
			"success: undirected, non-existent edge",
			&Graph{
				isDirected: false,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
					n2.id: {n1.id: n1},
				},
				tails: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
					n2.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(false, 1)},
					n2.id: {n1.id: e21(false, 1)},
				},
			},
			&Graph{
				isDirected: false,
				nodes:      map[ID]*Node{n1.id: n1, n2.id: n2, n3.id: n3},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
					n2.id: {n1.id: n1},
				},
				tails: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
					n2.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(false, 1)},
					n2.id: {n1.id: e21(false, 1)},
				},
			},
			input{n1.id, n3.id},
			output{nil},
		},
		{
			"failure: non-existent tail node",
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1)},
				},
			},
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1)},
				},
			},
			input{n3.id, n2.id},
			output{ErrNodeNotExist},
		},
		{
			"failure: non-existent head node",
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1)},
				},
			},
			&Graph{
				nodes: map[ID]*Node{n1.id: n1, n2.id: n2},
				heads: map[ID]map[ID]*Node{
					n1.id: {n2.id: n2},
				},
				tails: map[ID]map[ID]*Node{
					n2.id: {n1.id: n1},
				},
				edges: map[ID]map[ID]*Edge{
					n1.id: {n2.id: e12(true, 1)},
				},
			},
			input{n1.id, n3.id},
			output{ErrNodeNotExist},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, expected, in, out := tc.actual, tc.expected, tc.in, tc.out

			if err := actual.RemoveEdge(in.idTail, in.idHead); err != out.err {
				t.Errorf("expected: %v, actual: %v", out.err, err)
			}
			testGraphEquality(t, expected, actual)
		})
	}
}

func TestGraph_GetIndegreeDistribution(t *testing.T) {
	n1 := newTestNode("1")
	n2 := newTestNode("2")
	n3 := newTestNode("3")
	n4 := newTestNode("4")

	expected := &DegreeDistribution{
		m:       map[int]int{0: 1, 1: 2, 2: 1},
		degrees: []int{0, 1, 2},
	}
	g := &Graph{
		tails: map[ID]map[ID]*Node{
			n1.id: {n2.id: n2, n3.id: n3},
			n2.id: {n4.id: n4},
			n3.id: {n4.id: n4},
			n4.id: {},
		},
	}

	actual := g.GetIndegreeDistribution()
	sort.Sort(actual)
	testDegreeDistributionEquality(t, expected, actual)
}

func TestGraph_GetOutdegreeDistribution(t *testing.T) {
	n1 := newTestNode("1")
	n2 := newTestNode("2")
	n3 := newTestNode("3")
	n4 := newTestNode("4")

	expected := &DegreeDistribution{
		m:       map[int]int{0: 1, 1: 2, 2: 1},
		degrees: []int{0, 1, 2},
	}
	g := &Graph{
		heads: map[ID]map[ID]*Node{
			n1.id: {n2.id: n2, n3.id: n3},
			n2.id: {n4.id: n4},
			n3.id: {n4.id: n4},
			n4.id: {},
		},
	}

	actual := g.GetOutdegreeDistribution()
	sort.Sort(actual)
	testDegreeDistributionEquality(t, expected, actual)
}
