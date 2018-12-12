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

func (p Parabola) GetLeftChild() *Parabola {
	return p.l
}

func (p Parabola) GetRightChild() *Parabola {
	return p.r
}

func (p *Parabola) GetParent() *Parabola {
	return p.p
}

func (p *Parabola) SetPrev(par *Parabola) {
	if p.l == nil {
		p.l = par
		par.p = p
		return
	}

	p = p.l
	for p.r != nil {
		p = p.r
	}
	p.r = par
	par.p = p
}

func (p *Parabola) SetNext(par *Parabola) {
	if p.r == nil {
		p.r = par
		par.p = p
		return
	}

	p = p.r
	for p.l != nil {
		p = p.l
	}
	p.l = par
	par.p = p
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
		p = p.r
		for p.l != nil {
			p = p.l
		}
		return p
	}

	pp := p.p
	for pp != nil && p == pp.r {
		p, pp = pp, pp.p
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
	if p.l == nil && p.r == nil {
		return p.RecognizeParent(nil)
	} else if p.l == nil {
		p.r.p = p.p
		return p.RecognizeParent(p.r)
	} else if p.r == nil {
		p.l.p = p.p
		return p.RecognizeParent(p.l)
	}

	pRMin := p.r.GetMin()

	if pRMin.r != nil {
		pRMin.r.p = pRMin.p
		if pRMin.p.l == pRMin {
			pRMin.p.l = pRMin.r
		} else {
			pRMin.p.r = pRMin.r
		}
	}

	if pRMin != p.r {
		p.r.p = pRMin
		pRMin.r = p.r
	}
	p.l.p = pRMin
	pRMin.l = p.l
	pRMin.p = p.p
	return p.RecognizeParent(pRMin)
}

func (p *Parabola) RecognizeParent(q *Parabola) (bool, *Parabola) {
	if p.p == nil {
		return true, q
	} else if p.p.l == p {
		p.p.l = q
	} else {
		p.p.r = q
	}
	return false, nil
}

func (p Parabola) IsLeaf() bool {
	return p.l == nil && p.r == nil
}

func (p *Parabola) GetParabolaByX(point *Vector) *Parabola {
	par := p
	for !par.IsLeaf() {
		if par.GetLeftBreakPoint(point.Y)-point.X > 0 {
			par = par.GetLeftChild()
		} else if point.X-par.GetRightBreakPoint(point.X) > 0 {
			par = par.GetRightChild()
		} else {
			break
		}
	}

	return par
}

func (p Parabola) GetLeftBreakPoint(ly float64) float64 {
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
	plby2 := rFocY - ly
	if plby2 == 0 {
		return lFocX
	}

	if plby2 != prby2 {
		hl := lFocX - rFocX
		aby2 := 1/prby2 - 1/plby2
		b := hl / aby2
		return (-b+math.Sqrt(b*b-2*aby2*(hl*hl/(-2*plby2)-(lFocY-plby2/2)+(rFocY-prby2/2))))/aby2 + rFocX
	}
	return (lFocY + rFocX) / 2
}

func (p Parabola) GetRightBreakPoint(ly float64) float64 {
	parR := p.GetNext()
	if parR != nil {
		return parR.GetLeftBreakPoint(ly)
	}

	point := p.Point
	if point.Y == ly {
		return point.X
	}
	return posiInf
}

func (p *Parabola) Len() int {
	if p == nil {
		return 0
	}
	return 1 + p.l.Len() + p.r.Len()
}
