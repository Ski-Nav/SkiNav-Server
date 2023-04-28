package maps

import (
	"errors"

	"github.com/Ski-Nav/SkiNav-Server/pkg/common"
	"github.com/Ski-Nav/SkiNav-Server/pkg/lib/db"
)

type I interface {
	GetGraphByResortName(name string) (*common.Graph, error)
	GetAllResorts() *[]string
}

type ResortMap struct {
	Map        map[string]*common.Graph
	AllResorts *[]string
}

func Init(db *db.DB) I {
	resortMap := make(map[string]*common.Graph)
	allResorts := db.GetAllResort()
	for _, resort := range *allResorts {
		resortMap[resort] = db.GetGraphByResort(resort)
	}
	return &ResortMap{
		Map:        resortMap,
		AllResorts: allResorts,
	}
}

func (m *ResortMap) GetAllResorts() *[]string {
	return m.AllResorts
}

func (m *ResortMap) GetGraphByResortName(name string) (*common.Graph, error) {
	graph, ok := m.Map[name]
	if !ok {
		return nil, errors.New("resort not found")
	}
	return graph, nil
}
