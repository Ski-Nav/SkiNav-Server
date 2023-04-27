package db

import (
	"context"

	"github.com/Ski-Nav/SkiNav-Server/pkg/common"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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
	_, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Code within this function might be invoked more than once in case of
		// transient errors.
		readAllLoc := `
			MATCH (n) 
			WHERE n.resort = $resortName
			RETURN n.id, n.name, n.latitude, n.longitude
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
			vertex := common.NewVertex(result.Record().Values[0].(int64), result.Record().Values[1].(string), result.Record().Values[2].(string), result.Record().Values[3].(string))
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
			Match (n1)-[r]->(n2) 
			WHERE r.resort = $resortName
			RETURN n1.id, n2.id, r.name, type(r), r.difficulty, r.weight
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
			fromID, toID := result.Record().Values[0].(int64), result.Record().Values[1].(int64)
			fromV, toV := graph.GetVertex(fromID), graph.GetVertex(toID)
			edgeName := result.Record().Values[2].(string)
			edgeType := result.Record().Values[3].(string)
			edgeDiffculty := result.Record().Values[4].(int64)
			edgeWeight := result.Record().Values[5].(int64)
			edge := common.NewEdge(edgeName, edgeType, toV, edgeDiffculty, edgeWeight)
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
