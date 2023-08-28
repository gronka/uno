#!/bin/sh
port=$(cat buildah/PORT)

echo 'common commands:'
echo CREATE DATABASE papi;
echo 'list tables: \dt;'
echo '\c papi'
echo 'select * from users;'

echo UPDATE users SET status=50 WHERE id='2f617a75-...';

# connect from local psql
#psql -d papi -U postgres -h localhost -p 7002

# connect from psql inside the postgresql container
#sudo podman exec -it papi_container_1 /bin/bash -c "psql -d papi -U postgres -h localhost"

# -d papi to connect to database
#psql -U postgres -h localhost -p $port -d uf_user -W
psql postgresql://postgres:postgres@localhost:$port/uf_order
