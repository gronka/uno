#!/bin/sh
name=$(cat buildah/CONTAINERNAME)_local
podman stop $name
