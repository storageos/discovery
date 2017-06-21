package http

import (
	"net/http"
	"os"

	gorillaHandlers "github.com/gorilla/handlers"

	"github.com/storageos/discovery/handlers"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/gorilla/mux"
)

func Setup(etcdHost, discHost string) {
	handlers.Setup(etcdHost, discHost)
	r := mux.NewRouter()

	r.HandleFunc("/", handlers.HomeHandler)
	r.HandleFunc("/new", handlers.NewTokenHandler)
	r.HandleFunc("/health", handlers.HealthHandler)
	r.HandleFunc("/robots.txt", handlers.RobotsHandler)

	// Only allow exact tokens with GETs and PUTs
	r.HandleFunc("/{token:[a-f0-9]{32}}", handlers.TokenHandler).
		Methods("GET", "PUT")
	r.HandleFunc("/{token:[a-f0-9]{32}}/", handlers.TokenHandler).
		Methods("GET", "PUT")
	r.HandleFunc("/{token:[a-f0-9]{32}}/{machine}", handlers.TokenHandler).
		Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/{token:[a-f0-9]{32}}/_config/size", handlers.TokenHandler).
		Methods("GET")

	logH := gorillaHandlers.LoggingHandler(os.Stdout, r)

	http.Handle("/", logH)
	http.Handle("/metrics", prometheus.Handler())
}
