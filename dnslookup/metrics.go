package dnslookup

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func NewPrometheusMetricMiddleware(handler http.Handler) http.Handler {
	prometheusHandler := promhttp.Handler()
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/metrics" {
			prometheusHandler.ServeHTTP(rw, r)
		} else {
			handler.ServeHTTP(rw, r)
		}
	})
}
