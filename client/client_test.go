package client

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/storageos/discovery/cluster"
	"github.com/storageos/discovery/handlers"
	"github.com/storageos/discovery/store/boltdb"
	"github.com/storageos/discovery/types"
	"github.com/storageos/discovery/util/codecs"

	"testing"
)

const testServerPort = 4551
const testServerEndpoint = "http://127.0.0.1:4551"

func TestMain(m *testing.M) {

	srv := newTestingServer()
	go srv.Start()

	retCode := m.Run()

	srv.Stop()

	os.Exit(retCode)
}

func newTestingServer() *handlers.Server {
	dir, err := ioutil.TempDir("", "testcreatecluster")
	if err != nil {
		log.Fatalf("failed to get temp dir: %s", err)
	}

	db, err := boltdb.New(dir + "testdb")
	if err != nil {
		log.Fatalf("failed to create db: %s", err)
	}

	cm := cluster.New(db, codecs.DefaultSerializer())

	srv := handlers.NewServer(testServerPort, cm)
	return srv
}

func TestClientRegister(t *testing.T) {
	client := New(WithEndpoint(testServerEndpoint))

	nodeAdvertiseIP := "1.1.1.1"
	nodeName := "node-1"
	nodeID := "client-node-uuid"

	newCluster, err := client.ClusterCreate(types.ClusterCreateOps{Name: "new-1", Size: 3})
	if err != nil {
		t.Fatalf("failed to create cluster: %s", err)
	}

	cluster, err := client.ClusterRegisterNode(newCluster.ID, nodeID, nodeName, nodeAdvertiseIP)
	if err != nil {
		t.Fatalf("failed to register node: %s", err)
	}

	if len(cluster.Nodes) == 0 {
		t.Fatalf("cannot find node in the cluster nodes")
	}

	if cluster.Nodes[0].AdvertiseAddress != nodeAdvertiseIP {
		t.Errorf("unexpected advertise address: %s", cluster.Nodes[0].AdvertiseAddress)
	}
}
