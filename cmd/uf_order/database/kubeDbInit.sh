#!/bin/sh
export PGPASSWORD=postgres
export dbName=uf_order

name="order-pg"
pod=$(kubectl get pods | awk '{print $1}' | grep -e $name)
echo name is $name

files=$(ls initdb)
for file in ${files}; do
	echo "--->executing $file"
	cat "initdb/${file}" | \
	kubectl exec -it $pod \
		-- psql -U postgres \
		-d $dbName
done
