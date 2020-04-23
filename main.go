
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/GSabadini/go-prometheus/middleware"
	"github.com/GSabadini/go-prometheus/prometheus"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/negroni"
)

func main() {
	var r = mux.NewRouter()

	metricService, err := prometheus.NewPrometheusService()
	if err != nil {
		log.Fatal(err.Error())
	}

	var n = negroni.New(
		middleware.Metrics(metricService),
	)

	n.UseHandler(r)

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/", info)

	log.Println("Start HTTP server :3001")
	if err := http.ListenAndServe(":3001", nil); err != nil {
		panic(err)
	}
}

func info(w http.ResponseWriter, req *http.Request) {
	hostname, _ := os.Hostname()

	data := struct {
		Hostname string      `json:"hostname,omitempty"`
		IP       string      `json:"ip,omitempty"`
		Headers  http.Header `json:"headers,omitempty"`
		URL      string      `json:"url,omitempty"`
		Host     string      `json:"host,omitempty"`
		Method   string      `json:"method,omitempty"`
	}{
		Hostname: hostname,
		IP:       getIP(req),
		Headers:  req.Header,
		URL:      req.URL.RequestURI(),
		Host:     req.Host,
		Method:   req.Method,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getIP(req *http.Request) string {
	forwarded := req.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}

	return req.RemoteAddr
}