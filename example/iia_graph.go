package main

import (
	"fmt"
	"github.com/awalterschulze/gographviz"
	"github.com/cia-rana/golaunator/iia"
)

func main() {
	verticesForIaa := make([]*iia.Vertex, 0, len(vertices))
	for _, v := range vertices {
		verticesForIaa = append(verticesForIaa, &iia.Vertex{X: v[0], Y: v[1]})
	}

	t := iia.NewTriangulator(verticesForIaa)
	t.Triangulate()

	if true {
		g := gographviz.NewGraph()

		gName := "G"
		g.SetName(gName)
		g.SetDir(true)
		g.AddAttr(gName, "splines", "line")

		//verticesForIaa = append(verticesForIaa, iia.ComputeCoverTriangleVertices(verticesForIaa)...)

		for i, v := range verticesForIaa {
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
