package main

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/KurobaneShin/tolling/types"
)

type LogMiddleware struct {
	next DataProducer
}

func NewLogMiddleware(next DataProducer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) ProduceData(data types.OBUData) error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":  time.Since(start),
			"obuID": data.OBUID,
			"lat":   data.Lat,
			"long":  data.Long,
		}).Info("producing to kafka")
	}(time.Now())
	return l.next.ProduceData(data)
}
