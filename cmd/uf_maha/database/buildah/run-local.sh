#!/bin/sh
con=$(cat CONTAINERNAME)
port=$(cat PORT)
name=${con}_local
podman run -dt -p $port:5432 --name $name localhost:5000/$con
