package iia

import (
	"math"
)

type Triangulator struct {
	Vertices  []*Vertex
	HalfEdges []*HalfEdge
}

func NewTriangulator(vertices []*Vertex) *Triangulator {
	triangulator := &Triangulator{
		Vertices:  vertices,
		HalfEdges: make([]*HalfEdge, 0, 2*(len(vertices)+3)-5),
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

	// 点群を覆う三角形を作る
	t.AddCoverTriangle()

	for _, v := range t.Vertices {
		// Triangle を持つ half-edge の内、その Triangle の外接円の内部に v を含むものを集める
		var he0s []*HalfEdge
		for _, th := range t.HalfEdges {
			if th.Triangle != nil {
				if th.Triangle.IsVertexInsideCircumcircle(*v) {
					he0s = append(he0s, th)
					break
				}
			}
		}

		// insideTriangle を v で3つに分割する
		halfEdgeStack := NewHalfEdgeStack()
		for _, he0 := range he0s {
			he1 := he0.Next
			he2 := he1.Next

			he0.Triangle = nil
			he1.Triangle = nil
			he2.Triangle = nil

			nhe0, nhe1, nhe2, nhe3, nhe4, nhe5 := &HalfEdge{}, &HalfEdge{}, &HalfEdge{}, &HalfEdge{}, &HalfEdge{}, &HalfEdge{}

			t.HalfEdges = append(t.HalfEdges, nhe0)
			t.HalfEdges = append(t.HalfEdges, nhe1)
			t.HalfEdges = append(t.HalfEdges, nhe2)
			t.HalfEdges = append(t.HalfEdges, nhe3)
			t.HalfEdges = append(t.HalfEdges, nhe4)
			t.HalfEdges = append(t.HalfEdges, nhe5)

			nhe0.End = v
			nhe1.End = he0.End
			nhe2.End = v
			nhe3.End = he1.End
			nhe4.End = v
			nhe5.End = he2.End

			nhe0.Pair = nhe1
			nhe1.Pair = nhe0
			nhe2.Pair = nhe3
			nhe3.Pair = nhe2
			nhe4.Pair = nhe5
			nhe5.Pair = nhe4

			he0.Next = nhe0
			nhe0.Next = nhe5
			nhe5.Next = he0
			he1.Next = nhe2
			nhe2.Next = nhe1
			nhe1.Next = he1
			he2.Next = nhe4
			nhe4.Next = nhe3
			nhe3.Next = he2

			newTriangle0 := NewTriangle(he0.End, v, he2.End)
			newTriangle1 := NewTriangle(he1.End, v, he0.End)
			newTriangle2 := NewTriangle(he2.End, v, he1.End)
			he0.Triangle = newTriangle0
			nhe0.Triangle = newTriangle0
			nhe5.Triangle = newTriangle0
			he1.Triangle = newTriangle1
			nhe2.Triangle = newTriangle1
			nhe1.Triangle = newTriangle1
			he2.Triangle = newTriangle2
			nhe4.Triangle = newTriangle2
			nhe3.Triangle = newTriangle2

			halfEdgeStack.Push(he0)
			halfEdgeStack.Push(he1)
			halfEdgeStack.Push(he2)
		}

		// flip-flop-loop
		for !halfEdgeStack.IsEmpty() {
			he := halfEdgeStack.Pop()
			if he.Pair == nil {
				continue
			}

			if he.Triangle.IsVertexInsideCircumcircle(*he.Pair.Next.End) {
				he0 := he
				he1 := he.Next
				he2 := he.Next.Next
				he3 := he.Pair
				he4 := he.Pair.Next
				he5 := he.Pair.Next.Next

				he0.Triangle = nil
				he1.Triangle = nil
				he2.Triangle = nil
				he3.Triangle = nil
				he4.Triangle = nil
				he5.Triangle = nil

				he0.End = he4.End
				he3.End = he1.End

				he0.Next, he1.Next, he2.Next, he3.Next, he4.Next, he5.Next = he5, he0, he4, he2, he3, he1

				newTriangle0 := NewTriangle(he0.End, he5.End, he1.End)
				newTriangle3 := NewTriangle(he3.End, he2.End, he4.End)
				he0.Triangle = newTriangle0
				he5.Triangle = newTriangle0
				he1.Triangle = newTriangle0
				he3.Triangle = newTriangle3
				he2.Triangle = newTriangle3
				he4.Triangle = newTriangle3

				halfEdgeStack.Push(he2)
				halfEdgeStack.Push(he4)
				halfEdgeStack.Push(he5)
				halfEdgeStack.Push(he1)
			}
		}
	}

	t.RemoveCoverTriangle()

	return nil
}

func (t *Triangulator) Validate() error {
	return nil
}

func (t *Triangulator) AddCoverTriangle() {
	vs := ComputeCoverTriangleVertices(t.Vertices)
	v0, v1, v2 := vs[0], vs[1], vs[2]
	v0.Index, v1.Index, v2.Index = len(t.Vertices), len(t.Vertices)+1, len(t.Vertices)+2

	if false {
		for _, v := range t.Vertices {
			abx := v1.X - v0.X
			aby := v1.Y - v0.Y
			bcx := v2.X - v1.X
			bcy := v2.Y - v1.Y
			cax := v0.X - v2.X
			cay := v0.Y - v2.Y

			bpx := v.X - v1.X
			bpy := v.Y - v1.Y
			cpx := v.X - v2.X
			cpy := v.Y - v2.Y
			apx := v.X - v0.X
			apy := v.Y - v0.Y

			c0 := abx*bpy - aby*bpx
			c1 := bcx*cpy - bcy*cpx
			c2 := cax*apy - cay*apx

			if !((c0 > 0 && c1 > 0 && c2 > 0) || (c0 < 0 && c1 < 0 && c2 < 0)) {
				panic("OMG")
			}
		}
	}

	t.HalfEdges = append(t.HalfEdges, &HalfEdge{End: v0})
	t.HalfEdges = append(t.HalfEdges, &HalfEdge{End: v1})
	t.HalfEdges = append(t.HalfEdges, &HalfEdge{End: v2})

	t.HalfEdges[0].Next = t.HalfEdges[1]
	t.HalfEdges[1].Next = t.HalfEdges[2]
	t.HalfEdges[2].Next = t.HalfEdges[0]

	newTriangle := NewTriangle(v0, v1, v2)
	t.HalfEdges[0].Triangle = newTriangle
	t.HalfEdges[1].Triangle = newTriangle
	t.HalfEdges[2].Triangle = newTriangle

	return
}

func (t *Triangulator) RemoveCoverTriangle() {
	outerVertices := make([]*Vertex, 0, 3)
	for i := 0; i < 3; i++ {
		outerVertices = append(outerVertices, t.HalfEdges[i].End)
	}

	t.HalfEdges = t.HalfEdges[3:len(t.HalfEdges)]
	th := t.HalfEdges
L:
	for i := 0; i < len(th); {
		tri := th[i].Triangle
		for _, ov := range outerVertices {
			if tri.Vertex0 == ov || tri.Vertex1 == ov || tri.Vertex2 == ov {
				if th[i].Pair != nil {
					th[i].Pair.Pair = nil
				}
				if th[i].Next.Pair != nil {
					th[i].Next.Pair.Pair = nil
				}
				if th[i].Next.Next.Pair != nil {
					th[i].Next.Next.Pair.Pair = nil
				}

				copy(th[i:], th[i+1:])
				th[len(th)-1] = nil
				th = th[:len(th)-1]
				t.HalfEdges = th
				continue L
			}
		}
		i++
	}
}

func ComputeCoverTriangleVertices(vertices []*Vertex) []*Vertex {
	xMin, yMin, xMax, yMax := ComputeBounds(vertices)

	centerX, centerY := (xMax+xMin)/2, (yMax+yMin)/2 // <=> xMin+(xMax-xMin)/2, yMin+(yMax-yMin)/2
	r := math.Hypot(centerX, centerY)
	r2 := r * 2
	r3r := r * math.Sqrt(3)

	return []*Vertex{
		&Vertex{X: centerX, Y: centerY + r2},
		&Vertex{X: centerX - r3r, Y: centerY - r},
		&Vertex{X: centerX + r3r, Y: centerY - r},
	}
}

func ComputeBounds(vertices []*Vertex) (xMin, yMin, xMax, yMax float64) {
	xMin, yMin, xMax, yMax = math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64
	for _, v := range vertices {
		if v.X < xMin {
			xMin = v.X
		} else if v.X > xMax {
			xMax = v.X
		}
		if v.Y < yMin {
			yMin = v.Y
		} else if v.Y > yMax {
			yMax = v.Y
		}
	}
	return
}
