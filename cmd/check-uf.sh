#!/bin/bash
#name=${PWD##*/}
name=$1

if [ $name == "uf_aim" ] || \
	 [ $name == "uf_border" ] || \
	 [ $name == "uf_maha" ] || \
	 [ $name == "uf_order" ] || \
	 [ $name == "uf_public" ] || \
	 [ $name == "uf_user" ]; then
	echo $name
else
	echo invalid microservice: $name
	exit 1
fi
