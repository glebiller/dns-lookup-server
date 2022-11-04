package dnslookup

import (
	"fmt"
	"net/http"
	"time"
)

func NewAccessLogMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// TODO the format could be improved and following a more standard convention
		// Specifically with better time format
		//nolint:forbidigo // Should be using the log from the server instead
		fmt.Printf("%s [%s] %s %s\n", time.Now(), r.Method, r.URL, r.RemoteAddr)
		handler.ServeHTTP(rw, r)
	})
}
