package graph

import "testing"

func testInitialized(t *testing.T, g *graph) {
	if len(g.idToNodes) > 0 {
		t.Errorf("expected: %d, actual: %d", 0, len(g.idToNodes))
	}
	if len(g.idToTails) > 0 {
		t.Errorf("expected: %d, actual: %d", 0, len(g.idToTails))
	}
	if len(g.idToHeads) > 0 {
		t.Errorf("expected: %d, actual: %d", 0, len(g.idToHeads))
	}
}

func testGraphEquality(t *testing.T, expected, actual *graph) {
	if actual.isDirected != expected.isDirected {
		t.Errorf("expected: %t, actual: %t", expected.isDirected, actual.isDirected)
	}
	testIDToNodesEquality(t, expected.idToNodes, actual.idToNodes)
	testIDToEndsEquality(t, expected.idToTails, actual.idToTails)
	testIDToEndsEquality(t, expected.idToHeads, actual.idToHeads)
}

func testIDToNodesEquality(t *testing.T, expected, actual map[ID]Node) {
	if len(actual) != len(expected) {
		t.Errorf("expected: %d, actual: %d", len(expected), len(actual))
	}

	for id, n := range expected {
		if _, ok := actual[id]; !ok {
			t.Errorf("expected: %t, actual: %t", true, ok)
			continue
		}

		if actual[id].ID() != n.ID() {
			t.Errorf("expected: %q, actual: %q", n.ID(), actual[id].ID())
		}
	}
}

func testIDToEndsEquality(t *testing.T, expected, actual map[ID]map[ID]Edge) {
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
	g := NewDirected().(*graph)
	if !g.isDirected {
		t.Errorf("expected: %t, actual: %t", true, g.isDirected)
	}
	testInitialized(t, g)
}

func TestNewUndirecred(t *testing.T) {
	g := NewUndirected().(*graph)
	if g.isDirected {
		t.Errorf("expected: %t, actual: %t", false, g.isDirected)
	}
	testInitialized(t, g)
}

func TestGraph_IsDirected(t *testing.T) {
	g := &graph{
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
		id StringID
	}
	type output struct {
		node Node
		err  error
	}

	n1 := newTestNode("1")
	n2 := newTestNode("2")

	testCases := []struct {
		name  string
		graph *graph
		in    input
		out   output
	}{
		{
			"success",
			&graph{idToNodes: map[ID]Node{n1.id: n1}},
			input{n1.id},
			output{n1, nil},
		},
		{
			"failure: non-existent node",
			&graph{idToNodes: map[ID]Node{n1.id: n1}},
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
			}
			if out.node == nil {
				if n != nil {
					t.Errorf("expected: nil, actual: non-nil")
				}
			} else {
				if n == nil {
					t.Errorf("expected: non-nil, actual: nil")
				}
				if n.ID() != in.id {
					t.Errorf("expected: %q, actual: %q", in.id, n.ID())
				}
			}
		})
	}
}

func TestGraph_GetNodes(t *testing.T) {
	n1 := newTestNode("1")
	n2 := newTestNode("2")
	expected := map[ID]Node{n1.id: n1, n2.id: n2}
	g := &graph{
		idToNodes: expected,
	}

	actual, err := g.GetNodes()
	if err != nil {
		t.Fatal(err)
	}
	testIDToNodesEquality(t, expected, actual)
}

func TestGraph_GetTails(t *testing.T) {
	type input struct {
		id StringID
	}
	type output struct {
		idToNodes map[ID]Node
		err       error
	}

	n1 := newTestNode("1")
	n2 := newTestNode("2")
	n3 := newTestNode("3")
	e21 := newTestEdgeGenerator("2", "1")
	e31 := newTestEdgeGenerator("3", "1")

	testCases := []struct {
		name  string
		graph *graph
		in    input
		out   output
	}{
		{
			"success: empty",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1},
				idToTails: map[ID]map[ID]Edge{},
			},
			input{n1.id},
			output{map[ID]Node{}, nil},
		},
		{
			"success",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n1.id: {n2.id: e21(true, 1), n3.id: e31(true, 1)},
				},
			},
			input{n1.id},
			output{map[ID]Node{n2.id: n2, n3.id: n3}, nil},
		},
		{
			"failure: non-existent node",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]Edge{
					n1.id: {n2.id: e21(true, 1)},
				},
			},
			input{n3.id},
			output{nil, ErrNodeNotExist},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, in, out := tc.graph, tc.in, tc.out

			idToNodes, err := g.GetTails(in.id)
			if err != out.err {
				t.Errorf("expected: %v, actual: %v", out.err, err)
			}
			testIDToNodesEquality(t, out.idToNodes, idToNodes)
		})
	}
}

