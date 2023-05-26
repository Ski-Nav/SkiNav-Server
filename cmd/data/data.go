package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Ski-Nav/SkiNav-Server/pkg/data"
	"github.com/Ski-Nav/SkiNav-Server/pkg/lib/db"
	"github.com/umahmood/haversine"
)

type location struct {
	lat  float64
	long float64
}

func main() {
	db := db.Init()
	entries, err := os.ReadDir("./data")
	if err != nil {
		log.Fatal(err)
	}
	locToId := make(map[location]string)
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		resort := e.Name()
		log.Printf("delete existing nodes")
		err := db.DeleteNodesAndWays(resort)
		if err != nil {
			log.Fatal(err)
		}
		filePath := path.Join("./data", resort, "pisteNode.geojson")
		nodes := data.ExtractNode(filePath)
		log.Printf("%d piste nodes", len(nodes.Features))
		noNameCount := 1
		for i, f := range nodes.Features {
			if i%100 == 0 {
				log.Printf("%d done", i)
			}
			loc := location{long: f.Geometry.Point[0], lat: f.Geometry.Point[1]}
			if _, ok := locToId[loc]; !ok {
				locToId[loc] = f.Properties["id"].(string)
				err := db.InsertNode(f, resort)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		filePath = path.Join("./data", resort, "pisteWay.geojson")
		ways := data.ExtractWay(filePath)
		log.Printf("%d piste ways", len(ways.Features))
		for i, f := range ways.Features {
			if i%10 == 0 {
				log.Printf("%d done", i)
			}
			for i := 0; i < len(f.Geometry.LineString)-1; i++ {
				fromLoc := location{long: f.Geometry.LineString[i][0], lat: f.Geometry.LineString[i][1]}
				fromId, ok := locToId[fromLoc]
				if !ok {
					log.Fatal("node not found")
				}
				toLoc := location{long: f.Geometry.LineString[i+1][0], lat: f.Geometry.LineString[i+1][1]}
				toId, ok := locToId[toLoc]
				if !ok {
					log.Fatal("node not found")
				}
				mi, _ := haversine.Distance(haversine.Coord{Lat: fromLoc.lat, Lon: fromLoc.long}, haversine.Coord{Lat: toLoc.lat, Lon: toLoc.long})
				if _, ok := f.Properties["name"]; !ok {
					f.Properties["name"] = fmt.Sprintf("no name piste %d", noNameCount)
					noNameCount += 1
				}
				if _, ok := f.Properties["piste:difficulty"]; !ok {
					f.Properties["piste:difficulty"] = "novice"
				}
				err := db.InsertPiste(f, resort, fromId, toId, mi)
				if err != nil {
					log.Fatal(err)
				}
			}
			toploc := location{long: f.Geometry.LineString[0][0], lat: f.Geometry.LineString[0][1]}
			topId := locToId[toploc]
			err := db.AddNodeAlias(topId, fmt.Sprintf("top/%s", f.Properties["name"]))
			if err != nil {
				log.Fatal(err)
			}
			lastIdx := len(f.Geometry.LineString) - 1
			bottomloc := location{long: f.Geometry.LineString[lastIdx][0], lat: f.Geometry.LineString[lastIdx][1]}
			bottomId := locToId[bottomloc]
			err = db.AddNodeAlias(bottomId, fmt.Sprintf("bottom/%s", f.Properties["name"]))
			if err != nil {
				log.Fatal(err)
			}
		}

		filePath = path.Join("./data", resort, "liftNode.geojson")
		nodes = data.ExtractNode(filePath)
		log.Printf("%d lift nodes", len(nodes.Features))
		for i, f := range nodes.Features {
			if i%100 == 0 {
				log.Printf("%d done", i)
			}
			loc := location{long: f.Geometry.Point[0], lat: f.Geometry.Point[1]}
			if _, ok := locToId[loc]; !ok {
				locToId[loc] = f.Properties["id"].(string)
				err := db.InsertNode(f, resort)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		filePath = path.Join("./data", resort, "liftWay.geojson")
		ways = data.ExtractWay(filePath)
		log.Printf("%d liftWay ways", len(ways.Features))
		noNameCount = 1
		for i, f := range ways.Features {
			if i%10 == 0 {
				log.Printf("%d done", i)
			}
			if _, ok := f.Properties["name"]; !ok {
				f.Properties["name"] = fmt.Sprintf("no name lift %d", noNameCount)
				noNameCount += 1
			}
			for i := 0; i < len(f.Geometry.LineString)-1; i++ {
				fromLoc := location{long: f.Geometry.LineString[i][0], lat: f.Geometry.LineString[i][1]}
				fromId, ok := locToId[fromLoc]
				if !ok {
					log.Fatal("node not found")
				}
				toLoc := location{long: f.Geometry.LineString[i+1][0], lat: f.Geometry.LineString[i+1][1]}
				toId, ok := locToId[toLoc]
				if !ok {
					log.Fatal("node not found")
				}
				mi, _ := haversine.Distance(haversine.Coord{Lat: fromLoc.lat, Lon: fromLoc.long}, haversine.Coord{Lat: toLoc.lat, Lon: toLoc.long})
				err := db.InsertLift(f, resort, fromId, toId, mi)
				if err != nil {
					log.Fatal(err)
				}
			}
			toploc := location{long: f.Geometry.LineString[0][0], lat: f.Geometry.LineString[0][1]}
			topId := locToId[toploc]
			err := db.AddNodeAlias(topId, fmt.Sprintf("top/%s", f.Properties["name"]))
			if err != nil {
				log.Fatal(err)
			}
			lastIdx := len(f.Geometry.LineString) - 1
			bottomloc := location{long: f.Geometry.LineString[lastIdx][0], lat: f.Geometry.LineString[lastIdx][1]}
			bottomId := locToId[bottomloc]
			err = db.AddNodeAlias(bottomId, fmt.Sprintf("bottom/%s", f.Properties["name"]))
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
