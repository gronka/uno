#!/bin/sh

name=uf_aim_pg

podman stop $name
podman rm -f $name
podman rm -f ${name}_local
podman rmi -f $name
buildah from --name $name docker.io/postgres:14.5-alpine

echo config
buildah config \
	--env POSTGRES_USER=postgres \
	--env POSTGRES_PASSWORD=postgres \
	--env POSTGRES_DB=uf_aim \
	$name

#buildah copy $name 00-db-init.sql /docker-entrypoint-initdb.d/
buildah config --port 5432 $name

buildah config --created-by "textFridayy" $name
buildah config --author "Taylor" $name
buildah config --label name=$name $name
buildah commit $name $name

podman tag $name "localhost:5000/$name"
podman push "localhost:5000/$name"
