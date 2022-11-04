package dnslookup

import (
	"net/http"
)

func NewHealthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == httpGetMethod && r.URL.Path == "/health" {
			// Health could check if the database is up,
			// however the app will already return 500 if the log does not write correctly to the DB
			_, err := rw.Write([]byte("{}"))
			if err != nil {
				// TODO better handle error message
				rw.WriteHeader(httpServerError)
				return
			}
		} else {
			handler.ServeHTTP(rw, r)
		}
	})
}
