package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/storageos/discovery/cluster"
	"github.com/storageos/discovery/store"
	"github.com/storageos/discovery/types"
)

func TestRegisterNodeHandler(t *testing.T) {
	srv := setupTestServer(t)
	defer teardownTestServer(t, srv)
	// Create cluster as prerequisite
	c, err := srv.server.clusterManager.Create(types.ClusterCreateOps{Size: 3})
	if err != nil {
		t.Error(err)
	}
	// Pre-place cluster for name/address conflict
	srv.server.clusterManager.RegisterNode(
		c.ID,
		&types.Node{ID: "4", Name: "node4", AdvertiseAddress: "192.168.0.4", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	)
	// Parallel test cases
	testcases := []struct {
		name      string
		clusterID string
		node      types.Node
		code      int
		errMsg    string
	}{
		{
			name:      "ok registration - node 1",
			clusterID: c.ID,
			node:      types.Node{ID: "1", Name: "node1", AdvertiseAddress: "192.168.0.1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			code:      http.StatusOK,
			errMsg:    "",
		},
		{
			name:      "ok registration - node 2",
			clusterID: c.ID,
			node:      types.Node{ID: "2", Name: "node2", AdvertiseAddress: "192.168.0.2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			code:      http.StatusOK,
			errMsg:    "",
		},
		{
			name:      "ok registration - node 3",
			clusterID: c.ID,
			node:      types.Node{ID: "3", Name: "node3", AdvertiseAddress: "192.168.0.3", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			code:      http.StatusOK,
			errMsg:    "",
		},
		{
			name:      "name already present - node 4",
			clusterID: c.ID,
			node:      types.Node{ID: "5", Name: "node4", AdvertiseAddress: "192.168.0.5", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			code:      http.StatusConflict,
			errMsg: fmt.Sprintf(
				"%s: name node4 exists in cluster %s\n",
				cluster.ErrNodeNamePresent,
				c.ID,
			),
		},
		{
			name:      "address already present - 192.168.0.4",
			clusterID: c.ID,
			node:      types.Node{ID: "6", Name: "node6", AdvertiseAddress: "192.168.0.4", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			code:      http.StatusConflict,
			errMsg: fmt.Sprintf(
				"%s: address 192.168.0.4 exists in cluster %s\n",
				cluster.ErrNodeAddressPresent,
				c.ID,
			),
		},
		{
			name:      "address missing",
			clusterID: c.ID,
			node:      types.Node{ID: "5", Name: "node5", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			code:      http.StatusBadRequest,
			errMsg:    fmt.Sprintf("%s\n", cluster.ErrAddressMissing),
		},
		{
			name:      "name missing",
			clusterID: c.ID,
			node:      types.Node{ID: "5", AdvertiseAddress: "192.168.0.5", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			code:      http.StatusBadRequest,
			errMsg:    fmt.Sprintf("%s\n", cluster.ErrNameMissing),
		},
		{
			name:      "invalid address",
			clusterID: c.ID,
			node:      types.Node{ID: "5", Name: "node5", AdvertiseAddress: "192.168.0.5:5555", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			code:      http.StatusBadRequest,
			errMsg:    fmt.Sprintf("%s\n", cluster.ErrInvalidAddress),
		},
		{
			name:      "non-existent cluster",
			clusterID: "123", // non-existent cluster ID
			node:      types.Node{ID: "5", Name: "node5", AdvertiseAddress: "192.168.0.5", CreatedAt: time.Now(), UpdatedAt: time.Now()},
			code:      http.StatusNotFound,
			errMsg:    fmt.Sprintf("%s\n", store.ErrNotFound),
		},
	}

	// Register nodes correctly
	registerNode := func(t *testing.T, clusterID string, n types.Node) *httptest.ResponseRecorder {
		reqBody, err := json.Marshal(n)
		if err != nil {
			t.Fatalf("failed to marshal node: %v", err)
		}
		req, err := http.NewRequest(http.MethodPut,
			fmt.Sprintf("/clusters/%s", clusterID),
			bytes.NewBuffer(reqBody))
		if err != nil {
			t.Fatalf("failed to create request: %v", err)
		}

		rec := httptest.NewRecorder()

		srv.server.mux.ServeHTTP(rec, req)

		return rec
	}

	t.Run("register requests", func(t *testing.T) {
		for _, tc := range testcases {
			want := tc // capture
			t.Run(want.name, func(t *testing.T) {
				t.Parallel()
				resp := registerNode(t, want.clusterID, want.node)
				if resp.Code != want.code {
					t.Errorf("\ngot code %d\n wanted code %d", resp.Code, want.code)
				}
				if want.errMsg != "" {
					if resp.Body.String() != want.errMsg {
						t.Errorf("\ngot err:\n %s \n wanted err:\n %s", resp.Body.String(), want.errMsg)
					}
				}
			})
		}
	})
}
