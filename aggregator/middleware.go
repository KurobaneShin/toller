package main

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/KurobaneShin/tolling/types"
)

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
		)
		if inv != nil {
			distance, amount = inv.TotalDistance, inv.TotalAmount
		}
		logrus.WithFields(logrus.Fields{
			"took":        time.Since(start),
			"obuID":       inv.OBUID,
			"totalDist":   distance,
			"totalAmount": amount,
			"err":         err,
		}).Info("AggregateDistance")
	}(time.Now())
	inv, err = m.next.CalculateInvoice(obuId)
	return
}
