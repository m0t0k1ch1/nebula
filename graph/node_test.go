package graph

import "testing"

func newTestNode(id string) *node {
	return &node{
		id: StringID(id),
	}
}

func TestNewNode(t *testing.T) {
	id := StringID("1")

	n := NewNode(id.String())
	if n.ID() != id {
		t.Errorf("expected: %q, actual: %q", id, n.ID())
	}
}
