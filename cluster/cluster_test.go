package cluster

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/storageos/discovery/store/boltdb"
	"github.com/storageos/discovery/types"
	"github.com/storageos/discovery/util/codecs"
)

func TestClusterCreate(t *testing.T) {
	dir, err := ioutil.TempDir("", "testcreatecluster")
	if err != nil {
		t.Fatalf("failed to get temp dir: %s", err)
	}

	db, err := boltdb.New(dir + "testdb")
	if err != nil {
		t.Fatalf("failed to create db: %s", err)
	}

	cm := New(db, codecs.DefaultSerializer())

	cluster, err := cm.Create(types.ClusterCreateOps{AccountID: "123"})
	if err != nil {
		t.Errorf("failed to create cluster: %s", err)
	}

	if cluster.ID == "" {
		t.Errorf("expected cluster ID to be populated")
	}

	if cluster.Size != 3 {
		t.Errorf("unexpected cluster size: %d", cluster.Size)
	}

}

func TestClusterRegisterNode(t *testing.T) {
	dir, err := ioutil.TempDir("", "testcreatecluster")
	if err != nil {
		t.Fatalf("failed to get temp dir: %s", err)
	}

	db, err := boltdb.New(dir + "testdb")
	if err != nil {
		t.Fatalf("failed to create db: %s", err)
	}

	cm := New(db, codecs.DefaultSerializer())

	cluster, err := cm.Create(types.ClusterCreateOps{AccountID: "123"})
	if err != nil {
		t.Errorf("failed to create cluster: %s", err)
	}

	node := &types.Node{
		ID:          "controller-uuid-1",
		Name:        "node-1",
		AdvertiseIP: "10.0.1.4",
	}

	updatedCluster, err := cm.RegisterNode(cluster.ID, node)
	if err != nil {
		t.Errorf("failed to update cluster: %s", err)
	}

	if updatedCluster.Nodes[0].Name != node.Name {
		t.Errorf("expected node name %s but got %s", node.Name, updatedCluster.Nodes[0].Name)
	}

}

func TestClusterRegisterNodes(t *testing.T) {
	dir, err := ioutil.TempDir("", "testcreatecluster")
	if err != nil {
		t.Fatalf("failed to get temp dir: %s", err)
	}

	db, err := boltdb.New(dir + "testdb")
	if err != nil {
		t.Fatalf("failed to create db: %s", err)
	}

	cm := New(db, codecs.DefaultSerializer())

	cluster, err := cm.Create(types.ClusterCreateOps{AccountID: "123"})
	if err != nil {
		t.Errorf("failed to create cluster: %s", err)
	}

	for i := 0; i < 10; i++ {
		node := &types.Node{
			ID:          "controller-uuid-1",
			Name:        fmt.Sprintf("node-%d", i),
			AdvertiseIP: fmt.Sprintf("10.0.1.%d", i),
		}

		_, err := cm.RegisterNode(cluster.ID, node)
		if err != nil {
			t.Errorf("failed to update cluster: %s", err)
		}
	}

	updated, err := cm.Get(cluster.ID)
	if err != nil {
		t.Errorf("failed to get cluster: %s", err)
	}

	if len(updated.Nodes) != 10 {
		t.Errorf("unexpected number of nodes in cluster: %d", len(updated.Nodes))
	}

}
