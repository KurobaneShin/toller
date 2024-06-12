package main

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type DistanceCalculator struct{}

func main() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
		"group.id":          "foo",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatal(err)
	}

	consumer.SubscribeTopics([]string{"obudata"}, nil)

	for {
		msg, err := consumer.ReadMessage(time.Second)
		if err != nil {
			log.Println("error")
			continue
		}

		fmt.Printf("message on %s: %s\n", msg.TopicPartition, string(msg.Value))

	}
}
