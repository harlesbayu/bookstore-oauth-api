package cassandra

import (
	"github.com/gocql/gocql"
)

var (
	cluster *gocql.ClusterConfig
)

func GetSession() (*gocql.Session, error) {
	cluster = gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "bookstore_oauth"
	cluster.Consistency = gocql.Quorum
	return cluster.CreateSession()
}
