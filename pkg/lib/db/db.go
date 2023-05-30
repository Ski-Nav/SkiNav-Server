package db

import (
	"context"

	"github.com/Ski-Nav/SkiNav-Server/pkg/common"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	geojson "github.com/paulmach/go.geojson"
)

type DB struct {
	driver *neo4j.DriverWithContext
	ctx    context.Context
}

func Init() *DB {
	// Aura requires you to use "neo4j+s" protocol
	// (You need to replace your connection details, username and password)
	uri := "neo4j+s://f4c8b933.databases.neo4j.io:7687"
	auth := neo4j.BasicAuth("neo4j", "PjfCV1nixrFvsPTaqmJAkC1-FX4BNlSFIONpw2dARnc", "")
	// You typically have one driver instance for the entire application. The
	// driver maintains a pool of database connections to be used by the sessions.
	// The driver is thread safe.
	driver, err := neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		panic(err)
	}

	return &DB{
		driver: &driver,
		ctx:    context.Background(),
	}
}

func (db *DB) GetGraphByResort(resortName string) *common.Graph {
	driver := *db.driver
	ctx := db.ctx
	// Don't forget to close the driver connection when you are finished with it
	// Create a session to run transactions in. Sessions are lightweight to
	// create and use. Sessions are NOT thread safe.
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)
	// Now read the created persons. By using ExecuteRead method a connection
	// to a read replica can be used which reduces load on writer nodes in cluster.
	graph := common.NewGraph()
	idToIdx := make(map[string]int64, 0)
	var idx int64 = 0
	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Code within this function might be invoked more than once in case of
		// transient errors.
		readAllLoc := `
			MATCH (n) 
			WHERE n.resort = $resortName
			RETURN n.id, n.latitude, n.longitude, n.aliases
			ORDER BY n.id
		`
		result, err := tx.Run(ctx, readAllLoc, map[string]any{
			"resortName": resortName,
		})
		if err != nil {
			return nil, err
		}
		// Iterate over the result within the transaction instead of using
		// Collect (just to show how it looks...). Result.Next returns true
		// while a record could be retrieved, in case of error result.Err()
		// will return the error.
		for result.Next(ctx) {
			id := result.Record().Values[0].(string)
			idToIdx[id] = idx
			idx++
			vertex := common.NewVertex(idToIdx[id], id, result.Record().Values[1].(float64), result.Record().Values[2].(float64), result.Record().Values[3].([]interface{}))
			graph.AddNewVertex(vertex)
		}
		// Again, return any error back to driver to indicate rollback and
		// retry in case of transient error.
		return nil, result.Err()
	})
	if err != nil {
		panic(err)
	}

	_, err = session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Code within this function might be invoked more than once in case of
		// transient errors.
		readAllLoc := `
			Match (n1)-[r:SLOPE]->(n2) 
			WHERE r.resort = $resortName
			RETURN n1.id, n2.id, r.name, type(r), r.difficulty, r.distance, r.id
		`
		result, err := tx.Run(ctx, readAllLoc, map[string]any{
			"resortName": resortName,
		})
		if err != nil {
			return nil, err
		}
		// Iterate over the result within the transaction instead of using
		// Collect (just to show how it looks...). Result.Next returns true
		// while a record could be retrieved, in case of error result.Err()
		// will return the error.
		for result.Next(ctx) {
			fromID, toID := idToIdx[result.Record().Values[0].(string)], idToIdx[result.Record().Values[1].(string)]
			fromV, toV := graph.GetVertex(fromID), graph.GetVertex(toID)
			edgeName := result.Record().Values[2].(string)
			edgeType := result.Record().Values[3].(string)
			edgeDiffculty := result.Record().Values[4].(string)
			edgeWeight := result.Record().Values[5].(float64)
			edgeID := result.Record().Values[6].(string)
			edge := common.NewEdge(edgeName, edgeID, edgeType, toV, edgeDiffculty, edgeWeight)
			fromV.AddNewEdge(*edge)
		}
		// Again, return any error back to driver to indicate rollback and
		// retry in case of transient error.
		return nil, result.Err()
	})
	if err != nil {
		panic(err)
	}

	_, err = session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Code within this function might be invoked more than once in case of
		// transient errors.
		readAllLoc := `
			Match (n1)-[r:LIFT]->(n2) 
			WHERE r.resort = $resortName
			RETURN n1.id, n2.id, r.name, type(r), r.distance, r.id
		`
		result, err := tx.Run(ctx, readAllLoc, map[string]any{
			"resortName": resortName,
		})
		if err != nil {
			return nil, err
		}
		// Iterate over the result within the transaction instead of using
		// Collect (just to show how it looks...). Result.Next returns true
		// while a record could be retrieved, in case of error result.Err()
		// will return the error.
		for result.Next(ctx) {
			fromID, toID := idToIdx[result.Record().Values[0].(string)], idToIdx[result.Record().Values[1].(string)]
			fromV, toV := graph.GetVertex(fromID), graph.GetVertex(toID)
			edgeName := result.Record().Values[2].(string)
			edgeType := result.Record().Values[3].(string)
			edgeWeight := result.Record().Values[4].(float64)
			edgeID := result.Record().Values[5].(string)
			edge := common.NewEdge(edgeName, edgeID, edgeType, toV, "", edgeWeight)
			fromV.AddNewEdge(*edge)
		}
		// Again, return any error back to driver to indicate rollback and
		// retry in case of transient error.
		return nil, result.Err()
	})
	if err != nil {
		panic(err)
	}
	return graph
}

