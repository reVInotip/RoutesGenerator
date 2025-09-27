#!/bin/bash

exec_psql() {
	PGPASSWORD=${POSTGRES_PASSWORD} psql -v ON_ERROR_STOP=1 -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" <<-EOSQL
        $1
EOSQL
}

exec_psql "GRANT ALL PRIVILEGES ON DATABASE ${POSTGRES_DB} TO ${POSTGRES_USER};"

exec_psql "CREATE EXTENSION IF NOT EXISTS postgis;"
exec_psql "CREATE EXTENSION IF NOT EXISTS hstore;"
exec_psql "CREATE EXTENSION IF NOT EXISTS pgrouting;"

touch /tmp/flat_node

# Импорт данных в PostgreSQL
for map in $(ls /tmp/maps)
do
echo $map
nohup osm2pgsql -F /tmp/flat_node -d $POSTGRES_DB -U $POSTGRES_USER --cache=1000 --number-processes=4 --create --multi-geometry --slim --drop --hstore --proj 4326 /tmp/maps/$map
done

cat > "$PGDATA/pg_hba.conf" <<EOF
    # TYPE  DATABASE  USER  ADDRESS  METHOD
    local   all       all             scram-sha-256
    host    ${POSTGRES_DB}	${POSTGRES_USER}		   all			   scram-sha-256
EOF