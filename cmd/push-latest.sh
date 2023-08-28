#!/bin/sh
name=$1
check=$(./check-uf.sh $name)
if [ "$check" != "$name" ]; then
	echo $check
	exit 1
fi

buildah tag "localhost:5000/$name" "docker.io/gronka/$name"
podman push "docker.io/gronka/$name"

echo "pushed service: $name:latest"
