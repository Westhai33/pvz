package kafka

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestMain инициализирует тестовую среду перед запуском тестов
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// TestNewProducer проверяет создание нового продюсера Kafka
func TestNewProducer(t *testing.T) {
	t.Parallel()
	brokers := []string{"localhost:9092"}
	producer, err := NewProducer(brokers, "test-topic")
	assert.NoError(t, err)
	assert.NotNil(t, producer)
	defer producer.Close()
}

// TestSendOrderMessage проверяет отправку сообщения в Kafka
func TestSendOrderMessage(t *testing.T) {
	t.Parallel()
	brokers := []string{"localhost:9092"}
	producer, err := NewProducer(brokers, "test-topic")
	assert.NoError(t, err)
	defer producer.Close()

	orderMessage := OrderMessage{
		TimeStamp:   time.Now(),
		Operation:   "create",
		OrderID:     1,
		Description: "Test order",
	}

	err = producer.SendOrderMessage(orderMessage)
	assert.NoError(t, err)

}

// TestClose проверяет корректное закрытие продюсера
func TestClose(t *testing.T) {
	t.Parallel()
	brokers := []string{"localhost:9092"}
	producer, err := NewProducer(brokers, "test-topic")
	assert.NoError(t, err)

	err = producer.Close()
	assert.NoError(t, err)
}
