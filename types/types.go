package types

import (
	"bytes"
	"fmt"
	"time"
)

// ClusterCreateOps - optional fields when creating cluster
type ClusterCreateOps struct {
	AccountID string
	// optional value when to expire cluster
	TTL  int64
	Name string
	Size int
}

type Cluster struct {
	// cluster ID used for joining or getting cluster status
	ID string `json:"id,omitempty"`

	// cluster size, defaults to 3
	Size int `json:"size,omitempty"`

	Name string `json:"name,omitempty"`

	// optional account ID
	AccountID string `json:"accountID,omitempty"`

	// nodes participating in cluster
	Nodes []*Node `json:"nodes,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type Node struct {
	ID               string `json:"id,omitempty"` // node/controller UUID
	Name             string `json:"name,omitempty"`
	AdvertiseAddress string `json:"advertiseAddress,omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

// VersionInfo describes the server's version and runtime info.
type VersionInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	APIVersion   string `json:"apiVersion"`
	BuildDate    string `json:"buildDate"`
	Experimental bool   `json:"experimental"`
}

// String returns string representation of version info.
func (v VersionInfo) String() string {
	var versionString bytes.Buffer
	fmt.Fprintf(&versionString, "Discovery %s", v.Version)
	return versionString.String()
}
