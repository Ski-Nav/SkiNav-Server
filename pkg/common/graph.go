package common

type Graph struct {
	Verticies []*Vertex `json:"vertices"`
}

func NewGraph() *Graph {
	new := &Graph{Verticies: []*Vertex{}}
	return new
}

func (g *Graph) AddNewVertex(vertex *Vertex) {
	g.Verticies = append(g.Verticies, vertex)
}

func (g *Graph) GetVertex(id int64) *Vertex {
	return g.Verticies[id]
}
