package data

import (
	"os"

	geojson "github.com/paulmach/go.geojson"
)

func ExtractNode(path string) *geojson.FeatureCollection {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	// Feature Collection
	rawFeatureJSON := data

	fc, err := geojson.UnmarshalFeatureCollection(rawFeatureJSON)
	if err != nil {
		panic(err)
	}
	return fc
}

func ExtractWay(path string) *geojson.FeatureCollection {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	// Feature Collection
	rawFeatureJSON := data

	fc, err := geojson.UnmarshalFeatureCollection(rawFeatureJSON)
	if err != nil {
		panic(err)
	}

	return fc
}
