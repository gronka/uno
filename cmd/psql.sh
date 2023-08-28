#!/bin/sh
#name=$1
#check=$(./check-uf.sh $name)
#if [ "$check" != "$name" ]; then
	#echo $check
	#exit 1
#fi

export PGPASSWORD=postgres
#export dbName=$1
#export dbName=uf_order
#container=$($name//[_]/-)
#${orig//[xyz]/_}

pod=$(kubectl get pods | awk '{print $1}' | grep -e 'order-pg')
#echo container is $container

kubectl exec -it $pod \
	-- psql -U postgres 
	#-d $dbName
