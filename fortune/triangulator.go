package fortune

import (
	"math"
)

type Triangulator struct {
	Vertices []*Vertex
	Edges    []*Edge

	parRoot *Parabola // Root of parabola tree

	eventQueue *EventQueue

	ly float64 // Current y position of the sweep line on plain
}

func NewTriangulator(vertices []*Vertex) *Triangulator {
	triangulator := &Triangulator{
		Vertices:   vertices,
		Edges:      make([]*Edge, 0, 2*len(vertices)-5),
		eventQueue: NewEventQueue(),
	}

	for i, _ := range vertices {
		triangulator.Vertices[i].Index = i
	}

	return triangulator
}

func (t *Triangulator) Triangulate() error {
	if len(t.Vertices) < 3 {
		return nil
	}

	for _, v := range t.Vertices {
		t.eventQueue.Push(&Event{Type: Site, Point: Vertex2Vector(v)})
	}

	for !t.eventQueue.Empty() {
		event := t.eventQueue.Pull()

		t.ly = event.Point.Y

		switch event.Type {
		case Site:
			t.insertParabola(event)
		case Circle:
			t.removeParabola(event)
		default:
			panic("")
		}
	}

	return nil
}

func (t *Triangulator) insertParabola(event *Event) {
	point := event.Point

	if t.parRoot == nil {
		t.parRoot = NewParabolaWithPoint(point)
		return
	}

	if t.parRoot.IsLeaf() && event.Point.Y-t.parRoot.Point.Y < math.SmallestNonzeroFloat64 {
		rp := t.parRoot.Point
		t.parRoot.SetNext(NewParabolaWithPoint(point))
		t.Edges = append(t.Edges, NewEdge(Vector2Vertex(point), Vector2Vertex(rp)))
		return
	}

	par := t.parRoot.GetParabolaByX(point)

	edge := NewEdge(Vector2Vertex(par.Point), Vector2Vertex(point))

	t.Edges = append(t.Edges, edge)

	parA := par
	parB := NewParabolaWithPoint(point)
	parC := NewParabolaWithPoint(par.Point)

	parA.SetNext(parB)
	parB.SetNext(parC)

	t.CheckCircle(parA)
	t.CheckCircle(parC)
}

func (t *Triangulator) removeParabola(event *Event) {
	parB := event.Parabola
	parA := parB.GetPrev()
	parC := parB.GetNext()

	edge := NewEdge(Vector2Vertex(parA.Point), Vector2Vertex(parC.Point))

	t.Edges = append(t.Edges, edge)

	if parRootIsChanged, parRootCandidate := parB.Delete(); parRootIsChanged {
		t.parRoot = parRootCandidate
	}

	t.CheckCircle(parA)
	t.CheckCircle(parC)
}

func (t *Triangulator) CheckCircle(par *Parabola) {
	parA := par.GetPrev()
	parB := par
	parC := par.GetNext()
	if parA == nil || parC == nil {
		return
	}

	a := parA.Point
	b := parB.Point
	c := parC.Point
	if a == c {
		return
	}

	ax := a.X - b.X
	ay := a.Y - b.Y
	cx := c.X - b.X
	cy := c.Y - b.Y

	d := 2 * (ax*cy - ay*cx)
	//   ABxBC <= e-12
	// <=>
	//   (b.X - a.X) * (c.Y - b.Y) - (b.Y - a.Y) * (c.X - b.X)
	// = (-ax)*cy - (-ay)*cx
	// = -ax*cy + ay*cx <= e^-12
	// <=>
	//   2*(ax*cy - ay*cx) >= -2e^-12
	// <=>
	if d >= -2*math.SmallestNonzeroFloat64 {
		return
	}

	ha := ax*ax + ay*ay
	hc := cx*cx + cy*cy
	x := (cy*ha - ay*hc) / d
	y := (ax*hc - cx*ha) / d

	t.eventQueue.Push(&Event{
		Type: Circle,
		Point: &Vector{
			X: x + b.X,
			Y: y + b.Y + math.Sqrt(x*x+y*y),
		},
		Parabola: par,
	})
}
