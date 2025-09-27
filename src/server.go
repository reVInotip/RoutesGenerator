package main

import (
	"RoutesGenerator/db"
	"RoutesGenerator/utils"
	"fmt"
	"os"
)

func main() {
	var pgr db.PGRoutingQueries = db.PGRoutingQueries{}
	pgr.EstablishConnection("postgres://admin:ab4dsF5hpli1@localhost:5432/route_builder")// standard vars will be used
	defer pgr.FuckingDestroyConnection()

	points := []utils.Point{
		{Lat: 55.759471, Lon: 37.616917},
		{Lat: 55.743097, Lon: 37.614078},
		{Lat: 55.730466, Lon: 37.604430}}

	id, geomStr, length := pgr.BuildRout(&points)
	if id < 0 {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", id)
		return
	}

	fmt.Println(geomStr)
	fmt.Println(length)

	err := os.WriteFile("./utils/coords.json", []byte(geomStr), 0644)
    if err != nil {
        fmt.Println("Write error:", err)
        return
    }
}