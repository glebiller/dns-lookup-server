package dnslookup

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewPrometheusMetricMiddleware(handler http.Handler) http.Handler {
	prometheusHandler := promhttp.Handler()
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == httpGetMethod && r.URL.Path == "/metrics" {
			prometheusHandler.ServeHTTP(rw, r)
		} else {
			handler.ServeHTTP(rw, r)
		}
	})
}
