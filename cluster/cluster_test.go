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
		ID:               "controller-uuid-1",
		Name:             "node-1",
		AdvertiseAddress: "10.0.1.4",
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
			ID:               "controller-uuid-1",
			Name:             fmt.Sprintf("node-%d", i),
			AdvertiseAddress: fmt.Sprintf("10.0.1.%d", i),
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

func Test_nodeValid(t *testing.T) {
	type args struct {
		node *types.Node
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid node",
			args: args{
				node: &types.Node{
					AdvertiseAddress: "http://localhost:2333",
					Name:             "node-1",
				},
			},
			wantErr: false,
		},
		{
			name: "missing node name",
			args: args{
				node: &types.Node{
					AdvertiseAddress: "http://localhost:2333"},
			},
			wantErr: true,
		},
		{
			name: "missing node address",
			args: args{
				node: &types.Node{
					Name: "node-1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := nodeValid(tt.args.node); (err != nil) != tt.wantErr {
				t.Errorf("nodeValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
