package httperror

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func Error(w http.ResponseWriter, r *http.Request, error string, code int, httpReqs *prometheus.CounterVec) {
	http.Error(w, error, code)
	log.Println(fmt.Sprintf("%s: %s", http.StatusText(code), error))
	httpReqs.WithLabelValues(strconv.Itoa(code), r.Method).Add(1)
}
