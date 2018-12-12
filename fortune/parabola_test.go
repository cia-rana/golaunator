package fortune

import (
	"sort"
	"testing"
)

var points = []*Vector{
	{X: 2, Y: 0},
	{X: 8, Y: 0},
	{X: 1, Y: 0},
	{X: 3, Y: 0},
	{X: 4, Y: 0},
	{X: 7, Y: 0},
	{X: 5, Y: 0},
	{X: 9, Y: 0},
	{X: 0, Y: 0},
}

func TestParabolaInOrder(t *testing.T) {
	var par *Parabola
	var expected []float64

	parRoot := NewParabolaWithPoint(&Vector{X: 6, Y: 0})
	for _, point := range points {
		parRoot.set(NewParabolaWithPoint(point))
	}

	// Test for Parabola.Next
	par = parRoot.GetMin()
	expected = []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	for i, e := range expected {
		if par == nil {
			t.Fatalf("%d\n", i)
		} else if e != par.Point.X {
			t.Fatalf("expected: %f, actual: %f\n", e, par.Point.X)
		}
		par = par.GetNext()
	}

	// Test for Parabola.Prev
	par = parRoot.GetMax()
	sort.Sort(sort.Reverse(sort.Float64Slice(expected)))
	for i, e := range expected {
		if par == nil {
			t.Fatalf("%d\n", i)
		} else if e != par.Point.X {
			t.Fatalf("expected: %f, actual: %f\n", e, par.Point.X)
		}
		par = par.GetPrev()
	}
}

func TestParabolaDelete(t *testing.T) {
	var par *Parabola
	var expected []float64

	parRoot := NewParabolaWithPoint(&Vector{X: 6, Y: 0})
	for _, point := range points {
		parRoot.set(NewParabolaWithPoint(point))
	}

	if parRootIsChanged, parRootCandidate := parRoot.GetLeftChild().Delete(); parRootIsChanged {
		parRoot = parRootCandidate
	}

	expected = []float64{0, 1, 3, 4, 5, 6, 7, 8, 9}
	if l := parRoot.Len(); l != len(expected) {
		t.Fatalf("parRoot.Len(): %d, len(expected): %d\n", l, len(expected))
	}
	par = parRoot.GetMin()
	for i, e := range expected {
		if par == nil {
			t.Fatalf("%d\n", i)
		} else if e != par.Point.X {
			t.Fatalf("no: %d, expected: %f, actual: %f\n", i, e, par.Point.X)
		}
		par = par.GetNext()
	}
}

// Use only for root of parabola tree
func (p *Parabola) set(par *Parabola) {
	if par == nil || p.Point == nil || par.Point == nil {
		return
	}

	for {
		if par.Point.X < p.Point.X {
			if p.l == nil {
				p.l = par
				par.p = p
				break
			}
			p = p.l
		} else {
			if p.r == nil {
				p.r = par
				par.p = p
				break
			}
			p = p.r
		}
	}
}
