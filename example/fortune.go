package main

import (
	"fmt"
	"github.com/awalterschulze/gographviz"
	"github.com/cia-rana/golaunator/fortune"
)

func main() {
	vertices := []*fortune.Vertex{
		{X: 0, Y: 0},
		{X: 100, Y: 500},
		{X: 500, Y: 0},
		{X: 125, Y: 125},
	}
	t := fortune.NewTriangulator(vertices)
	t.Triangulate()

	if !true {
		g := gographviz.NewGraph()

		gName := "G"
		g.SetName(gName)
		g.AddAttr(gName, "splines", "line")

		for i, v := range vertices {
			g.AddNode(
				gName,
				fmt.Sprint(i),
				map[string]string{
					"pos": "\"" + fmt.Sprint(v.X) + "," + fmt.Sprint(-v.Y) + "\"",
				},
			)
		}

		for _, e := range t.Edges {
			g.AddEdge(fmt.Sprint(e.Start.Index), fmt.Sprint(e.End.Index), false, nil)
		}

		fmt.Println(g.String())
	} else {
		for _, e := range t.Edges {
			fmt.Println(e.Start.Index, "--", e.End.Index)
		}
	}
}
