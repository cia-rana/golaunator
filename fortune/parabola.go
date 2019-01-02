package fortune

import (
	"math"
)

var (
	posiInf = math.Inf(1)
	negaInf = -posiInf
)

type Parabola struct {
	Point   *Vector   // Forcus point of this
	l, r, p *Parabola // Left, Right, Parent Parabola-Node

	Edge *Edge
}

func NewParabola() *Parabola {
	return &Parabola{}
}

func NewParabolaWithPoint(point *Vector) *Parabola {
	return &Parabola{Point: point}
}

func (p *Parabola) GetLeftChild() *Parabola {
	return p.l
}

func (p *Parabola) GetRightChild() *Parabola {
	return p.r
}

func (p *Parabola) GetParent() *Parabola {
	return p.p
}

func (p *Parabola) SetPrev(q *Parabola) {
	if p.l == nil {
		p.l = q
		q.p = p
		return
	}

	par := p.l
	for par.r != nil {
		par = par.r
	}
	par.r = q
	q.p = par
}

func (p *Parabola) SetNext(q *Parabola) {
	if p.r == nil {
		p.r = q
		q.p = p
		return
	}

	par := p.r
	for par.l != nil {
		par = par.l
	}
	par.l = q
	q.p = par
}

func (p *Parabola) GetPrev() *Parabola {
	if p.l != nil {
		p = p.l
		for p.r != nil {
			p = p.r
		}
		return p
	}

	pp := p.p
	for pp != nil && p == pp.l {
		p, pp = pp, pp.p
	}
	return pp
}

func (p *Parabola) GetNext() *Parabola {
	if p.r != nil {
		par := p.r
		if par.l != nil && par.l.l != nil && par.l.l.l != nil {
			//fmt.Println(par.Point.Index, par.l.Point.Index, par.l.l.Point.Index, par.l.l.l.Point.Index)
		}
		for par.l != nil {
			par = par.l
		}
		return par
	}

	par, pp := p, p.p
	for pp != nil && par == pp.r {
		par, pp = pp, pp.p
	}
	return pp
}

func (p *Parabola) GetMin() *Parabola {
	for par := p; par != nil; par = par.l {
		p = par
	}
	return p
}

func (p *Parabola) GetMax() *Parabola {
	for par := p; par != nil; par = par.r {
		p = par
	}
	return p
}

func (p *Parabola) Delete() (bool, *Parabola) {
	defer func() {
		p = nil
	}()

	if p.l == nil && p.r == nil {
		return p.RecognizeParent(nil)
	}
	if p.l == nil {
		p.r.p = p.p
		return p.RecognizeParent(p.r)
	}
	if p.r == nil {
		p.l.p = p.p
		return p.RecognizeParent(p.l)
	}
	if p.r.l == nil {
		p.r.p = p.p
		p.r.l = p.l
		p.l.p = p.r
		return p.RecognizeParent(p.r)
	}

	pRMin := p.r.GetMin()

	if pRMin.r != nil {
		pRMin.r.p = pRMin.p
	}

	pRMin.p.l = pRMin.r
	p.l.p = pRMin
	p.r.p = pRMin
	pRMin.l = p.l
	pRMin.r = p.r
	pRMin.p = p.p

	return p.RecognizeParent(pRMin)
}

func (p *Parabola) RecognizeParent(q *Parabola) (bool, *Parabola) {
	if p.p == nil {
		return true, q
	}

	if p.p.l == p {
		p.p.l = q
	} else {
		p.p.r = q
	}

	return false, nil
}

func (p *Parabola) IsLeaf() bool {
	return p.l == nil && p.r == nil
}

func (p *Parabola) GetParabolaByX(point *Vector) *Parabola {
	par := p
	for !par.IsLeaf() {
		if point.X < par.GetLeftBreakPoint(point.Y) {
			if par.l == nil {
				break
			}
			par = par.l
		} else if par.GetRightBreakPoint(point.Y) < point.X {
			if par.r == nil {
				break
			}
			par = par.r
		} else {
			break
		}
	}

	return par
}

func (p *Parabola) GetLeftBreakPoint(ly float64) float64 {
	point := p.Point
	rFocX, rFocY := point.X, point.Y
	prby2 := rFocY - ly
	if prby2 == 0 {
		return rFocX
	}

	parL := p.GetPrev()
	if parL == nil {
		return negaInf
	}

	point = parL.Point
	lFocX, lFocY := point.X, point.Y
	plby2 := lFocY - ly
	if plby2 == 0 {
		return lFocX
	}

	if plby2 != prby2 {
		hl := lFocX - rFocX
		aby2 := 1/prby2 - 1/plby2
		b := hl / plby2
		return (-b+math.Sqrt(b*b-2*aby2*(hl*hl/(-2*plby2)-(lFocY-plby2/2)+(rFocY-prby2/2))))/aby2 + rFocX
	}
	return (lFocX + rFocX) / 2
}

func (p *Parabola) GetRightBreakPoint(ly float64) float64 {
	if parR := p.GetNext(); parR != nil {
		return parR.GetLeftBreakPoint(ly)
	}

	if p.Point.Y == ly {
		return p.Point.X
	}
	return posiInf
}

func (p *Parabola) Len() int {
	if p == nil {
		return 0
	}
	return 1 + p.l.Len() + p.r.Len()
}

func (p *Parabola) IsExist(q *Parabola) bool {
	if p == nil || q == nil {
		return false
	}
	return p == q || p.l.IsExist(q) || p.r.IsExist(q)
}
