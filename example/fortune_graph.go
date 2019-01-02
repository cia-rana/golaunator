package main

import (
	"fmt"
	"github.com/awalterschulze/gographviz"
	"github.com/cia-rana/golaunator/fortune"
)

func main() {
	verticesForFortune := make([]*fortune.Vertex, 0, len(vertices))
	for _, v := range vertices {
		verticesForFortune = append(verticesForFortune, &fortune.Vertex{X: v[0], Y: v[1]})
	}

	t := fortune.NewTriangulator(verticesForFortune)
	t.EnableDraw()
	t.Triangulate()

	if true {
		g := gographviz.NewGraph()

		gName := "G"
		g.SetName(gName)
		g.AddAttr(gName, "splines", "line")

		for i, v := range verticesForFortune {
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
