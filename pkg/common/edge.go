package common

import "fmt"

type Edge struct {
	Name       string  `json:"name"`
	EdgeType   string  `json:"edgeType"`
	To         *Vertex `json:"-"`
	ToID       int64   `json:"to"`
	Difficulty int64   `json:"difficulty"`
	Weight     int64   `json:"weight"`
}

func NewEdge(name string, edgeType string, to *Vertex, difficulty int64, weight int64) *Edge {
	return &Edge{
		Name:       name,
		EdgeType:   edgeType,
		To:         to,
		ToID:       to.ID,
		Difficulty: difficulty,
		Weight:     weight,
	}
}

func (e *Edge) Print() {
	fmt.Printf("Type: %s To: %d\n", e.EdgeType, e.To.ID)
}
