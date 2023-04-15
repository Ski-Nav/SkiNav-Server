package common

import "fmt"

type Edge struct {
	EdgeType string  `json:"edgeType"`
	To       *Vertex `json:"-"`
	ToID     int64   `json:"to"`
}

func NewEdge(edgeType string, to *Vertex) *Edge {
	return &Edge{EdgeType: edgeType, To: to, ToID: to.ID}
}

func (e *Edge) Print() {
	fmt.Printf("Type: %s To: %d\n", e.EdgeType, e.To.ID)
}
