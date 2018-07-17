package handlers

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/storageos/discovery/cluster"
	"github.com/storageos/discovery/store/boltdb"
	"github.com/storageos/discovery/util/codecs"
	"github.com/storageos/discovery/util/uuid"
)

type TestServer struct {
	path   string
	server *Server
	store  *boltdb.Store
}

func setupTestServer(t *testing.T) *TestServer {
	file, err := ioutil.TempFile("", "discovery-test-"+uuid.Generate())
	if err != nil {
		t.Fatal(err)
	}

	store, err := boltdb.New(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	cm := cluster.New(store, codecs.DefaultSerializer())

	server := NewServer(1, cm)

	return &TestServer{
		path:   file.Name(),
		server: server,
		store:  store,
	}
}

func teardownTestServer(t *testing.T, s *TestServer) {
	s.store.Close()
	os.Remove(s.path)
}
