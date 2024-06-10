package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"

	"github.com/KurobaneShin/tolling/types"
)

var kafkaTopic = "obudata"

func main() {
	recv, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", recv.handleWs)
	http.ListenAndServe(":30000", nil)
}

type DataReceiver struct {
	msgch    chan types.OBUData
	conn     *websocket.Conn
	producer *kafka.Producer
}

func NewDataReceiver() (*DataReceiver, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:29092",
		"client.id":         "foo",
		"acks":              "all",
	})
	if err != nil {
		return nil, err
	}
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				fmt.Printf("message received: %+v\n", ev.TopicPartition)

			case *kafka.Error:
				fmt.Printf("%+s\n", e)
			}
		}
	}()
	return &DataReceiver{
		msgch:    make(chan types.OBUData, 128),
		producer: p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = dr.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kafkaTopic,
			Partition: kafka.PartitionAny,
		},
		Value: b,
	}, nil)
	return err
}

func (dr *DataReceiver) handleWs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	dr.conn = conn
	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("New OBU connected client connected")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error: ", err)
			continue
		}
		if err := dr.produceData(data); err != nil {
			fmt.Println("kafka produce error: ", err)
		}
	}
}
