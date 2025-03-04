package db

import (
	"log"

	"github.com/IBM/sarama"
)

func InitKafka(brokers []string) sarama.SyncProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Version = sarama.V4_0_0_0

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Cannot connect to Kafka: %v", err)
	}

	return producer
}
