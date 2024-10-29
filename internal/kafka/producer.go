package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

// OrderMessage представляет структуру сообщения о заказе
type OrderMessage struct {
	TimeStamp   time.Time
	Operation   string
	OrderID     int
	Description string
}

// Producer представляет Kafka продюсера
type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewProducer создает нового продюсера Kafka
func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Idempotent = true
	config.Producer.Return.Successes = true
	config.Net.MaxOpenRequests = 1

	syncProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("sarama.NewSyncProducer: %w", err)
	}

	producer := &Producer{
		producer: syncProducer,
		topic:    topic,
	}

	return producer, nil
}

// SendOrderMessage отправляет сообщение о заказе в Kafka
func (p Producer) SendOrderMessage(message OrderMessage) error {
	// Сериализация сообщения
	msg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(msg),
	}

	partition, offset, err := p.producer.SendMessage(kafkaMsg)
	if err != nil {
		return fmt.Errorf("p.producer.SendMessage: %w", err)
	}

	log.Printf("Сообщение успешно отправлено в Kafka. Тема: %s, Раздел: %d, Смещение: %d, Сообщение: %+v\n", p.topic, partition, offset, message)

	return nil
}

// SendKafkaErrorMessage отправляет сообщение об ошибке в Kafka
func (p Producer) SendKafkaErrorMessage(operation string, orderID int, description string) error {
	errorMessage := OrderMessage{
		TimeStamp:   time.Now(),
		Operation:   operation,
		OrderID:     orderID,
		Description: fmt.Sprintf("Ошибка: %s", description),
	}

	msg, err := json.Marshal(errorMessage)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	kafkaMsg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(msg),
	}

	partition, offset, err := p.producer.SendMessage(kafkaMsg)
	if err != nil {
		return fmt.Errorf("p.producer.SendMessage: %w", err)
	}

	log.Printf("Сообщение об ошибке успешно отправлено в Kafka. Тема: %s, Раздел: %d, Смещение: %d, Сообщение: %+v\n", p.topic, partition, offset, errorMessage)

	return nil
}

// Close закрывает продюсера
func (p Producer) Close() error {
	err := p.producer.Close()
	if err != nil {
		return fmt.Errorf("p.producer.Close: %w", err)
	}

	return nil
}