func (db *DB) GetAllResort() *[]string {
	driver := *db.driver
	ctx := db.ctx
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)
	allResort := make([]string, 0)
	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		readAllLoc := `
			match (n) 
			return distinct(n.resort)
		`
		result, err := tx.Run(ctx, readAllLoc, map[string]any{})
		if err != nil {
			return nil, err
		}
		for result.Next(ctx) {
			resort := result.Record().Values[0].(string)
			allResort = append(allResort, resort)
		}
		return nil, result.Err()
	})
	if err != nil {
		panic(err)
	}
	return &allResort
}

func (db *DB) InsertNode(node *geojson.Feature, resort string) error {
	driver := *db.driver
	ctx := db.ctx
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(ctx,
			`CREATE (n:Node {id: $id, longitude: $longitude, latitude: $latitude, resort: $resort, aliases: []}) 
			RETURN n.name + ', from node ' + id(n)`,
			map[string]any{
				"id":        node.Properties["id"],
				"longitude": node.Geometry.Point[0],
				"latitude":  node.Geometry.Point[1],
				"resort":    resort,
			})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) InsertPiste(node *geojson.Feature, resort string, fromId string, toId string, distance float64) error {
	driver := *db.driver
	ctx := db.ctx
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(ctx,
			`MATCH
			(a:Node),
			(b:Node)
		  	WHERE a.id = $fromId AND b.id = $toId
		  	CREATE (a)-[r:SLOPE {resort: $resort, name: $name, difficulty: $difficulty, distance: $distance, id: $id}]->(b)
		  	RETURN type(r)`,
			map[string]any{
				"name":       node.Properties["name"],
				"id":         node.Properties["id"],
				"fromId":     fromId,
				"toId":       toId,
				"resort":     resort,
				"difficulty": node.Properties["piste:difficulty"],
				"distance":   distance,
			})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) InsertLift(node *geojson.Feature, resort string, fromId string, toId string, distance float64) error {
	driver := *db.driver
	ctx := db.ctx
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(ctx,
			`MATCH
			(a:Node),
			(b:Node)
		  	WHERE a.id = $fromId AND b.id = $toId
		  	CREATE (a)-[r:LIFT {resort: $resort, name: $name, distance: $distance, id: $id}]->(b)
		  	RETURN type(r)`,
			map[string]any{
				"name":     node.Properties["name"],
				"id":       node.Properties["id"],
				"fromId":   fromId,
				"toId":     toId,
				"resort":   resort,
				"distance": distance,
			})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return err
	}
	return nil
}

func (db *DB) DeleteNodesAndWays(resort string) error {
	driver := *db.driver
	ctx := db.ctx
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(ctx,
			`MATCH (a)-[r]->(b)
			WHERE r.resort = $resort
			DELETE r`,
			map[string]any{
				"resort": resort,
			})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return err
	}
	_, err = session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(ctx,
			`MATCH (n)
			WHERE n.resort = $resort
			DELETE n`,
			map[string]any{
				"resort": resort,
			})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) AddNodeAlias(id string, alias string) error {
	driver := *db.driver
	ctx := db.ctx
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, func(transaction neo4j.ManagedTransaction) (any, error) {
		result, err := transaction.Run(ctx,
			`MATCH (n:Node)
		  	WHERE n.id = $id
			SET n.aliases = n.aliases + $alias`,
			map[string]any{
				"id":    id,
				"alias": alias,
			})
		if err != nil {
			return nil, err
		}

		if result.Next(ctx) {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return err
	}
	return nil
}
