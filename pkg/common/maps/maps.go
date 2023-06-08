package maps

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/Ski-Nav/SkiNav-Server/pkg/common"
	"github.com/Ski-Nav/SkiNav-Server/pkg/lib/db"
	"github.com/bradfitz/gomemcache/memcache"
)

type I interface {
	GetGraphByResortName(name string) (*common.Graph, error)
	GetAllResorts() *[]string
}

type ResortMap struct {
	Map        map[string]*common.Graph
	AllResorts *[]string
	mc         *memcache.Client
	db         *db.DB
}

func Init(db *db.DB) I {
	resortMap := make(map[string]*common.Graph)
	allResorts := db.GetAllResort()
	for _, resort := range *allResorts {
		resortMap[resort] = db.GetGraphByResort(resort)
	}

	// mc := memcache.New("127.0.0.1:11211")
	mc := memcache.New("memcached.xzdlgx.0001.use2.cache.amazonaws.com:11211")
	return &ResortMap{
		Map:        resortMap,
		AllResorts: allResorts,
		mc:         mc,
		db:         db,
	}
}

func (m *ResortMap) GetAllResorts() *[]string {
	return m.AllResorts
}

func (m *ResortMap) GetGraphByResortName(name string) (*common.Graph, error) {
	// Try to retrieve the data from cache
	item, err := m.mc.Get(name)
	if err == nil {
		fmt.Println("Cache hit!")
		b := bytes.NewReader(item.Value)

		var graph common.Graph

		if err := gob.NewDecoder(b).Decode(&graph); err != nil {
			return &common.Graph{}, err
		}
		return &graph, nil
	}
	if err == memcache.ErrCacheMiss {
		// Cache miss, fetch the data from the database
		fmt.Println("Cache miss!")

		graph := m.db.GetGraphByResort(name)
		var b bytes.Buffer

		if err := gob.NewEncoder(&b).Encode(graph); err != nil {
			return &common.Graph{}, err
		}

		err := m.mc.Set(&memcache.Item{
			Key:        name,
			Value:      b.Bytes(),
			Expiration: int32(time.Now().Add(25 * time.Second).Unix()),
		})
		if err != nil {
			return &common.Graph{}, err
		}

		return graph, nil
	}
	return &common.Graph{}, err
}
