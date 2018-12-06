package fortune

type Vertex struct {
	X, Y  float64
	Index int
}

func Vector2Vertex(vector *Vector) *Vertex {
	if vector == nil {
		return nil
	}
	return &Vertex{X: vector.X, Y: vector.Y, Index: vector.Index}
}
