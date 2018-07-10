package handlers

import (
	"io/ioutil"
	"net"
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
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal("could not get free port to run test server on")
	}

	port := listener.Addr().(*net.TCPAddr).Port

	if err := listener.Close(); err != nil {
		t.Fatalf("couldn't close listener: %v", err)
	}

	file, err := ioutil.TempFile("", "discovery-test-"+uuid.Generate())
	if err != nil {
		t.Fatal(err)
	}

	store, err := boltdb.New(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	cm := cluster.New(store, codecs.DefaultSerializer())

	server := NewServer(port, cm)

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
