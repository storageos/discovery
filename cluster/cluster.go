package cluster

import (
	"fmt"
	"sync"
	"time"

	"github.com/storageos/discovery/store"
	"github.com/storageos/discovery/types"
	"github.com/storageos/discovery/util/codecs"
	"github.com/storageos/discovery/util/uuid"
)

// Manager - cluster manager
type Manager interface {
	// create new cluster
	Create(opts types.ClusterCreateOps) (*types.Cluster, error)

	// get cluster by ID
	Get(ref string) (*types.Cluster, error)
	// register node
	RegisterNode(clusterID string, node *types.Node) (updated *types.Cluster, err error)
	// update cluster details
	Update(cluster *types.Cluster) error
	// delete cluster
	Delete(id string) error
}

type DefaultManager struct {
	mu         *sync.Mutex
	store      store.Store
	serializer codecs.Serializer
}

func New(store store.Store, serializer codecs.Serializer) *DefaultManager {
	return &DefaultManager{
		mu:         &sync.Mutex{},
		store:      store,
		serializer: serializer,
	}
}

func (m *DefaultManager) Create(opts types.ClusterCreateOps) (*types.Cluster, error) {
	cluster := types.Cluster{
		ID:        uuid.Generate(),
		AccountID: opts.AccountID,
		Name:      opts.Name,
		Size:      opts.Size,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if cluster.Size == 0 {
		cluster.Size = 3
	}

	bts, err := m.serializer.Encode(&cluster)
	if err != nil {
		return nil, err
	}

	_, err = m.store.Create(cluster.ID, bts, opts.TTL)
	if err != nil {
		return nil, err
	}

	return &cluster, nil
}

func (m *DefaultManager) Get(ref string) (*types.Cluster, error) {
	kvp, err := m.store.Get(ref)
	if err != nil {
		return nil, err
	}
	var cluster types.Cluster
	err = m.serializer.Decode(kvp.Value, &cluster)
	return &cluster, err
}

func (m *DefaultManager) RegisterNode(clusterID string, node *types.Node) (updated *types.Cluster, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cluster, err := m.Get(clusterID)
	if err != nil {
		return nil, err
	}

	// looking for duplicates
	for _, n := range cluster.Nodes {
		if n.Name == node.Name && n.ID != node.ID {
			return nil, fmt.Errorf("node with name %s already exists in cluster %s", node.Name, clusterID)
		}

		if n.AdvertiseIP == node.AdvertiseIP && n.ID != node.ID {
			return nil, fmt.Errorf("node with advertise IP %s already exists in cluster %s", node.AdvertiseIP, clusterID)
		}
	}

	cluster.Nodes = append(cluster.Nodes, node)

	bts, err := m.serializer.Encode(cluster)
	if err != nil {
		return nil, err
	}

	_, err = m.store.Put(cluster.ID, bts, 0)
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

func (m *DefaultManager) Update(cluster *types.Cluster) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	bts, err := m.serializer.Encode(cluster)
	if err != nil {
		return err
	}

	_, err = m.store.Put(cluster.ID, bts, 0)

	return err
}

func (m *DefaultManager) Delete(id string) error {
	return m.store.Delete(id)
}
