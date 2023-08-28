#!/bin/sh
name=$(cat buildah/CONTAINERNAME)_local
podman start $name
