package graph

import "testing"

func TestNewNode(t *testing.T) {
	type input struct {
		id string
	}
	type output struct {
		node Node
	}
	testCases := []struct {
		name string
		in   input
		out  output
	}{
		{
			"success",
			input{"1"},
			output{&node{StringID("1")}},
		},
	}

	for i, tc := range testCases {
		t.Logf("[%d] %s", i, tc.name)
		in, out := tc.in, tc.out

		n := NewNode(in.id)
		if n.ID() != out.node.ID() {
			t.Errorf("expected: %q, actual: %q", out.node.ID(), n.ID())
		}
	}
}
