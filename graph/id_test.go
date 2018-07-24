package graph

import "testing"

func TestID(t *testing.T) {
	s := "1"

	id := ID(s)
	if id.String() != s {
		t.Errorf("expected: %q, actual: %q", s, id.String())
	}
}
