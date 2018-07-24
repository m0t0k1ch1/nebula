package graph

import "testing"

func newTestEdgeGenerator(idTail, idHead string) func(bool, float64) *Edge {
	return func(isDirected bool, weight float64) *Edge {
		return newTestEdge(isDirected, idTail, idHead, weight)
	}
}

func newTestEdge(isDirected bool, idTail, idHead string, weight float64) *Edge {
	return &Edge{
		isDirected: isDirected,
		tail:       newTestNode(idTail),
		head:       newTestNode(idHead),
		weight:     weight,
	}
}

func testEdgeEquality(t *testing.T, expected, actual *Edge) {
	if actual.IsDirected() != expected.IsDirected() {
		t.Errorf("expected: %t, actual: %t", expected.IsDirected(), actual.IsDirected())
	}
	if actual.Tail().ID() != expected.Tail().ID() {
		t.Errorf("expected: %q, actual: %q", expected.Tail().ID(), actual.Tail().ID())
	}
	if actual.Head().ID() != expected.Head().ID() {
		t.Errorf("expected: %q, actual: %q", expected.Head().ID(), actual.Head().ID())
	}
	if actual.Weight() != expected.Weight() {
		t.Errorf("expected: %f, actual: %f", expected.Weight(), actual.Weight())
	}
}

func TestNewDirectedEdge(t *testing.T) {
	nTail := newTestNode("tail")
	nHead := newTestNode("head")
	weight := 1.0

	e := NewDirectedEdge(nTail, nHead, weight)
	if !e.IsDirected() {
		t.Errorf("expected: %t, actual: %t", true, e.IsDirected())
	}
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

func TestNewUndirecredEdge(t *testing.T) {
	nTail := newTestNode("tail")
	nHead := newTestNode("head")
	weight := 1.0

	e := NewUndirectedEdge(nTail, nHead, weight)
	if e.IsDirected() {
		t.Errorf("expected: %t, actual: %t", false, e.IsDirected())
	}
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

func TestSetWeight(t *testing.T) {
	e := &Edge{
		weight: 1.0,
	}

	e.SetWeight(2.0)
	if e.Weight() != 2.0 {
		t.Errorf("expected: %f, actual: %f", 2.0, e.Weight())
	}
}
