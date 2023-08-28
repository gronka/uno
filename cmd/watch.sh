#!/bin/bash
name=$1
check=$(./check-uf.sh $name)
if [ "$check" != "$name" ]; then
	echo $check
	exit 1
fi

cd ../

export Environment=local
air --build.cmd "go build -o build/${name}_bin cmd/$name/main.go" \
	--build.bin "./build/${name}_bin"
