package graph

import "testing"

func TestStringID(t *testing.T) {
	s := "1"
	id := StringID(s)
	if id.String() != s {
		t.Errorf("expected: %q, actual: %q", s, id.String())
	}
}
