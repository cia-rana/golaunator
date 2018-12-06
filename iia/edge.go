package iia

type HalfEdge struct {
	End      *Vertex
	Pair     *HalfEdge
	Next     *HalfEdge
	Triangle *Triangle
}

type HalfEdgeStack []*HalfEdge

func NewHalfEdgeStack() *HalfEdgeStack {
	hes := HalfEdgeStack([]*HalfEdge{})
	return &hes
}

func (hes *HalfEdgeStack) Pop() *HalfEdge {
	if hes.IsEmpty() {
		return nil
	}

	preHes := *hes
	n := len(preHes)
	he := preHes[n-1]
	*hes = preHes[0 : n-1]
	return he
}

func (hes *HalfEdgeStack) Push(he *HalfEdge) {
	*hes = append(*hes, he)
}

func (hes HalfEdgeStack) IsEmpty() bool {
	return len(hes) == 0
}
