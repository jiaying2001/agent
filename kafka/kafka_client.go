package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jiaying2001/agent/store"
	"os"
)

var kp *kafka.Producer

func init() {
	var err error
	kp, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": store.C.Kafka.HostName + `:` + store.C.Kafka.Port,
		"client.id":         store.C.Kafka.Client.Id,
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
}

func Send(topic string, msg []byte) {
	// Delivery report handler for produced messages
	go func() {
		for e := range kp.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	kp.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          msg,
	}, nil)
}
