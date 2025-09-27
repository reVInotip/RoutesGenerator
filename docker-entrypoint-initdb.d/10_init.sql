-- Создаем основную таблицу для маршрутизации
CREATE TABLE IF NOT EXISTS routing_roads (
    id SERIAL PRIMARY KEY,
    osm_id BIGINT,
    name TEXT,
    highway TEXT,
    geom GEOMETRY(LineString, 4326),
    source INTEGER,
    target INTEGER,
    x1 DOUBLE PRECISION,
    y1 DOUBLE PRECISION,
    x2 DOUBLE PRECISION,
    y2 DOUBLE PRECISION,
    cost DOUBLE PRECISION,
    reverse_cost DOUBLE PRECISION
);

-- Заполняем таблицу данными дорог из OSM
INSERT INTO routing_roads (osm_id, name, highway, geom)
SELECT 
    osm_id, 
    name, 
    highway, 
    ST_Transform(way, 4326) AS geom
FROM planet_osm_line
WHERE highway IN (
    'motorway', 'trunk', 'primary', 'secondary',
    'tertiary', 'unclassified', 'residential',
    'motorway_link', 'trunk_link', 'primary_link',
    'secondary_link', 'tertiary_link'
);

/* -- create the vertices table */
SELECT * INTO routing_roads_vertices_pgr
FROM pgr_extractVertices('SELECT id, geom FROM routing_roads ORDER BY id');
/* -- set the source information */
UPDATE routing_roads AS e
SET source = v.id, x1 = x, y1 = y
FROM routing_roads_vertices_pgr AS v
WHERE ST_StartPoint(e.geom) = v.geom;
/* -- set the target information */
UPDATE routing_roads AS e
SET target = v.id, x2 = x, y2 = y
FROM routing_roads_vertices_pgr AS v
WHERE ST_EndPoint(e.geom) = v.geom;

-- Рассчитываем стоимость проезда
UPDATE routing_roads SET
    cost = ST_Length(geom::geography) / 
        CASE 
            WHEN highway IN ('motorway','trunk') THEN 130.0 -- 130 км/ч
            WHEN highway IN ('primary') THEN 90.0
            WHEN highway IN ('secondary') THEN 70.0
            ELSE 50.0 -- городские дороги
        END,
    reverse_cost = ST_Length(geom::geography) / 
        CASE 
            WHEN highway IN ('motorway','trunk','motorway_link','trunk_link') THEN 130.0
            WHEN highway IN ('primary','primary_link') THEN 90.0
            WHEN highway IN ('secondary','secondary_link') THEN 70.0
            ELSE 50.0 -- городские дороги
        END;

-- Настраиваем односторонние дороги
UPDATE routing_roads 
SET reverse_cost = -1 
WHERE highway IN ('motorway_link', 'trunk_link');

-- Создаем индексы для ускорения работы
CREATE INDEX IF NOT EXISTS routing_roads_geom_idx ON routing_roads USING GIST(geom);
CREATE INDEX IF NOT EXISTS routing_roads_source_idx ON routing_roads(source);
CREATE INDEX IF NOT EXISTS routing_roads_target_idx ON routing_roads(target);

-- Создание таблицы маршрутов
CREATE TABLE IF NOT EXISTS routes (
    route_id SERIAL PRIMARY KEY,
    route_name VARCHAR(255) NOT NULL,
    geom GEOMETRY(MultiLineString, 4326), -- Геометрия маршрута
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX routes_geom_idx ON routes USING GIST (geom);
CREATE INDEX routes_route_name_idx ON routes(route_name);

CREATE TABLE route_segments (
    route_id INT REFERENCES routes(route_id) ON DELETE CASCADE,
    edge_id BIGINT NOT NULL, -- ID из таблицы routing_roads
    seq_order INT NOT NULL, -- Порядок следования сегментов
    PRIMARY KEY (route_id, edge_id, seq_order)
);