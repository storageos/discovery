package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/storageos/discovery/types"
)

// DefaultEndpoint - default endpoint address
const DefaultEndpoint = "http://localhost:8081"

// Client - generic client interface
type Client interface {
	ClusterGet(ref string) (*types.Cluster, error)
	ClusterRegisterNode(clusterID, nodeID, name, advertiseIP string) (*types.Cluster, error)
}

type DefaultClient struct {
	endpoint string
	client   *http.Client
}

func New(options ...Option) *DefaultClient {
	client := DefaultClient{
		endpoint: DefaultEndpoint,
		client:   &http.Client{},
	}

	for _, opt := range options {
		opt.Configure(&client)
	}

	return &client
}

func (c *DefaultClient) ClusterGet(ref string) (*types.Cluster, error) {
	req, err := http.NewRequest("GET", c.endpoint+"/cluster/"+ref, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var cluster types.Cluster
	if err := json.NewDecoder(resp.Body).Decode(&cluster); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response from discovery service: %s", err)
	}

	return &cluster, nil
}

func WithEndpoint(endpoint string) Option {
	return OptionFn(func(c *DefaultClient) error {
		c.endpoint = endpoint
		return nil
	})
}

// Option is used to pass optional arguments to
// the DefaultRecorder constructor
type Option interface {
	Configure(*DefaultClient) error
}

// OptionFn is a type of Option that is represented
// by a single function that gets called for Configure()
type OptionFn func(*DefaultClient) error

// Configure - configures specific variable
func (o OptionFn) Configure(dc *DefaultClient) error {
	return o(dc)
}
