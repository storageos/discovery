package handlers

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/storageos/discovery.etcd.io/handlers/httperror"
	"github.com/storageos/discovery.etcd.io/pkg/lockstring"
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	tokenCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "endpoint_token_requests_total",
			Help: "How many /token requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method"},
	)
	prometheus.MustRegister(tokenCounter)
}

var (
	currentLeader lockstring.LockString
	tokenCounter  *prometheus.CounterVec
)

func proxyRequest(r *http.Request) (*http.Response, error) {
	body, _ := ioutil.ReadAll(r.Body)

	for i := 0; i <= 10; i++ {
		u := url.URL{
			Scheme:   "http",
			Host:     currentLeader.String(),
			Path:     path.Join("v2", "keys", "_etcd", "registry", r.URL.Path),
			RawQuery: r.URL.RawQuery,
		}

		buf := bytes.NewBuffer(body)
		outreq, err := http.NewRequest(r.Method, u.String(), buf)
		if err != nil {
			return nil, err
		}

		copyHeader(outreq.Header, r.Header)

		client := http.Client{}
		resp, err := client.Do(outreq)
		if err != nil {
			return nil, err
		}

		// Try again on the next host
		if resp.StatusCode == 307 &&
			(r.Method == "PUT" || r.Method == "DELETE") {
			u, err := resp.Location()
			if err != nil {
				return nil, err
			}
			currentLeader.Set(u.Host)
			continue
		}

		return resp, nil
	}

	return nil, errors.New("All attempts at proxying to etcd failed")
}

// copyHeader copies all of the headers from dst to src.
func copyHeader(dst, src http.Header) {
	for k, v := range src {
		for _, q := range v {
			dst.Add(k, q)
		}
	}
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := proxyRequest(r)
	if err != nil {
		log.Printf("Error making request: %v", err)
		httperror.Error(w, r, "", 500, tokenCounter)
	}

	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
	tokenCounter.WithLabelValues(strconv.Itoa(resp.StatusCode), r.Method).Add(1)
}
