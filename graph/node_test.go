package graph

import "testing"

func newTestNode(id string) *node {
	return &node{
		id: StringID(id),
	}
}

func testNodeEquality(t *testing.T, expected, actual Node) {
	if actual.ID() != expected.ID() {
		t.Errorf("expected: %q, actual: %q", expected.ID(), actual.ID())
	}
}

func TestNewNode(t *testing.T) {
	id := StringID("1")

	n := NewNode(id.String())
	if n.ID() != id {
		t.Errorf("expected: %q, actual: %q", id, n.ID())
	}
}
