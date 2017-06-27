package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/storageos/discovery/handlers/httperror"

	"github.com/storageos/discovery/version"
)

var versionCounter *prometheus.CounterVec

func init() {
	versionCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "endpoint_version_requests_total",
			Help: "How many /version requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method"},
	)
	prometheus.MustRegister(versionCounter)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	v := version.GetVersion()

	bts, err := json.Marshal(&v)
	if err != nil {
		httperror.Error(w, r, "health failed to get version", 500, versionCounter)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bts)
	versionCounter.WithLabelValues("200", r.Method).Add(1)
}
