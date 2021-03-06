package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/storageos/discovery/cluster"
	"github.com/storageos/discovery/handlers/httperror"
	"github.com/storageos/discovery/store"
	"github.com/storageos/discovery/types"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	tokenCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "endpoint_cluster_requests_total",
			Help: "How many /cluster requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method"},
	)
	prometheus.MustRegister(tokenCounter)
}

var (
	tokenCounter *prometheus.CounterVec
)

func (s *Server) clusterHandler(w http.ResponseWriter, r *http.Request) {

	cluster, err := s.clusterManager.Get(getParam(paramCluster, r))
	if err != nil {
		if err == store.ErrNotFound {
			httperror.Error(w, r, err.Error(), http.StatusNotFound, newCounter)
			return
		}
		httperror.Error(w, r, err.Error(), http.StatusInternalServerError, newCounter)
		return
	}

	bts, err := json.Marshal(cluster)
	if err != nil {
		httperror.Error(w, r, err.Error(), http.StatusInternalServerError, newCounter)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bts)
	tokenCounter.WithLabelValues(strconv.Itoa(http.StatusOK), r.Method).Add(1)
}

func (s *Server) registerNodeHandler(w http.ResponseWriter, r *http.Request) {
	clusterID := getParam(paramCluster, r)

	var node types.Node
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		httperror.Error(w, r, err.Error(), http.StatusInternalServerError, newCounter)
		return
	}

	updated, err := s.clusterManager.RegisterNode(clusterID, &node)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			httperror.Error(w, r, err.Error(), http.StatusNotFound, newCounter)
			return
		case cluster.ErrAddressMissing, cluster.ErrInvalidAddress, cluster.ErrNameMissing:
			httperror.Error(w, r, err.Error(), http.StatusBadRequest, newCounter)
			return
		case cluster.ErrNodeNamePresent:
			httperror.Error(w, r, err.Error()+fmt.Sprintf(": name %s exists in cluster %s", node.Name, clusterID), http.StatusUnprocessableEntity, newCounter)
			return
		case cluster.ErrNodeAddressPresent:
			httperror.Error(w, r, err.Error()+fmt.Sprintf(": address %s exists in cluster %s", node.AdvertiseAddress, clusterID), http.StatusUnprocessableEntity, newCounter)
			return
		}

		httperror.Error(w, r, err.Error(), http.StatusInternalServerError, newCounter)
		return
	}

	bts, err := json.Marshal(updated)
	if err != nil {
		httperror.Error(w, r, err.Error(), http.StatusInternalServerError, newCounter)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bts)
	newCounter.WithLabelValues("200", r.Method).Add(1)

}

func (s *Server) deleteClusterHandler(w http.ResponseWriter, r *http.Request) {
	err := s.clusterManager.Delete(getParam(paramCluster, r))
	if err != nil {
		httperror.Error(w, r, err.Error(), http.StatusInternalServerError, newCounter)
		return
	}

	newCounter.WithLabelValues("200", r.Method).Add(1)
}
