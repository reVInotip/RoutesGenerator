package main

import (
	"encoding/json"
	"RoutesGenerator/db"
	"RoutesGenerator/utils"
	"fmt"
	"math/rand/v2"
	"os"
)

type Route struct {
	Id int64 `json:"id"`
	Geometry json.RawMessage `json:"geom"`
	Length float64 `json:"length"`
}

func randomLat(minLat, maxLat float64) float64 {
	return minLat + rand.Float64()*(maxLat-minLat)
}

func randomLon(minLon, maxLon float64) float64 {
	return minLon + rand.Float64()*(maxLon-minLon)
}

func main() {
	var pgr db.PGRoutingQueries = db.PGRoutingQueries{}
	// this is bad practice but it is only test project
	pgr.EstablishConnection("postgres://admin:ab4dsF5hpli1@localhost:5432/route_builder")
	defer pgr.FuckingDestroyConnection()

	// Moscow bounds
	maxLat := 56.021389;
	minLat := 55.143833;
	maxLon := 37.967778;
	minLon := 36.80325;

	var routes [1000]Route = [1000]Route{};

	k := 0;
	for range 1000 {
		points := []utils.Point{
			{Lat: randomLat(minLat, maxLat), Lon: randomLon(minLon, maxLon)},
			{Lat: randomLat(minLat, maxLat), Lon: randomLon(minLon, maxLon)},
			{Lat: randomLat(minLat, maxLat), Lon: randomLon(minLon, maxLon)}}
		
		fmt.Println(points)

		id, geomStr, length := pgr.BuildRout(&points)
		if id < 0 {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", id)
			return
		}

		fmt.Println(geomStr)
		fmt.Println(length)

		if id > 0 {
			routes[k] = Route{
				Id: id,
				Geometry: json.RawMessage(geomStr),
				Length: length,
			}
			k++;
		}
	}
    
    jsonPretty, err := json.MarshalIndent(routes[:k], "", "  ")
    if err != nil {
        panic(err)
    }
	
    fmt.Println(string(jsonPretty))
	err = os.WriteFile("./utils/coords.json", []byte(jsonPretty), 0644);
	if err != nil {
		fmt.Println("Write error:", err)
		return
	}
}