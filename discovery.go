package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/storageos/discovery/cluster"
	"github.com/storageos/discovery/handlers"
	"github.com/storageos/discovery/store/boltdb"
	"github.com/storageos/discovery/util/codecs"
)

// DefaultPort - default port to run
const DefaultPort = 8081

// EnvPort - port to run application on
const EnvPort = "PORT"

// EnvDatabasePath - database path
const EnvDatabasePath = "DATABASE_PATH"

func main() {
	port := DefaultPort
	if os.Getenv(EnvPort) != "" {
		p, err := strconv.Atoi(os.Getenv(EnvPort))
		if err != nil {
			log.Fatalf("invalid port: %s", err)
		}
		port = p
	}

	path := "discovery.db"
	if os.Getenv(EnvDatabasePath) != "" {		
		path = filepath.Join(os.Getenv(EnvDatabasePath), "discovery.db")
	}

	db, err := boltdb.New(path)
	if err != nil {
		log.Fatalf("failed to init database: %s", err)
	}

	clusterManager := cluster.New(db, codecs.DefaultSerializer())

	srv := handlers.NewServer(port, clusterManager)
	log.Fatal(srv.Start())
}
