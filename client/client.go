package client

import (
	"github.com/storageos/discovery/types"
)

// Client - generic client interface
type Client interface {
	ClusterStatus(ref string) (*types.Cluster, error)
	RegisterNode(clusterID, name, advertiseIP string) (*types.Cluster, error)
}
