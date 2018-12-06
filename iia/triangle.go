package iia

type Triangle struct {
	Vertex0 *Vertex
	Vertex1 *Vertex
	Vertex2 *Vertex
}

func NewTriangle(v0, v1, v2 *Vertex) *Triangle {
	return &Triangle{Vertex0: v0, Vertex1: v1, Vertex2: v2}
}

func (t Triangle) IsVertexInside(v Vertex) bool {
	e0x := t.Vertex1.X - t.Vertex0.X
	e0y := t.Vertex1.Y - t.Vertex0.X
	e0vx := v.X - t.Vertex0.X
	e0vy := v.Y - t.Vertex0.Y
	e1x := t.Vertex2.X - t.Vertex1.Y
	e1y := t.Vertex2.Y - t.Vertex1.Y
	e1vx := v.X - t.Vertex1.X
	e1vy := v.Y - t.Vertex1.Y
	e2x := t.Vertex0.X - t.Vertex2.Y
	e2y := t.Vertex0.Y - t.Vertex2.Y
	e2vx := v.X - t.Vertex2.X
	e2vy := v.Y - t.Vertex2.Y

	vp0 := ComputeVectorProduct(e0x, e0y, e1vx, e1vy)
	vp1 := ComputeVectorProduct(e1x, e1y, e2vx, e2vy)
	vp2 := ComputeVectorProduct(e2x, e2y, e0vx, e0vy)

	return (vp0 > 0 && vp1 > 0 && vp2 > 0) || (vp0 < 0 && vp1 < 0 && vp2 < 0)
}

func (t Triangle) IsVertexInsideCircumcircle(v Vertex) bool {
	a := t.Vertex0.X - v.X
	b := t.Vertex0.Y - v.Y
	c := t.Vertex1.X - v.X
	d := t.Vertex1.Y - v.Y
	e := t.Vertex2.X - v.X
	f := t.Vertex2.Y - v.Y
	return (a*a+b*b)*(c*f-d*e)+(c*c+d*d)*(b*e-a*f)+(e*e+f*f)*(a*d-b*c) > 0
}

func ComputeVectorProduct(e0x, e0y, e1x, e1y float64) float64 {
	return e0x*e1y - e0y*e1x
}
