#!/bin/bash

services=("uf_aim" "uf_maha" "uf_order" "uf_user")

for service in ${services[@]}; do
	echo ===== reiniting database for $service =====
	cd $service/database
	./dbInit.sh
	cd -
done
