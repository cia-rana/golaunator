package main

import (
	"fmt"
	"github.com/awalterschulze/gographviz"
	"github.com/cia-rana/golaunator/iia"
)

func main() {
	vertices := []*iia.Vertex{
		{X: 0, Y: 0},
		{X: 100, Y: 500},
		{X: 500, Y: 0},
		{X: 125, Y: 125},
	}
	t := iia.NewTriangulator(vertices)
	t.Triangulate()

	if true {
		g := gographviz.NewGraph()

		gName := "G"
		g.SetName(gName)
		g.SetDir(true)
		g.AddAttr(gName, "splines", "line")

		vertices = append(vertices, iia.ComputeCoverTriangleVertices(vertices)...)

		for i, v := range vertices {
			g.AddNode(
				gName,
				fmt.Sprint(i),
				map[string]string{
					"pos": "\"" + fmt.Sprint(v.X) + "," + fmt.Sprint(-v.Y) + "\"",
				},
			)
		}

		for _, he := range t.HalfEdges {
			g.AddEdge(fmt.Sprint(he.Next.Next.End.Index), fmt.Sprint(he.End.Index), true, nil)
		}

		fmt.Println(g.String())
	} else {
		for _, he := range t.HalfEdges {
			fmt.Println(he.Next.Next.End.Index, "->", he.End.Index)
		}
	}
}
