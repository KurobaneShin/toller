package main

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/sirupsen/logrus"

	"github.com/KurobaneShin/tolling/types"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
}

func NewKafkaConsumer(topic string, svc CalculatorServicer) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
		"group.id":          "foo",
	})
	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	return &KafkaConsumer{
		consumer:    c,
		calcService: svc,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("kafka transport started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume error %s", err)
			continue
		}
		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON seriaization error: %s", err)
			continue
		}

		distance, err := c.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("calculation error: %s", err)
			continue

		}

		_ = distance
	}
}
