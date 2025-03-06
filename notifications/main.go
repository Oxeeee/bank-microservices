package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/Oxeeee/bank-microservices/notifications/config"
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
		Log.Error("cannot upgrade", "error", err)
		return
	}

	clientsMu.Lock()
	clients[conn] = true
	clientsMu.Unlock()
	Log.Debug("new client connected")

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}

	clientsMu.Lock()
	delete(clients, conn)
	clientsMu.Unlock()
	conn.Close()
	Log.Debug("client disconnected")
}

func broadcastMessage(data []byte) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			Log.Error("error while sending message", "error", err)
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
			Log.Error("cannot serizalize", "error", err)
			continue
		}

		Log.Debug("received message", "topic", message.Topic, "data", string(message.Value))

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
		Log.Error("error while creating consumer group", "error", err)
		panic(err)
	}

	go func() {
		defer client.Close()
		for {
			err := client.Consume(ctx, topics, consumer)
			if err != nil {
				Log.Error("error while consuming messages", "error", err)
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var Log *slog.Logger

func setupLogger(env string) {

	switch env {
	case envLocal:
		Log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		Log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		Log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
}

func main() {
	cfg := config.MustLoad()
	setupLogger(cfg.Env)

	brokers := cfg.Kafka.Brokers
	groupID := cfg.Kafka.ConsumerGroupID
	topics := cfg.Kafka.Topics

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startKafkaConsumer(ctx, brokers, groupID, topics)

	http.HandleFunc("/ws", wsHandler)
	go func() {
		Log.Info("HTTP server started", "host", fmt.Sprintf("%v:%v", cfg.WebSocket.Host, cfg.WebSocket.Port))
		if err := http.ListenAndServe(fmt.Sprintf(":%v", cfg.WebSocket.Port), nil); err != nil {
			panic("http server error: " + err.Error())
		}
	}()

	Log.Info("Notification Service started. Listening Kafka-topics and WebSocket conns...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	cancel()
}
