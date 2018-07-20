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
	testIDToNodesEquality(t, expected.idToNodes, actual.idToNodes)
	testIDToEndsEquality(t, expected.idToTails, actual.idToTails)
	testIDToEndsEquality(t, expected.idToHeads, actual.idToHeads)
	if expected.isDirected != actual.isDirected {
		t.Errorf("expected: %t, actual: %t", expected.isDirected, actual.isDirected)
	}
}

func testIDToNodesEquality(t *testing.T, expected, actual map[ID]Node) {
	if len(actual) != len(expected) {
		t.Errorf("expected: %d, actual: %d", len(expected), len(actual))
	}
	for id, n := range expected {
		if _, ok := actual[id]; !ok {
			t.Errorf("expected: %t, actual: %t", true, ok)
		} else {
			if actual[id].ID() != n.ID() {
				t.Errorf("expected: %q, actual: %q", n.ID(), actual[id].ID())
			}
		}
	}
}

func testIDToEndsEquality(t *testing.T, expected, actual map[ID]map[ID]float64) {
	if len(actual) != len(expected) {
		t.Errorf("expected: %d, actual: %d", len(expected), len(actual))
	}
	for id1, ends := range expected {
		if _, ok := actual[id1]; !ok {
			t.Errorf("expected: %t, actual: %t", true, ok)
		} else {
			if len(actual[id1]) != len(ends) {
				t.Errorf("expected: %d, actual: %d", len(ends), len(actual[id1]))
			}
			for id2, weight := range ends {
				if _, ok := actual[id1][id2]; !ok {
					t.Errorf("expected: %t, actual: %t", true, ok)
				} else {
					if actual[id1][id2] != weight {
						t.Errorf("expected: %f, actual: %f", weight, actual[id1][id2])
					}
				}
			}
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
				idToTails: map[ID]map[ID]float64{},
			},
			input{n1.id},
			output{map[ID]Node{}, nil},
		},
		{
			"success",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0, n3.id: 1.0},
				},
			},
			input{n1.id},
			output{map[ID]Node{n2.id: n2, n3.id: n3}, nil},
		},
		{
			"failure: non-existent node",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
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
				idToHeads: map[ID]map[ID]float64{},
			},
			input{n1.id},
			output{map[ID]Node{}, nil},
		},
		{
			"success",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0, n3.id: 1.0},
				},
			},
			input{n1.id},
			output{map[ID]Node{n2.id: n2, n3.id: n3}, nil},
		},
		{
			"failure: non-existent node",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
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
				idToTails: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0, n3.id: 1.0},
					n2.id: {n1.id: 1.0, n3.id: 1.0},
					n3.id: {n1.id: 1.0, n2.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0, n3.id: 1.0},
					n2.id: {n1.id: 1.0, n3.id: 1.0},
					n3.id: {n1.id: 1.0, n2.id: 1.0},
				},
			},
			&graph{
				idToNodes: map[ID]Node{n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n2.id: {n3.id: 1.0},
					n3.id: {n2.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n2.id: {n3.id: 1.0},
					n3.id: {n2.id: 1.0},
				},
			},
			input{n1.id},
			output{nil},
		},
		{
			"success: non-existent node",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]float64{
					n2.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
				},
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]float64{
					n2.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
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
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails:  map[ID]map[ID]float64{},
				idToHeads:  map[ID]map[ID]float64{},
				isDirected: true,
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n2.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
				},
				isDirected: true,
			},
			input{n1.id, n2.id, 1.0},
			output{nil},
		},
		{
			"success: directed, second edge",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n2.id: {n1.id: 1.0},
					n3.id: {n2.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
					n2.id: {n3.id: 1.0},
				},
				isDirected: true,
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n2.id: {n1.id: 1.0},
					n3.id: {n1.id: 1.0, n2.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0, n3.id: 1.0},
					n2.id: {n3.id: 1.0},
				},
				isDirected: true,
			},
			input{n1.id, n3.id, 1.0},
			output{nil},
		},
		{
			"success: directed, existent edge",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n2.id: {n1.id: 1.0},
					n3.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0, n3.id: 1.0},
				},
				isDirected: true,
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n2.id: {n1.id: 2.0},
					n3.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 2.0, n3.id: 1.0},
				},
				isDirected: true,
			},
			input{n1.id, n2.id, 1.0},
			output{nil},
		},
		{
			"success: undirected, first edge",
			&graph{
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails:  map[ID]map[ID]float64{},
				idToHeads:  map[ID]map[ID]float64{},
				isDirected: false,
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
					n2.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
					n2.id: {n1.id: 1.0},
				},
				isDirected: false,
			},
			input{n1.id, n2.id, 1.0},
			output{nil},
		},
		{
			"success: undirected, second edge",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
					n2.id: {n1.id: 1.0, n3.id: 1.0},
					n3.id: {n2.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
					n2.id: {n1.id: 1.0, n3.id: 1.0},
					n3.id: {n2.id: 1.0},
				},
				isDirected: false,
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0, n3.id: 1.0},
					n2.id: {n1.id: 1.0, n3.id: 1.0},
					n3.id: {n1.id: 1.0, n2.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0, n3.id: 1.0},
					n2.id: {n1.id: 1.0, n3.id: 1.0},
					n3.id: {n1.id: 1.0, n2.id: 1.0},
				},
				isDirected: false,
			},
			input{n1.id, n3.id, 1.0},
			output{nil},
		},
		{
			"success: undirected, existent edge",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0, n3.id: 1.0},
					n2.id: {n1.id: 1.0},
					n3.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0, n3.id: 1.0},
					n2.id: {n1.id: 1.0},
					n3.id: {n1.id: 1.0},
				},
				isDirected: false,
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n1.id: {n2.id: 2.0, n3.id: 1.0},
					n2.id: {n1.id: 2.0},
					n3.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 2.0, n3.id: 1.0},
					n2.id: {n1.id: 2.0},
					n3.id: {n1.id: 1.0},
				},
				isDirected: false,
			},
			input{n1.id, n2.id, 1.0},
			output{nil},
		},
		{
			"failure: non-existent tail node",
			&graph{
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails:  map[ID]map[ID]float64{},
				idToHeads:  map[ID]map[ID]float64{},
				isDirected: false,
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]float64{},
				idToHeads: map[ID]map[ID]float64{},
			},
			input{n3.id, n2.id, 1.0},
			output{ErrNodeNotExist},
		},
		{
			"failure: non-existent head node",
			&graph{
				idToNodes:  map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails:  map[ID]map[ID]float64{},
				idToHeads:  map[ID]map[ID]float64{},
				isDirected: false,
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]float64{},
				idToHeads: map[ID]map[ID]float64{},
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
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n2.id: {n1.id: 1.0},
					n3.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0, n3.id: 1.0},
				},
				isDirected: true,
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n2.id: {},
					n3.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n3.id: 1.0},
				},
				isDirected: true,
			},
			input{n1.id, n2.id},
			output{nil},
		},
		{
			"success: directed, non-existent edge",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n2.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
				},
				isDirected: true,
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n2.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
				},
				isDirected: true,
			},
			input{n1.id, n3.id},
			output{nil},
		},
		{
			"success: undirected",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0, n3.id: 1.0},
					n2.id: {n1.id: 1.0},
					n3.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0, n3.id: 1.0},
					n2.id: {n1.id: 1.0},
					n3.id: {n1.id: 1.0},
				},
				isDirected: false,
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n1.id: {n3.id: 1.0},
					n2.id: {},
					n3.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n3.id: 1.0},
					n2.id: {},
					n3.id: {n1.id: 1.0},
				},
				isDirected: false,
			},
			input{n1.id, n2.id},
			output{nil},
		},
		{
			"success: undirected, non-existent edge",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
					n2.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
					n2.id: {n1.id: 1.0},
				},
				isDirected: false,
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
					n2.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
					n2.id: {n1.id: 1.0},
				},
				isDirected: false,
			},
			input{n1.id, n3.id},
			output{nil},
		},
		{
			"failure: non-existent tail node",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
					n2.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
					n2.id: {n1.id: 1.0},
				},
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
					n2.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
					n2.id: {n1.id: 1.0},
				},
			},
			input{n3.id, n2.id},
			output{ErrNodeNotExist},
		},
		{
			"failure: non-existent head node",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
					n2.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
					n2.id: {n1.id: 1.0},
				},
			},
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
					n2.id: {n1.id: 1.0},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.0},
					n2.id: {n1.id: 1.0},
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
				idToTails: map[ID]map[ID]float64{
					n2.id: {n1.id: 1.2},
					n3.id: {n1.id: 1.3},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.2, n3.id: 1.3},
				},
			},
			input{n1.id, n2.id},
			output{1.2, nil},
		},
		{
			"failure: non-existent tail node",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]float64{
					n2.id: {n1.id: 1.2},
					n3.id: {n1.id: 1.3},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.2, n3.id: 1.3},
				},
			},
			input{n3.id, n2.id},
			output{0, ErrNodeNotExist},
		},
		{
			"failure: non-existent head node",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2},
				idToTails: map[ID]map[ID]float64{
					n2.id: {n1.id: 1.2},
					n3.id: {n1.id: 1.3},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.2, n3.id: 1.3},
				},
			},
			input{n1.id, n3.id},
			output{0, ErrNodeNotExist},
		},
		{
			"failure: non-existent edge",
			&graph{
				idToNodes: map[ID]Node{n1.id: n1, n2.id: n2, n3.id: n3},
				idToTails: map[ID]map[ID]float64{
					n2.id: {n1.id: 1.2},
					n3.id: {n1.id: 1.3},
				},
				idToHeads: map[ID]map[ID]float64{
					n1.id: {n2.id: 1.2, n3.id: 1.3},
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
