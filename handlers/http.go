package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	gorillaHandlers "github.com/gorilla/handlers"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/storageos/discovery/cluster"

	"github.com/gorilla/mux"
)

type Server struct {
	clusterManager cluster.Manager
	port           int
	server         *http.Server
	mux            *mux.Router
}

func NewServer(port int, cm cluster.Manager) *Server {
	srv := &Server{
		clusterManager: cm,
		port:           port,
	}

	srv.registerHandlers()

	return srv
}

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
	return s.server.ListenAndServe()
}

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

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/health", s.healthHandler)
	r.HandleFunc("/robots.txt", robotsHandler)

	r.HandleFunc("/clusters", s.newClusterHandler).Methods("POST")
	r.HandleFunc("/clusters/{ref}", s.clusterHandler).Methods("GET")
	r.HandleFunc("/clusters/{ref}", s.registerNodeHandler).Methods("PUT")
	r.HandleFunc("/clusters/{ref}", s.deleteClusterHandler).Methods("DELETE")

	logH := gorillaHandlers.LoggingHandler(os.Stdout, r)

	http.Handle("/", logH)
	http.Handle("/metrics", prometheus.Handler())

	s.mux = r
}
