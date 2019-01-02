package main

import (
	"math"

	"github.com/fogleman/gg"
)

func main() {
	w, h := 512, 512
	dc := gg.NewContext(w, h)
	dc.SetHexColor("fff")
	dc.Clear()

	var focX, focY float64 = float64(w) * 0.5, float64(h) * 0.5
	directrix := float64(h)*0.5 + 10

	dc.SetHexColor("000")
	dc.DrawPoint(focX, focY, 2)
	dc.Fill()

	dc.SetHexColor("f00")
	dc.SetLineWidth(1)
	DrawParabola(dc, focX, focY, directrix, 0, float64(w))

	for i := 0; i < 512; i++ {
		dc.LineTo(float64(i), float64(i)*0.001+256)
	}
	dc.Stroke()

	dc.SavePNG("gg.png")
}

func DrawParabola(dc *gg.Context, focX, focY, directrix, startX, endX float64) {
	if math.Abs(focY-directrix) < math.SmallestNonzeroFloat64 {
		dc.LineTo(focX, 0)
		dc.LineTo(focX, directrix)
		dc.Stroke()
		return
	}

	tq := 1.0 / (focY - directrix)
	a := 0.5 * tq
	b := -focX * tq
	c := (focX*focX + focY*focY - directrix*directrix) * tq

	x0 := startX
	y0 := a*x0*x0 + b*x0 + c
	x2 := endX
	y2 := a*x2*x2 + b*x2 + c
	x1 := (x0 + x2) * 0.5
	y1 := a*x0*x2 + 0.5*b*(x0+x2) + c
	points := gg.QuadraticBezier(x0, y0, x1, y1, x2, y2)

	for _, p := range points {
		dc.LineTo(p.X, p.Y)
	}
	dc.Stroke()
}
