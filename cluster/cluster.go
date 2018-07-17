package cluster

import (
	"errors"
	"net/url"
	"sync"
	"time"

	"github.com/storageos/discovery/store"
	"github.com/storageos/discovery/types"
	"github.com/storageos/discovery/util/codecs"
	"github.com/storageos/discovery/util/uuid"
)

// node registration errors
var (
	ErrAddressMissing     = errors.New("node address missing")
	ErrInvalidAddress     = errors.New("invalid node address")
	ErrNameMissing        = errors.New("node name missing")
	ErrNodeNamePresent    = errors.New("node name already present")
	ErrNodeAddressPresent = errors.New("node address already present")
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

// DefaultManager - default cluster manager
type DefaultManager struct {
	mu         *sync.Mutex
	store      store.Store
	serializer codecs.Serializer
}

// New - create new cluster manager
func New(store store.Store, serializer codecs.Serializer) *DefaultManager {
	return &DefaultManager{
		mu:         &sync.Mutex{},
		store:      store,
		serializer: serializer,
	}
}

// Create - create new cluster
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

// Get - get cluster by ID
func (m *DefaultManager) Get(ref string) (*types.Cluster, error) {
	kvp, err := m.store.Get(ref)
	if err != nil {
		return nil, err
	}
	var cluster types.Cluster
	err = m.serializer.Decode(kvp.Value, &cluster)
	return &cluster, err
}

func nodeValid(node *types.Node) error {
	if node.AdvertiseAddress == "" {
		return ErrAddressMissing
	}

	if node.Name == "" {
		return ErrNameMissing
	}

	_, err := url.Parse(node.AdvertiseAddress)
	if err != nil {
		return ErrInvalidAddress
	}
	return nil
}

// RegisterNode - register new node to the cluster
func (m *DefaultManager) RegisterNode(clusterID string, node *types.Node) (updated *types.Cluster, err error) {

	err = nodeValid(node)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	cluster, err := m.Get(clusterID)
	if err != nil {
		return nil, err
	}

	// looking for duplicates
	for _, n := range cluster.Nodes {
		if n.Name == node.Name && n.AdvertiseAddress == node.AdvertiseAddress && n.ID == node.ID {
			// node already registered, nothing to do
			return cluster, nil
		}

		if n.Name == node.Name {
			return nil, ErrNodeNamePresent
		}

		if n.AdvertiseAddress == node.AdvertiseAddress {
			return nil, ErrNodeAddressPresent
		}
	}

	node.CreatedAt = time.Now()
	node.UpdatedAt = time.Now()

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

// Update - update cluster
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

// Delete - delete cluster by ID
func (m *DefaultManager) Delete(id string) error {
	return m.store.Delete(id)
}
