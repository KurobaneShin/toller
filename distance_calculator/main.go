package main

import (
	"log"
)

type DistanceCalculator struct{}

const kafkaTopic = "obudata"

func main() {
	var svc CalculatorServicer
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
