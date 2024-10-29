package main

import (
	"homework1/internal/config"
	"homework1/internal/kafka"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	consumerGroup, err := kafka.NewConsumerGroup(cfg.KafkaBrokers, cfg.KafkaGroupID)
	if err != nil {
		log.Fatalf("Ошибка создания consumer группы: %v", err)
	}

	// Обработка системных сигналов для корректного завершения работы
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		log.Println("Получен сигнал завершения. Завершаем работу...")
		os.Exit(0)
	}()

	handler := kafka.NotifierHandler{}

	// Запуск Kafka Consumer в отдельной горутине
	kafka.StartConsumer(consumerGroup, []string{cfg.KafkaTopic}, handler)

	select {}
}