func TestGraph_GetHeads(t *testing.T) {
	type input struct {
		id StringID
	}
	type output struct {
		idToNodes map[ID]Node
		err       error
	}

	n1 := newTestNode("1")
	n2 := newTestNode("2")
	n3 := newTestNode("3")
	e12 := newTestEdgeGenerator("1", "2")
	e13 := newTestEdgeGenerator("1", "3")

	testCases := []struct {
		name  string
		graph *graph
		in    input
		out   output
	}{
		{
			"success: empty",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1},
				idToHeads: map[ID]map[ID]Edge{},
			},
			input{n1.id},
			output{map[ID]Node{}, nil},
		},
		{
			"success",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 1), n3.id: e13(true, 1)},
				},
			},
			input{n1.id},
			output{map[ID]Node{n2.id: n2, n3.id: n3}, nil},
		},
		{
			"failure: non-existent node",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 1)},
				},
			},
			input{n3.id},
			output{nil, ErrNodeNotExist},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, in, out := tc.graph, tc.in, tc.out

			idToNodes, err := g.GetHeads(in.id)
			if err != out.err {
				t.Errorf("expected: %v, actual: %v", out.err, err)
			}
			testIDToNodesEquality(t, out.idToNodes, idToNodes)
		})
	}
}

func TestGraph_AddNode(t *testing.T) {
	type input struct {
		node *node
	}
	type output struct {
		err error
	}

	n1 := newTestNode("1")

	testCases := []struct {
		name     string
		actual   *graph
		expected *graph
		in       input
		out      output
	}{
		{
			"success",
			&graph{idToNodes: map[ID]Node{}},
			&graph{idToNodes: map[ID]Node{n1.id: n1}},
			input{n1},
			output{nil},
		},
		{
			"success: existent node",
			&graph{idToNodes: map[ID]Node{n1.id: n1}},
			&graph{idToNodes: map[ID]Node{n1.id: n1}},
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
		id StringID
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
		actual   *graph
		expected *graph
		in       input
		out      output
	}{
		{
			"success",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n1.id: {n2.id: e21(true, 1), n3.id: e31(true, 1)},
					n2.id: {n1.id: e12(true, 1), n3.id: e32(true, 1)},
					n3.id: {n1.id: e13(true, 1), n2.id: e23(true, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 1), n3.id: e13(true, 1)},
					n2.id: {n1.id: e21(true, 1), n3.id: e23(true, 1)},
					n3.id: {n1.id: e31(true, 1), n2.id: e32(true, 1)},
				},
			},
			&graph{
				idToNodes: map[ID]Node{n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n3.id: e32(true, 1)},
					n3.id: {n2.id: e23(true, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n2.id: {n3.id: e23(true, 1)},
					n3.id: {n2.id: e32(true, 1)},
				},
			},
			input{n1.id},
			output{nil},
		},
		{
			"success: non-existent node",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n1.id: e12(true, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 1)},
				},
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n1.id: e12(true, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
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

func TestGraph_AddEdge(t *testing.T) {
	type input struct {
		idTail StringID
		idHead StringID
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
		actual   *graph
		expected *graph
		in       input
		out      output
	}{
		{
			"success: directed, first edge",
			&graph{
				isDirected: true,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails:  map[ID]map[ID]Edge{},
				idToHeads:  map[ID]map[ID]Edge{},
			},
			&graph{
				isDirected: true,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n1.id: e12(true, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 1)},
				},
			},
			input{n1.id, n2.id, 1.0},
			output{nil},
		},
		{
			"success: directed, second edge",
			&graph{
				isDirected: true,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n1.id: e12(true, 1)},
					n3.id: {n2.id: e23(true, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 1)},
					n2.id: {n3.id: e23(true, 1)},
				},
			},
			&graph{
				isDirected: true,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n1.id: e12(true, 1)},
					n3.id: {n1.id: e13(true, 1), n2.id: e23(true, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 1), n3.id: e13(true, 1)},
					n2.id: {n3.id: e23(true, 1)},
				},
			},
			input{n1.id, n3.id, 1.0},
			output{nil},
		},
		{
			"success: directed, existent edge",
			&graph{
				isDirected: true,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n1.id: e12(true, 1)},
					n3.id: {n1.id: e13(true, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 1), n3.id: e13(true, 1)},
				},
			},
			&graph{
				isDirected: true,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n1.id: e12(true, 2)},
					n3.id: {n1.id: e13(true, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 2), n3.id: e13(true, 1)},
				},
			},
			input{n1.id, n2.id, 1.0},
			output{nil},
		},
		{
			"success: undirected, first edge",
			&graph{
				isDirected: false,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails:  map[ID]map[ID]Edge{},
				idToHeads:  map[ID]map[ID]Edge{},
			},
			&graph{
				isDirected: false,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n1.id: {n2.id: e21(false, 1)},
					n2.id: {n1.id: e12(false, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(false, 1)},
					n2.id: {n1.id: e21(false, 1)},
				},
			},
			input{n1.id, n2.id, 1.0},
			output{nil},
		},
		{
			"success: undirected, second edge",
			&graph{
				isDirected: false,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n1.id: {n2.id: e21(false, 1)},
					n2.id: {n1.id: e12(false, 1), n3.id: e32(false, 1)},
					n3.id: {n2.id: e23(false, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(false, 1)},
					n2.id: {n1.id: e21(false, 1), n3.id: e23(false, 1)},
					n3.id: {n2.id: e32(false, 1)},
				},
			},
			&graph{
				isDirected: false,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n1.id: {n2.id: e21(false, 1), n3.id: e31(false, 1)},
					n2.id: {n1.id: e12(false, 1), n3.id: e32(false, 1)},
					n3.id: {n1.id: e13(false, 1), n2.id: e23(false, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
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
			&graph{
				isDirected: false,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n1.id: {n2.id: e21(false, 1), n3.id: e31(false, 1)},
					n2.id: {n1.id: e12(false, 1)},
					n3.id: {n1.id: e13(false, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(false, 1), n3.id: e13(false, 1)},
					n2.id: {n1.id: e21(false, 1)},
					n3.id: {n1.id: e31(false, 1)},
				},
			},
			&graph{
				isDirected: false,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n1.id: {n2.id: e21(false, 2), n3.id: e31(false, 1)},
					n2.id: {n1.id: e12(false, 2)},
					n3.id: {n1.id: e13(false, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(false, 2), n3.id: e13(false, 1)},
					n2.id: {n1.id: e21(false, 2)},
					n3.id: {n1.id: e31(false, 1)},
				},
			},
			input{n1.id, n2.id, 1.0},
			output{nil},
		},
		{
			"failure: non-existent tail node",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]Edge{},
				idToHeads: map[ID]map[ID]Edge{},
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]Edge{},
				idToHeads: map[ID]map[ID]Edge{},
			},
			input{n3.id, n2.id, 1.0},
			output{ErrNodeNotExist},
		},
		{
			"failure: non-existent head node",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]Edge{},
				idToHeads: map[ID]map[ID]Edge{},
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]Edge{},
				idToHeads: map[ID]map[ID]Edge{},
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
		idTail StringID
		idHead StringID
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
	e31 := newTestEdgeGenerator("3", "1")

	testCases := []struct {
		name     string
		actual   *graph
		expected *graph
		in       input
		out      output
	}{
		{
			"success: directed",
			&graph{
				isDirected: true,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n1.id: e12(true, 1)},
					n3.id: {n1.id: e13(true, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 1), n3.id: e13(true, 1)},
				},
			},
			&graph{
				isDirected: true,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n3.id: {n1.id: e13(true, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n3.id: e13(true, 1)},
				},
			},
			input{n1.id, n2.id},
			output{nil},
		},
		{
			"success: directed, non-existent edge",
			&graph{
				isDirected: true,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n1.id: e12(true, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 1)},
				},
			},
			&graph{
				isDirected: true,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n1.id: e12(true, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 1)},
				},
			},
			input{n1.id, n3.id},
			output{nil},
		},
		{
			"success: undirected",
			&graph{
				isDirected: false,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n1.id: {n2.id: e21(false, 1), n3.id: e31(false, 1)},
					n2.id: {n1.id: e12(false, 1)},
					n3.id: {n1.id: e13(false, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(false, 1), n3.id: e13(false, 1)},
					n2.id: {n1.id: e21(false, 1)},
					n3.id: {n1.id: e31(false, 1)},
				},
			},
			&graph{
				isDirected: false,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n1.id: {n3.id: e31(false, 1)},
					n3.id: {n1.id: e13(false, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n3.id: e13(false, 1)},
					n3.id: {n1.id: e31(false, 1)},
				},
			},
			input{n1.id, n2.id},
			output{nil},
		},
		{
			"success: undirected, non-existent edge",
			&graph{
				isDirected: false,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n1.id: {n2.id: e21(false, 1)},
					n2.id: {n1.id: e12(false, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(false, 1)},
					n2.id: {n1.id: e21(false, 1)},
				},
			},
			&graph{
				isDirected: false,
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n1.id: {n2.id: e21(false, 1)},
					n2.id: {n1.id: e12(false, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(false, 1)},
					n2.id: {n1.id: e21(false, 1)},
				},
			},
			input{n1.id, n3.id},
			output{nil},
		},
		{
			"failure: non-existent tail node",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n1.id: e12(true, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 1)},
				},
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n1.id: e12(true, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 1)},
				},
			},
			input{n3.id, n2.id},
			output{ErrNodeNotExist},
		},
		{
			"failure: non-existent head node",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n1.id: e12(true, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 1)},
				},
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n1.id: e12(true, 1)},
				},
				idToHeads: map[ID]map[ID]Edge{
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

func TestGraph_GetWeight(t *testing.T) {
	type input struct {
		idTail StringID
		idHead StringID
	}
	type output struct {
		weight float64
		err    error
	}

	n1 := newTestNode("1")
	n2 := newTestNode("2")
	n3 := newTestNode("3")
	e12 := newTestEdgeGenerator("1", "2")
	e13 := newTestEdgeGenerator("1", "3")

	testCases := []struct {
		name  string
		graph *graph
		in    input
		out   output
	}{
		{
			"success",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n1.id: e12(true, 1.2)},
					n3.id: {n1.id: e13(true, 1.3)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 1.2), n3.id: e13(true, 1.3)},
				},
			},
			input{n1.id, n2.id},
			output{1.2, nil},
		},
		{
			"failure: non-existent tail node",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n1.id: e12(true, 1.2)},
					n3.id: {n1.id: e13(true, 1.3)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 1.2), n3.id: e13(true, 1.3)},
				},
			},
			input{n3.id, n2.id},
			output{0, ErrNodeNotExist},
		},
		{
			"failure: non-existent head node",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n1.id: e12(true, 1.2)},
					n3.id: {n1.id: e13(true, 1.3)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 1.2), n3.id: e13(true, 1.3)},
				},
			},
			input{n1.id, n3.id},
			output{0, ErrNodeNotExist},
		},
		{
			"failure: non-existent edge",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]Edge{
					n2.id: {n1.id: e12(true, 1.2)},
					n3.id: {n1.id: e13(true, 1.3)},
				},
				idToHeads: map[ID]map[ID]Edge{
					n1.id: {n2.id: e12(true, 1.2), n3.id: e13(true, 1.3)},
				},
			},
			input{n2.id, n3.id},
			output{0, ErrEdgeNotExist},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g, in, out := tc.graph, tc.in, tc.out

			weight, err := g.GetWeight(in.idTail, in.idHead)
			if err != out.err {
				t.Errorf("expected: %v, actual: %v", out.err, err)
			}
			if weight != out.weight {
				t.Errorf("expected: %f, actual: %f", out.weight, weight)
			}
		})
	}
}
