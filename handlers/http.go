package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/storageos/discovery/cluster"

	"github.com/gorilla/mux"
)

// Server - main discovery service server container
type Server struct {
	clusterManager cluster.Manager
	port           int
	server         *http.Server
	mux            *mux.Router
}

// NewServer - new discovery http server
func NewServer(port int, cm cluster.Manager) *Server {
	srv := &Server{
		clusterManager: cm,
		port:           port,
	}

	srv.registerHandlers()

	return srv
}

// Start - configures and starts HTTP server
func (s *Server) Start() error {

	s.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", s.port),
		Handler:           s.mux,
		IdleTimeout:       time.Second * 120,
		ReadTimeout:       time.Second * 30,
		ReadHeaderTimeout: time.Second * 30,
		WriteTimeout:      time.Second * 25,
	}
	log.Printf("server starting on port %d", s.port)

	s.server.Handler = gorillaHandlers.LoggingHandler(os.Stdout, s.mux)

	return s.server.ListenAndServe()
}

// Stop - stops HTTP server
func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(5))
	defer cancel()
	s.server.Shutdown(ctx)
}

func getParam(param string, req *http.Request) string {
	return mux.Vars(req)[param]
}

const (
	paramCluster string = "ref"
	paramNode    string = "node"
)

func (s *Server) registerHandlers() {
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/version", versionHandler)
	r.HandleFunc("/health", s.healthHandler)
	r.HandleFunc("/robots.txt", robotsHandler)

	r.HandleFunc("/clusters", s.newClusterHandler).Methods("POST")
	r.HandleFunc("/clusters/{ref}", s.clusterHandler).Methods("GET")
	r.HandleFunc("/clusters/{ref}", s.registerNodeHandler).Methods("PUT")
	r.HandleFunc("/clusters/{ref}", s.deleteClusterHandler).Methods("DELETE")

	r.Handle("/metrics", promhttp.Handler())

	logH := gorillaHandlers.LoggingHandler(os.Stdout, r)

	http.Handle("/", logH)

	s.mux = r
}
