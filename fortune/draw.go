package fortune

import (
	"math"

	"github.com/fogleman/gg"
)

func DrawParabola(dc *gg.Context, focX, focY, directrix, startX, endX float64) {
	if math.Abs(focY-directrix) < math.SmallestNonzeroFloat64 {
		dc.DrawLine(focX, 0, focX, directrix)
		return
	}

	tq := 1.0 / (focY - directrix)
	a := 0.5 * tq
	b := -focX * tq
	c := (focX*focX + focY*focY - directrix*directrix) * tq * 0.5

	x0 := startX
	y0 := a*x0*x0 + b*x0 + c
	x2 := endX
	y2 := a*x2*x2 + b*x2 + c
	x1 := (x0 + x2) * 0.5
	y1 := a*x0*x2 + 0.5*b*(x0+x2) + c
	points := gg.QuadraticBezier(x0, y0, x1, y1, x2, y2)

	for i := 1; i < len(points); i++ {
		dc.LineTo(points[i].X, points[i].Y)
	}
}
