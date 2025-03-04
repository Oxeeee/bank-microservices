package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

// Consumer реализует интерфейс ConsumerGroupHandler
type Consumer struct{}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Здесь можно выполнить инициализацию
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	// Здесь можно выполнить очистку
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var data map[string]interface{}
		err := json.Unmarshal(message.Value, &data)
		if err != nil {
			log.Printf("Ошибка десериализации: %v", err)
			continue
		}
		fmt.Printf("Получено сообщение из топика '%s': %v\n", message.Topic, data)
		// Отмечаем сообщение как обработанное
		session.MarkMessage(message, "")
	}
	return nil
}

func main() {
	brokers := []string{"localhost:9092"}
	groupID := "notification_service_group"
	topics := []string{"auth_topic", "transaction_topic", "billing_topic"}

	config := sarama.NewConfig()
	// Используем версию брокера не ниже 2.1.0
	config.Version = sarama.V4_0_0_0
	// Если нет смещений, читаем с начала
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer := &Consumer{}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Ошибка создания consumer group: %v", err)
	}
	defer client.Close()

	// Запускаем потребление в отдельной горутине
	go func() {
		for {
			err := client.Consume(ctx, topics, consumer)
			if err != nil {
				log.Fatalf("Ошибка в процессе потребления: %v", err)
			}
			// Если контекст завершён, выходим
			if ctx.Err() != nil {
				return
			}
		}
	}()

	fmt.Println("Notification Service запущен и слушает топики...")

	// Ожидание сигнала для завершения работы
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	cancel()
}
