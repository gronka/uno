#!/bin/bash

cd global
kubectl delete -f global-config-prod.yaml
kubectl apply -f global-config-prod.yaml
cd ..

services=("uf_aim" "uf_border" "uf_maha" "uf_order" "uf_public" "uf_user")

for service in ${services[@]}; do
	echo ===== redeploying $service =====
	cd $service
	kubectl delete -f deploy-latest.yaml
	kubectl apply -f deploy-latest.yaml
	cd ..
done
