package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/KurobaneShin/tolling/types"
)

type HTTPMetricHandler struct {
	reqCounter prometheus.Counter
}

func newHTTPMetricsHandler(reqName string) *HTTPMetricHandler {
	reqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_counter"),
		Name:      "aggregator",
	})

	return &HTTPMetricHandler{
		reqCounter: reqCounter,
	}
}

func (h *HTTPMetricHandler) instrument(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.reqCounter.Inc()
		next(w, r)
	}
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		obuId := r.URL.Query().Get("obu")
		if len(obuId) == 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing obu query param"})
			return
		}

		obuID, err := strconv.Atoi(obuId)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "obuId is not a number"})
			return
		}

		inv, err := svc.CalculateInvoice(obuID)
		if err != nil {

			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, inv)
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance

		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}
