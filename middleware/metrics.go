package middleware

import (
	"net/http"
	"strconv"

	"github.com/GSabadini/go-prometheus/prometheus"

	"github.com/codegangsta/negroni"
)

//Metrics to prometheus
func Metrics(mService prometheus.TypeMetric) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		appMetric := prometheus.NewHTTP(r.URL.Path, r.Method)
		appMetric.Started()
		next(w, r)
		res := w.(negroni.ResponseWriter)
		appMetric.Finished()
		appMetric.StatusCode = strconv.Itoa(res.Status())
		mService.HTTP(appMetric)
	}
}
