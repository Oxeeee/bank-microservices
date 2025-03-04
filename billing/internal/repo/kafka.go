package repo

import (
	"encoding/json"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/Oxeeee/bank-microservices/billing/internal/models/domain"
)

type BillingKafkaRepo interface {
	SendPaymentStatus(payment *domain.BillPayment) error
}

type billingKafkaRepository struct {
	producer sarama.SyncProducer
	topic    string
	log      *slog.Logger
}

func NewBillingKafkaRepo(producer sarama.SyncProducer, topic string, log *slog.Logger) BillingKafkaRepo {
	return &billingKafkaRepository{
		producer: producer,
		topic:    topic,
		log:      log,
	}
}

func (r *billingKafkaRepository) SendPaymentStatus(payment *domain.BillPayment) error {
	const op = "kafkarepo.sendPaymentStatus"
	log := r.log.With(slog.String("op", op))
	msgBytes, err := json.Marshal(payment)
	if err != nil {
		log.Error("cannot marshal json", "error", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: r.topic,
		Value: sarama.ByteEncoder(msgBytes),
	}

	partition, offset, err := r.producer.SendMessage(msg)
	if err != nil {
		log.Error("cannot send message", "error", err)
		return err
	}

	log.Debug("payment status sent", "topic", r.topic, "partition", partition, "offset", offset)
	return nil
}
