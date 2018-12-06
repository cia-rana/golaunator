package fortune

import (
	"fmt"
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
		fmt.Println(event, t.ly)

		switch event.Type {
		case Site:
			t.insertParabola(event)
		case Circle:
			t.removeParabola(event)
		default:
			// Do nothing!
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
		t.parRoot.SetRightChild(NewParabolaWithPoint(point))
		t.Edges = append(t.Edges, NewEdge(Vector2Vertex(point), Vector2Vertex(rp)))
		return
	}

	par := t.parRoot.GetParabolaByX(point)
	fmt.Println(par.Point, point)

	edge := NewEdge(Vector2Vertex(par.Point), Vector2Vertex(point))

	t.Edges = append(t.Edges, edge)

	parA := par
	parB := NewParabolaWithPoint(point)
	parC := NewParabolaWithPoint(par.Point)

	parA.SefRightChild(parB)
	parB.SetRightChild(parC)

	t.CheckCircle(parA)
	t.CheckCircle(parC)
}

func (t *Triangulator) removeParabola(event *Event) {
	parB := event.Parabola

	xl := parB.GetLeftParent()
	xr := parB.GetRightParent()

	parA := xl.GetLeftLeaf()
	parC := xr.GetRightLeaf()

	if parA == parC {
		panic("")
		return
	}

	gparent := parB.GetParent()
	ggparent := gparent.GetParent()
	if gparent.GetLeftChild() == parB {
		if ggparent.GetLeftChild() == gparent {
			ggparent.SetLeftChild(gparent.GetRightChild())
		}
		if ggparent.GetRightChild() == gparent {
			ggparent.SetRightChild(gparent.GetRightChild())
		}
	} else {
		if ggparent.GetLeftChild() == gparent {
			ggparent.SetLeftChild(gparent.GetLeftChild())
		}
		if ggparent.GetRightChild() == gparent {
			ggparent.SetRightChild(gparent.GetLeftChild())
		}
	}

	a, b, c := parA.Point, parB.Point, parC.Point
	x, y := calcCircleCenter(a, b, c)
	dx, dy := b.X-x, b.Y-y
	d := math.Sqrt(dx*dx + dy*dy)
	fmt.Println(a, b, c, y+d)

	t.CheckCircle(parA)
	t.CheckCircle(parC)

	return
}

func (t *Triangulator) CheckCircle(par *Parabola) {
	parA := par.GetPrev()
	parB := par
	parC := par.GetNext()
	if parA == nil || parC == nil {
		return
	}

	a := parA.Point
	b := par.Point
	c := parC.Point
	if a == c {
		return
	}

	ax := a.X - b.X
	ay := a.Y - b.Y
	cx := c.X - b.X
	xy := c.Y - b.Y

	d := 2 * (ax*cy - ay*cx)
	if d >= -2*math.SmallestNonzeroFloat64 { // <=> ay * cx - ax * cy <= e^-12
		return
	}

	/*
		x, y := calcCircleCenter(a, b, c)

		dx := b.X - x
		dy := b.Y - y
		d := math.Sqrt(dx*dx + dy*dy)

		event := &Event{Type: Circle, Point: &Vector{X: x, Y: y + d}}
		event.Arch = par
		t.eventQueue.Push(event)
	*/
}

func calcCircleCenter(a, b, c *Vector) (float64, float64) {
	pbABx := -(a.Y - b.Y) / (a.X - b.X)
	pbABy := (a.X + a.Y + b.X + b.Y) / 2
	pbBCx := -(b.Y - c.Y) / (b.X - c.X)

	x := -(a.X + a.Y - (c.X + c.Y)) / (pbABx - pbBCx)
	y := pbABx*x + pbABy

	return x, y
}
