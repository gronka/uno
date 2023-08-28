#!/bin/bash
echo "don't forget to make sure your local docker repository is running"

echo login to docker:
podman login docker.io

services=("uf_aim" "uf_border" "uf_maha" "uf_order" "uf_public" "uf_user")

for service in ${services[@]}; do
	echo ===== preparing $service =====
	./arch-build.sh $service
	echo finished building $service
	./push-latest.sh $service
	echo finished pushing $service:latest
done
