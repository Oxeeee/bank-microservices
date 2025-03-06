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

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("cannot upgrade:", err)
		return
	}

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()
	log.Println("new client connected")

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}

	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()
	conn.Close()
	log.Println("client disconnected")
}

func broadcastMessage(data []byte) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("error while sending message:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}

type Consumer struct{}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var data map[string]interface{}
		err := json.Unmarshal(message.Value, &data)
		if err != nil {
			log.Printf("cannot deserialization: %v", err)
			continue
		}
		
		fmt.Printf("received new message from topic '%s': %v\n", message.Topic, data)
		
		broadcastMessage(message.Value)
		
		session.MarkMessage(message, "")
	}
	return nil
}


func startKafkaConsumer(ctx context.Context, brokers []string, groupID string, topics []string) {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer := &Consumer{}
	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("error while creating consumer group: %v", err)
	}

	go func() {
		defer client.Close()
		for {
			err := client.Consume(ctx, topics, consumer)
			if err != nil {
				log.Fatalf("error while consuming messages: %v", err)
			}
			
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startKafkaConsumer(ctx, brokers, groupID, topics)

	http.HandleFunc("/ws", wsHandler)
	go func() {
		log.Println("Websocket server started :8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("HTTP Server error: %v", err)
		}
	}()

	log.Println("Notification Service started. Listening Kafka-topics and WebSocket conns...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	cancel()
}
