package common

type Edge struct {
	Name       string  `json:"name"`
	ID         string  `json:"id"`
	EdgeType   string  `json:"edgeType"`
	ToID       int64   `json:"to"`
	Difficulty string  `json:"difficulty"`
	Weight     float64 `json:"weight"`
}

func NewEdge(name string, id string, edgeType string, to *Vertex, difficulty string, weight float64) *Edge {
	return &Edge{
		Name:       name,
		ID:         id,
		EdgeType:   edgeType,
		ToID:       to.ID,
		Difficulty: difficulty,
		Weight:     weight,
	}
}
