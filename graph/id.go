package graph

type ID interface {
	String() string
}

type StringID string

func (s StringID) String() string {
	return string(s)
}
