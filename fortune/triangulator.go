package fortune

import (
	"fmt"
	"math"

	"github.com/fogleman/gg"
)

type Triangulator struct {
	Vertices []*Vertex
	Edges    []*Edge

	parabolaTree *ParabolaTree

	eventQueue *EventQueue

	dc           *gg.Context
	enableDraw   bool
	drawingCount int
}

func NewTriangulator(vertices []*Vertex) *Triangulator {
	triangulator := &Triangulator{
		Vertices:     vertices,
		Edges:        make([]*Edge, 0, 2*len(vertices)-5),
		parabolaTree: NewParabolaTree(),
		eventQueue:   NewEventQueue(),
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

		switch event.Type {
		case Site:
			t.insertParabola(event)
		case Circle:
			t.removeParabola(event)
		default:
			panic(fmt.Sprintf("undifined type: %v", event.Type))
		}
	}

	return nil
}

func (t *Triangulator) insertParabola(event *Event) {
	point := event.Point

	if t.parabolaTree.GetRoot() == nil {
		t.parabolaTree.SetRoot(NewParabolaWithPoint(point))
		return
	}

	if parRoot := t.parabolaTree.GetRoot(); parRoot.IsLeaf() && event.Point.Y-parRoot.Point.Y < math.SmallestNonzeroFloat64 {
		rp := parRoot.Point
		parRoot.SetNext(NewParabolaWithPoint(point))
		t.Edges = append(t.Edges, NewEdge(Vector2Vertex(point), Vector2Vertex(rp)))
		return
	}

	par := t.parabolaTree.GetParabolaByX(point)

	t.Edges = append(t.Edges, NewEdge(Vector2Vertex(par.Point), Vector2Vertex(point)))

	parA := NewParabolaWithPoint(par.Point)
	parB := NewParabolaWithPoint(point)
	parC := NewParabolaWithPoint(par.Point)

	par.SetNext(parB)
	parB.SetPrev(parA)
	parB.SetNext(parC)

	t.parabolaTree.Delete(par)

	if t.enableDraw {
		t.draw(point)
	}

	t.CheckCircle(parA)
	t.CheckCircle(parC)
}

func (t *Triangulator) removeParabola(event *Event) {
	parB := event.Parabola

	if !t.parabolaTree.IsExist(parB) {
		return
	}

	parA := parB.GetPrev()
	parC := parB.GetNext()

	t.Edges = append(t.Edges, NewEdge(Vector2Vertex(parA.Point), Vector2Vertex(parC.Point)))

	t.parabolaTree.Delete(parB)

	if t.enableDraw {
		t.draw(event.Point)
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

func (t *Triangulator) EnableDraw() {
	if t.dc == nil {
		t.dc = gg.NewContext(1000, 1000)
		t.dc.SetHexColor("fff")
		t.dc.Clear()
		t.drawingCount = 0
	}
	t.enableDraw = true
}

func (t *Triangulator) DisenableDraw() {
	t.enableDraw = false
}

func (t *Triangulator) draw(point *Vector) {

	t.dc.SetHexColor("000")
	t.dc.DrawString(fmt.Sprintf("%d", t.drawingCount), 20, 20)
	{
		var i int
		for p := t.parabolaTree.GetMin(); p != nil; p = p.GetNext() {
			t.dc.SetHexColor("f00")
			bpl := p.GetLeftBreakPoint(point.Y)
			bpr := p.GetRightBreakPoint(point.Y)
			DrawParabola(t.dc, p.Point.X, p.Point.Y, point.Y, math.Max(bpl, 0), Clamp(0, float64(t.dc.Width()), bpr))
			t.dc.DrawString(fmt.Sprintf("%dL", i), bpl-10, p.Point.Y+10*float64(i))
			t.dc.DrawString(fmt.Sprintf("%dR", i), bpr+10, p.Point.Y+10*float64(i))
			t.dc.Stroke()
			i++
		}
	}

	t.dc.SetHexColor("00f")
	for _, e := range t.Edges {
		t.dc.DrawLine(e.Start.X, e.Start.Y, e.End.X, e.End.Y)
	}

	t.dc.SetHexColor("00f")
	for i, v := range t.Vertices {
		t.dc.DrawPoint(v.X, v.Y, 3)
		t.dc.DrawString(fmt.Sprintf("%d", i), v.X, v.Y+30)
		t.dc.Stroke()
	}

	t.dc.SetHexColor("0f0")
	t.dc.DrawLine(0, point.Y, 1000, point.Y)
	t.dc.Stroke()

	t.savePNG()
}

func (t *Triangulator) savePNG() {
	err := t.dc.SavePNG(fmt.Sprintf("img/fortune%03d.png", t.drawingCount))
	if err != nil {
		fmt.Println(err)
	}
	t.drawingCount++
	t.dc.SetHexColor("fff")
	t.dc.Clear()
}

func print(i ...interface{}) {
	fmt.Println(i...)
}
