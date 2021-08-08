package cassandra

import (
	"github.com/gocql/gocql"
)

func GetSession() *gocql.Session {
	var session *gocql.Session
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "bookstore_oauth"
	cluster.Consistency = gocql.Quorum

	var err error
	if session, err = cluster.CreateSession(); err != nil {
		panic(err)
	}
	return session
}
