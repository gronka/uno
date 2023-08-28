#!/bin/sh
name=$1
check=$(./check-uf.sh $name)
if [ "$check" != "$name" ]; then
	echo $check
	exit 1
fi

version=$(cat VERSION)
buildah tag "localhost/$name" "docker.io/gronka/$name:$version"
podman push "docker.io/gronka/$name:$version"
