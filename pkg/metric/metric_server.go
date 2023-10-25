package metric

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Listen() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
