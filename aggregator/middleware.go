package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"

	"github.com/KurobaneShin/tolling/types"
)

type MetricsMiddleware struct {
	next Aggregator

	reqCounterAggErr  prometheus.Counter
	reqCounterCalcErr prometheus.Counter
	reqCounterAgg     prometheus.Counter
	reqCounterCalc    prometheus.Counter
	reqLatencyAgg     prometheus.Histogram
	reqLatencyCalc    prometheus.Histogram
}

func NewMetricsMiddleware(next Aggregator) *MetricsMiddleware {
	reqCounterAgg := promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "aggregator_request_counter",
			Name:      "aggregate",
		},
	)
	reqCounterCalc := promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "aggregator_request_counter",
			Name:      "calculate",
		},
	)
	reqCounterCalcErr := promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "aggregator_request_counter",
			Name:      "calculate_error",
		},
	)

	reqCounterAggErr := promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: "aggregator_request_counter",
			Name:      "aggregator_error",
		},
	)
	reqLatencyAgg := promauto.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "aggregator_request_latency",
			Name:      "aggregate",
			Buckets:   []float64{0.1, 0.5, 1},
		},
	)

	reqLatencyCalc := promauto.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "aggregator_request_latency",
			Name:      "calculate",
			Buckets:   []float64{0.1, 0.5, 1},
		},
	)
	return &MetricsMiddleware{
		next:              next,
		reqCounterAgg:     reqCounterAgg,
		reqCounterCalc:    reqCounterCalc,
		reqLatencyAgg:     reqLatencyAgg,
		reqLatencyCalc:    reqLatencyCalc,
		reqCounterAggErr:  reqCounterAggErr,
		reqCounterCalcErr: reqCounterCalcErr,
	}
}

func (m *MetricsMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		m.reqLatencyAgg.Observe(time.Since(start).Seconds())
		m.reqCounterAgg.Inc()
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	if err != nil {
		m.reqCounterAggErr.Inc()
	}
	return
}

func (m *MetricsMiddleware) CalculateInvoice(obuId int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		m.reqLatencyCalc.Observe(time.Since(start).Seconds())
		m.reqCounterCalc.Inc()
	}(time.Now())
	inv, err = m.next.CalculateInvoice(obuId)
	if err != nil {
		m.reqCounterCalcErr.Inc()
	}
	return
}

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("AggregateDistance")
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}

func (m *LogMiddleware) CalculateInvoice(obuId int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			distance float64
			amount   float64
			id       int
		)
		if inv != nil {
			distance, amount, id = inv.TotalDistance, inv.TotalAmount, inv.OBUID
		}
		logrus.WithFields(logrus.Fields{
			"took":        time.Since(start),
			"obuID":       id,
			"totalDist":   distance,
			"totalAmount": amount,
			"err":         err,
		}).Info("AggregateDistance")
	}(time.Now())
	inv, err = m.next.CalculateInvoice(obuId)
	return
}
