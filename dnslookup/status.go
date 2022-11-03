package dnslookup

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime/debug"
	"time"
)

const kubernetesServiceHostKey = "KUBERNETES_SERVICE_HOST"

// ModelStatus model status
type ModelStatus struct {
	Version    string `json:"version"`
	Date       int64  `json:"date"`
	Kubernetes bool   `json:"kubernetes"`
}

func NewStatusMiddleware(handler http.Handler) http.Handler {
	runningOnKubernetes := false
	// Could also be check for /var/run/secrets/kubernetes.io, however with automountServiceAccountToken: false
	// this would not work, hence I prefer to check for the env variable that is automatically set by K8s.
	if _, ok := os.LookupEnv(kubernetesServiceHostKey); ok {
		runningOnKubernetes = true
	}
	return StatusMiddleware{
		nextHandler:         handler,
		runningOnKubernetes: runningOnKubernetes,
	}
}

type StatusMiddleware struct {
	nextHandler         http.Handler
	runningOnKubernetes bool
}

func (m StatusMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" && r.URL.Path == "/" {
		m.writeApplicationStatus(rw, r)
	} else {
		m.nextHandler.ServeHTTP(rw, r)
	}
}

func (m StatusMiddleware) writeApplicationStatus(rw http.ResponseWriter, r *http.Request) {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		rw.WriteHeader(500)
		return
	}

	status := ModelStatus{
		Version:    buildInfo.Main.Version,
		Date:       time.Now().Unix(),
		Kubernetes: m.runningOnKubernetes,
	}
	output, err := json.Marshal(status)
	if err != nil {
		// TODO improve error message
		rw.WriteHeader(500)
		return
	}

	_, err = rw.Write(output)
	if err != nil {
		// TODO improve error message
		rw.WriteHeader(500)
		return
	}
}
