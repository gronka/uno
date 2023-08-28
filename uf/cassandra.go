package uf

import (
	"time"

	"github.com/gocql/gocql"
)

var csb *gocql.Session // binary partitioner
var csm *gocql.Session // murmur3 partitioner

func connectCassandra(initial bool) {
	if !initial {
		csm.Close()
	}

	casshost := "127.0.0.1"
	cluster := gocql.NewCluster(casshost)
	cluster.Keyspace = "fridayy"
	cluster.Consistency = gocql.One
	cluster.ProtoVersion = 3
	//TODO: set up a separate cassandra that runs murmur3

	var err error
	csb, err = cluster.CreateSession()
	csm, err = cluster.CreateSession()
	for err != nil {
		time.Sleep(1 * time.Second)
		Error("CQL session could not be established. Retrying in 1 second.")
		// TODO: email  warning every so many loops here
		csb, err = cluster.CreateSession()
		csm, err = cluster.CreateSession()
	}
	Info("CQL session established!")
}
