package example

import (
	"math/rand"
	"os"
	"testing"

	"github.com/cia-rana/golaunator/fortune"
	"github.com/cia-rana/golaunator/iia"
)

const (
	vertexDomainMax = 1000
	verticesNum     = 400
)

var vertices = make([][]float64, 0, verticesNum)

func TestMain(m *testing.M) {
	//rand.Seed(time.Now().UnixNano())
	rand.Seed(4)
	for i := 0; i < verticesNum; i++ {
		vertices = append(vertices, []float64{float64(rand.Intn(vertexDomainMax)), float64(rand.Intn(vertexDomainMax))})
	}

	os.Exit(m.Run())
}

func TestEdgesMatchInAllAlgorithms(t *testing.T) {
	verticesForFortune := make([]*fortune.Vertex, 0, len(vertices))
	verticesForIaa := make([]*iia.Vertex, 0, len(vertices))
	for _, v := range vertices {
		verticesForFortune = append(verticesForFortune, &fortune.Vertex{X: v[0], Y: v[1]})
		verticesForIaa = append(verticesForIaa, &iia.Vertex{X: v[0], Y: v[1]})
	}

	triangulatorForFortune := fortune.NewTriangulator(verticesForFortune)
	triangulatorForFortune.Triangulate()

	triangulatorForIaa := iia.NewTriangulator(verticesForIaa)
	triangulatorForIaa.Triangulate()

	edgesByFortune := make(map[int]struct{})
	edgesByIaa := make(map[int]struct{})
	for _, e := range triangulatorForFortune.Edges {
		a, b := IntMax(e.Start.Index, e.End.Index)
		t.Log(a, b)
		edgesByFortune[a*vertexDomainMax+b] = struct{}{}
	}
	for _, e := range triangulatorForIaa.HalfEdges {
		a, b := IntMax(e.Next.Next.End.Index, e.End.Index)
		t.Log(a, b)
		edgesByIaa[a*vertexDomainMax+b] = struct{}{}
	}
	/*
		if len(edgesByFortune) != len(edgesByIaa) {
			t.Fatalf("edgesByFortune: %d, edgesByIaa: %d", len(edgesByFortune), len(edgesByIaa))
		}
	*/
	for key, _ := range edgesByFortune {
		if k, ok := edgesByIaa[key]; !ok {
			t.Logf("%d, %d", key, k)
		}
	}
	t.Fatalf("")
}

func IntMax(a, b int) (int, int) {
	if a > b {
		return b, a
	}
	return a, b
}
