package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/storageos/discovery/handlers/httperror"
	"github.com/storageos/discovery/types"
)

var healthCounter *prometheus.CounterVec

func init() {
	healthCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "endpoint_health_requests_total",
			Help: "How many /health requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method"},
	)
	prometheus.MustRegister(healthCounter)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	cluster, err := s.clusterManager.Create(types.ClusterCreateOps{})

	if err != nil {
		log.Printf("health failed to create cluster %v", err)
		httperror.Error(w, r, "health failed to create cluster", 400, healthCounter)
		return
	}

	err = s.clusterManager.Delete(cluster.ID)
	if err != nil {
		log.Printf("health failed to delete cluster %v", err)
		httperror.Error(w, r, "health failed to delete cluster", 400, healthCounter)
		return
	}

	fmt.Fprintf(w, "OK")
	healthCounter.WithLabelValues("200", r.Method).Add(1)
}
