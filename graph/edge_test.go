package graph

import "testing"

func TestNewEdge(t *testing.T) {
	nTail := NewNode("1")
	nHead := NewNode("2")
	weight := 1.0

	e := NewEdge(nTail, nHead, weight)
	if e.Tail().ID() != nTail.ID() {
		t.Errorf("expected: %q, actual: %q", nTail.ID(), e.Tail().ID())
	}
	if e.Head().ID() != nHead.ID() {
		t.Errorf("expected: %q, actual: %q", nHead.ID(), e.Head().ID())
	}
	if e.Weight() != weight {
		t.Errorf("expected: %f, actual: %f", weight, e.Weight())
	}
}
