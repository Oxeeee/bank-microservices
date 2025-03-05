package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/gorilla/websocket"
)

// Глобальные переменные для хранения подключённых клиентов
var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

// wsHandler – HTTP-хендлер для WebSocket-соединений
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Апгрейдим HTTP-соединение до WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка апгрейда:", err)
		return
	}

	// Регистрируем клиента
	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()
	log.Println("Новый клиент подключён")

	// Ждём разрыва соединения
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}

	// Удаляем клиента после разрыва соединения
	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()
	conn.Close()
	log.Println("Клиент отключён")
}

// broadcastMessage отсылает переданные данные всем подключённым клиентам
func broadcastMessage(data []byte) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("Ошибка отправки сообщения:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}

// Consumer реализует интерфейс ConsumerGroupHandler для Kafka
type Consumer struct{}

// Setup вызывается перед началом потребления сообщений
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup вызывается после завершения потребления
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim обрабатывает сообщения из топика
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var data map[string]interface{}
		err := json.Unmarshal(message.Value, &data)
		if err != nil {
			log.Printf("Ошибка десериализации: %v", err)
			continue
		}
		// Выводим полученное сообщение в консоль
		fmt.Printf("Получено сообщение из топика '%s': %v\n", message.Topic, data)
		// Отправляем сообщение всем подключённым клиентам
		broadcastMessage(message.Value)
		// Отмечаем сообщение как обработанное
		session.MarkMessage(message, "")
	}
	return nil
}

// startKafkaConsumer запускает потребление сообщений из Kafka
func startKafkaConsumer(ctx context.Context, brokers []string, groupID string, topics []string) {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0 // используем нужную версию брокера
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer := &Consumer{}
	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Ошибка создания consumer group: %v", err)
	}

	// Запускаем потребление в горутине.
	// Перемещаем defer client.Close() внутрь горутины, чтобы закрытие произошло только при завершении цикла.
	go func() {
		defer client.Close()
		for {
			err := client.Consume(ctx, topics, consumer)
			if err != nil {
				log.Fatalf("Ошибка в процессе потребления: %v", err)
			}
			// Если контекст завершён – выходим из цикла
			if ctx.Err() != nil {
				return
			}
		}
	}()
}

func main() {
	brokers := []string{"localhost:9092"}
	groupID := "notification_service_group"
	topics := []string{"auth_topic", "transaction_topic", "billing_topic"}

	// Создаем контекст для управления жизненным циклом Kafka-консьюмера
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем Kafka-консьюмер
	startKafkaConsumer(ctx, brokers, groupID, topics)

	// Регистрируем WebSocket-хендлер
	http.HandleFunc("/ws", wsHandler)
	go func() {
		log.Println("WebSocket сервер запущен на :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("Ошибка HTTP сервера: %v", err)
		}
	}()

	log.Println("Notification Service запущен. Слушаем Kafka-топики и WebSocket подключения...")

	// Ожидаем сигнал для завершения работы
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	cancel()
}
