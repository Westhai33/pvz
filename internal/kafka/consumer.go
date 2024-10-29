package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

// NotifierHandler представляет обработчик для consumer группы
type NotifierHandler struct{}

// Setup вызывается при запуске consumer группы
func (NotifierHandler) Setup(_ sarama.ConsumerGroupSession) error {
	log.Println("Consumer группа инициализирована")
	return nil
}

// Cleanup вызывается при завершении работы consumer группы
func (NotifierHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	log.Println("Consumer группа завершила работу")
	return nil
}

// ConsumeClaim отвечает за обработку сообщений из Kafka
func (NotifierHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Получено сообщение: Key = %s, Topic = %s, Partition = %d, Offset = %d",
			string(message.Key), message.Topic, message.Partition, message.Offset)

		err := processMessage(message)
		if err != nil {
			log.Printf("Ошибка при обработке сообщения с Offset = %d: %v. Начинаем повторные попытки.", message.Offset, err)

			for i := 1; i <= 3; i++ {
				time.Sleep(2 * time.Second)
				log.Printf("Повторная попытка %d для сообщения с Offset = %d", i, message.Offset)
				err = processMessage(message)
				if err == nil {
					log.Printf("Сообщение успешно обработано после %d попытки. Offset = %d", i, message.Offset)
					break
				}
				log.Printf("Ошибка при повторной обработке сообщения с Offset = %d на %d попытке: %v", message.Offset, i, err)
			}

			if err != nil {
				log.Printf("Сообщение не удалось обработать после 3 попыток. Offset = %d. Пропускаем сообщение.", message.Offset)
				continue
			}
		}

		sess.MarkMessage(message, "")
		sess.Commit()

		log.Printf("Смещение зафиксировано и сообщение с Offset = %d успешно обработано", message.Offset)
	}
	return nil
}

// processMessage обрабатывает сообщение из Kafka
func processMessage(message *sarama.ConsumerMessage) error {
	var orderMsg OrderMessage
	err := json.Unmarshal(message.Value, &orderMsg)
	if err != nil {
		return fmt.Errorf("json.Unmarshal error: %v", err)
	}

	switch orderMsg.Operation {
	case "error":
		return fmt.Errorf("произошла ошибка во время обработки заказа ID = %d", orderMsg.OrderID)
	case "create":
		log.Printf("Обработан заказ ID = %d, операция: create", orderMsg.OrderID)
	case "update":
		log.Printf("Обновлен заказ ID = %d, операция: update", orderMsg.OrderID)
	case "delete":
		log.Printf("Удален заказ ID = %d, операция: delete", orderMsg.OrderID)
	case "create_return":
		log.Printf("Обработан возврат для заказа ID = %d, операция: create_return", orderMsg.OrderID)
	case "update_return":
		log.Printf("Обновлен возврат для заказа ID = %d, операция: update_return", orderMsg.OrderID)
	case "delete_return":
		log.Printf("Удален возврат для заказа ID = %d, операция: delete_return", orderMsg.OrderID)
	default:
		return fmt.Errorf("неизвестная операция: %s для заказа ID = %d", orderMsg.Operation, orderMsg.OrderID)
	}

	log.Printf("Сообщение успешно обработано: %+v", orderMsg)
	return nil
}

// NewConsumerGroup создаёт новый Consumer Group для Kafka
func NewConsumerGroup(brokers []string, groupID string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0

	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRange(),
	}

	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания consumer группы: %w", err)
	}

	return consumerGroup, nil
}

// StartConsumer запускает процесс потребления сообщений Kafka в отдельной горутине
func StartConsumer(consumerGroup sarama.ConsumerGroup, topics []string, handler NotifierHandler) {
	go func() {
		for {
			if err := consumerGroup.Consume(context.Background(), topics, handler); err != nil {
				log.Fatalf("Ошибка при потреблении сообщений: %v", err)
			}
		}
	}()
}
