package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/storageos/discovery/types"
)

// DefaultEndpoint - default endpoint address
const DefaultEndpoint = "http://discovery.storageos.cloud"

// Client - generic client interface
type Client interface {
	ClusterGet(ref string) (*types.Cluster, error)
	ClusterRegisterNode(clusterID, nodeID, name, advertiseIP string) (*types.Cluster, error)
}

// DefaultClient - default discovery client
type DefaultClient struct {
	endpoint string
	client   *http.Client
}

// New - create new discovery client
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

// ClusterGet - get specified cluster by ID
func (c *DefaultClient) ClusterGet(ref string) (*types.Cluster, error) {
	req, err := http.NewRequest("GET", c.endpoint+"/clusters/"+ref, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var cluster types.Cluster
	if err := json.NewDecoder(resp.Body).Decode(&cluster); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response from discovery service: %s", err)
	}

	return &cluster, nil
}

// ClusterCreate - create cluster
func (c *DefaultClient) ClusterCreate(opts types.ClusterCreateOps) (*types.Cluster, error) {

	path := c.endpoint + "/clusters"
	vals := url.Values{}
	vals.Set("size", fmt.Sprintf("%d", opts.Size))
	vals.Set("name", fmt.Sprintf("%s", opts.Name))
	path = fmt.Sprintf("%s?%s", path, vals.Encode())

	req, err := http.NewRequest("POST", path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respMsg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("unexpected status code: %d, response body unavailable", resp.StatusCode)
		}
		return nil, fmt.Errorf("unexpected status code: %d (%s)", resp.StatusCode, string(respMsg))
	}

	var cluster types.Cluster
	if err := json.NewDecoder(resp.Body).Decode(&cluster); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response from discovery service: %s", err)
	}

	return &cluster, nil
}

// ClusterRegisterNode - register node to cluster
func (c *DefaultClient) ClusterRegisterNode(clusterID, nodeID, name, advertiseIP string) (*types.Cluster, error) {
	node := types.Node{
		ID:          nodeID,
		Name:        name,
		AdvertiseIP: advertiseIP,
	}
	reqBody, err := json.Marshal(&node)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", c.endpoint+"/clusters/"+clusterID, ioutil.NopCloser(bytes.NewBuffer(reqBody)))
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respMsg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("unexpected status code: %d, response body unavailable", resp.StatusCode)
		}
		return nil, fmt.Errorf("unexpected status code: %d (%s)", resp.StatusCode, string(respMsg))
	}

	var cluster types.Cluster
	if err := json.NewDecoder(resp.Body).Decode(&cluster); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response from discovery service: %s", err)
	}

	return &cluster, nil
}

// WithEndpoint - override default endpoint
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
