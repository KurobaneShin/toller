package main

import (
	"log"

	"github.com/KurobaneShin/tolling/aggregator/client"
)

type DistanceCalculator struct{}

const (
	kafkaTopic         = "obudata"
	aggregatorEndpoint = "http://127.0.0.1:3000/aggregate"
)

func main() {
	var svc CalculatorServicer
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, client.NewClient(aggregatorEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
