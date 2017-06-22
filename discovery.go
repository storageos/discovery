package main

import (
	"log"

	"github.com/storageos/discovery/cluster"
	"github.com/storageos/discovery/handlers"
	"github.com/storageos/discovery/store/boltdb"
	"github.com/storageos/discovery/util/codecs"
)

func main() {

	db, err := boltdb.New("discovery.db")
	if err != nil {
		log.Fatalf("failed to init database: %s", err)
	}

	clusterManager := cluster.New(db, codecs.DefaultSerializer())

	srv := handlers.NewServer(8081, clusterManager)
	log.Fatal(srv.Start())
}
