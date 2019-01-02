package fortune

type Edge struct {
	Start *Vertex
	End   *Vertex
}

func NewEdge(start, end *Vertex) *Edge {
	if start == nil || end == nil {
		return nil
	}

	e := &Edge{
		Start: start,
		End:   end,
	}
	return e
}
