package graph

import "testing"

func TestNewEdge(t *testing.T) {
	nTail := newTestNode("tail")
	nHead := newTestNode("head")
	weight := 1.0

	e := NewEdge(nTail, nHead, weight)
	if e.Tail().ID() != nTail.id {
		t.Errorf("expected: %q, actual: %q", nTail.id, e.Tail().ID())
	}
	if e.Head().ID() != nHead.id {
		t.Errorf("expected: %q, actual: %q", nHead.id, e.Head().ID())
	}
	if e.Weight() != weight {
		t.Errorf("expected: %f, actual: %f", weight, e.Weight())
	}
}
