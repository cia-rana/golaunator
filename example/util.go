package main

import "math/rand"

const (
	vertexDomainMax = 1000
	verticesNum     = 500
)

var vertices = [][]float64{}

func init() {
	rand.Seed(4)
	if !true {
		vertices = [][]float64{
			{0, 0},
			{100, 500},
			{500, 50},
			{125, 125},
		}
	} else {
		// Create unique vertices.
		tmpVertices := make(map[int]struct{})
		for i := 0; len(tmpVertices) < verticesNum; i++ {
			x, y := rand.Intn(vertexDomainMax), rand.Intn(vertexDomainMax)
			tmpVertices[x*vertexDomainMax+y] = struct{}{}
		}

		vertices = make([][]float64, 0, verticesNum)
		for k := range tmpVertices {
			vertices = append(vertices, []float64{float64(k / vertexDomainMax), float64(k % vertexDomainMax)})
		}
	}
}
