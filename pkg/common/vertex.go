package common

import "fmt"

type Vertex struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Edges     []Edge `json:"edges"`
}

func NewVertex(id int64, name string, latitude string, longitude string) *Vertex {
	return &Vertex{
		ID:        id,
		Name:      name,
		Edges:     []Edge{},
		Latitude:  latitude,
		Longitude: longitude,
	}
}

func (v *Vertex) AddNewEdge(edge Edge) {
	v.Edges = append(v.Edges, edge)
}

func (v *Vertex) Print() {
	fmt.Printf("ID: %d Name: %s\n", v.ID, v.Name)
	for _, edge := range v.Edges {
		edge.Print()
	}
}
