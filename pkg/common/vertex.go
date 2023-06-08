package common

type Vertex struct {
	ID        int64         `json:"id"`
	Name      string        `json:"name"`
	Latitude  float64       `json:"latitude"`
	Longitude float64       `json:"longitude"`
	Edges     []Edge        `json:"edges"`
	Aliases   []interface{} `json:"aliases"`
}

func NewVertex(id int64, name string, latitude float64, longitude float64, aliases []interface{}) *Vertex {
	return &Vertex{
		ID:        id,
		Name:      name,
		Edges:     []Edge{},
		Latitude:  latitude,
		Longitude: longitude,
		Aliases:   aliases,
	}
}

func (v *Vertex) AddNewEdge(edge Edge) {
	v.Edges = append(v.Edges, edge)
}
