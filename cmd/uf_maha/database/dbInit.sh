#!/bin/sh
export PGPASSWORD=postgres
export dbName=uf_maha
port=$(cat buildah/PORT)

files=$(ls initdb)
for file in ${files}; do
	echo "--->executing $file"
	psql -h localhost -U postgres -d $dbName -p $port -f "initdb/${file}"
	#psql -h localhost -U postgres -d lbapi_test -f "initdb/${file}"
done
