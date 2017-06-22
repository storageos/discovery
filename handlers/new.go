package handlers

import (
	"log"
	"net/http"
	"strconv"

	"encoding/json"

	"github.com/storageos/discovery/types"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/storageos/discovery/handlers/httperror"
)

var newCounter *prometheus.CounterVec

// var cfg *client.Config
var discHost string

func init() {
	newCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "endpoint_new_requests_total",
			Help: "How many /new requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method"},
	)
	prometheus.MustRegister(newCounter)
}

func (s *Server) newClusterHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	size := 3
	sz := r.FormValue("size")
	if sz != "" {
		size, err = strconv.Atoi(sz)
		if err != nil {
			httperror.Error(w, r, err.Error(), http.StatusBadRequest, newCounter)
			return
		}
	}
	cluster, err := s.clusterManager.Create(types.ClusterCreateOps{Size: size})
	if err != nil {
		httperror.Error(w, r, err.Error(), http.StatusInternalServerError, newCounter)
		return
	}

	log.Println("New cluster created", cluster.ID)

	bts, err := json.Marshal(cluster)
	if err != nil {
		httperror.Error(w, r, err.Error(), http.StatusInternalServerError, newCounter)
		return
	}

	w.Write(bts)
	newCounter.WithLabelValues("200", r.Method).Add(1)
}
