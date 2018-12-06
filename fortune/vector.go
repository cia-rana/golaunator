package fortune

type Vector struct {
	X, Y  float64
	Index int
}

func NewVector(x, y float64) *Vector {
	return &Vector{X: x, Y: y}
}

func Vertex2Vector(vertex *Vertex) *Vector {
	if vertex == nil {
		return nil
	}
	return &Vector{X: vertex.X, Y: vertex.Y, Index: vertex.Index}
}
