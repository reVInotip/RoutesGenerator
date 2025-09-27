package db

import (
	"RoutesGenerator/utils"
	"container/list"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type PGRoutingQueries struct {
	Connection *pgx.Conn
}

func (pgr *PGRoutingQueries) EstablishConnection(connString string) {
	var err error

	pgr.Connection, err = pgx.Connect(context.Background(), connString);
	if err != nil {
		panic(err)
	}
}

func (pgr *PGRoutingQueries) FuckingDestroyConnection() {
	pgr.Connection.Close(context.Background())
}

func (pgr *PGRoutingQueries) BuildRout(points *[]utils.Point) (int64, string, float64) {
	var vertices *list.List = list.New();

	var vertexId int64
	for _, point := range *points {
		query := fmt.Sprintf(
			`SELECT id
			FROM routing_roads_vertices_pgr
			WHERE ST_DWithin(
				geom,
				ST_SetSRID(ST_Point(%f, %f), 4326),
				100
			)
			ORDER BY geom <-> ST_SetSRID(ST_Point(%f, %f), 4326)
			LIMIT 1`, point.Lon, point.Lat, point.Lon, point.Lat)

		err := pgr.Connection.QueryRow(context.Background(), query).Scan(&vertexId);
		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			return -1, "", -1
		}

		vertices.PushBack(vertexId)
	}

	if vertices.Len() <= 1 {
		fmt.Fprintf(os.Stderr, "Can not create route of one vertex")
	}

	verticesString := "ARRAY" + ListToString(*vertices)

	query := fmt.Sprintf(
		`WITH dijkstra AS (
			SELECT *
			FROM pgr_dijkstraVia(
				'SELECT id, source, target, cost, reverse_cost FROM routing_roads',
				%s,
				directed := false
			)
		),
		path_edges AS (
			SELECT 
				d.path_id AS PathId,
				d.path_seq AS Sequence,
				d.edge AS EdgeId,
				d.node AS NodeId,
				r.geom AS edge_geom
			FROM dijkstra d
			LEFT JOIN routing_roads r ON d.edge = r.id
			WHERE d.edge > 0
		)
		SELECT 
			PathId,
			ST_AsGeoJSON(ST_LineMerge(ST_Collect(edge_geom))) AS RouteGeometry,
			ST_Length(ST_LineMerge(ST_Collect(edge_geom))) AS RouteLength
		FROM path_edges
		GROUP BY PathId
		ORDER BY PathId`, verticesString)
	
	
	var pathId int64
	var routeGeom string
	var routeLength float64
	err := pgr.Connection.QueryRow(context.Background(), query).Scan(&pathId, &routeGeom, &routeLength);
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return -1, "", -1
	}

	return pathId, routeGeom, routeLength
}

func ListToString(list list.List) string {
	var builder strings.Builder = strings.Builder{}
	builder.WriteString("[")

	item := list.Front()
	for ; ; {
		switch val := item.Value.(type) {
		case int64:
			_, _ = builder.WriteString(strconv.FormatInt(val, 10))
		default:
			fmt.Println("Can not convert list item to string")
			return ""
		}
		

		item = item.Next()
		if item == nil {
			builder.WriteString("]")
			break
		} else {
			builder.WriteString(", ")
		}
	}

	return builder.String()
}