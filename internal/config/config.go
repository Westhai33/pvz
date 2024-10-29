package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// Config представляет структуру для конфигурации сервиса
type Config struct {
	KafkaBrokers []string // Список брокеров Kafka
	KafkaGroupID string   // ID группы Kafka
	KafkaTopic   string   // Топик Kafka
	DBUser       string   // Пользователь базы данных
	DBPassword   string   // Пароль базы данных
	DBName       string   // Имя базы данных
	DBHost       string   // Хост базы данных
	DBPort       string   // Порт базы данных
	GrpcPort     string   // Порт gRPC
	HttpPort     string   // Порт HTTP
	RedisAddr    string   // Адрес Redis
	RedisDB      int      // Номер базы данных Redis
	MetricsAddr  string   // Адрес сервера метрик
	TracingURL   string   // URL для экспорта трейсинга (Jaeger или другой провайдер)
	ServiceName  string   // Название сервиса для трейсинга
}

// LoadConfig загружает конфигурацию из переменных окружения
func LoadConfig() *Config {
	kafkaBrokers := getEnvAsSlice("KAFKA_BROKERS", []string{"localhost:9092"})
	kafkaGroupID := getEnv("KAFKA_GROUP_ID", "notifier_group")
	kafkaTopic := getEnv("KAFKA_TOPIC", "pvz.events-log")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "postgres")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	grpcPort := getEnv("GRPC_PORT", "50051")
	httpPort := getEnv("HTTP_PORT", "8080")
	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisDB := getEnvAsInt("REDIS_DB", 0)
	metricsAddr := getEnv("METRICS_ADDR", ":8099")
	tracingURL := getEnv("TRACING_URL", "http://localhost:14268/api/traces")
	serviceName := getEnv("SERVICE_NAME", "my-go-service")

	log.Printf("Конфигурация загружена: brokers=%v, groupID=%s, topic=%s", kafkaBrokers, kafkaGroupID, kafkaTopic)
	log.Printf("Настройки базы данных: user=%s, dbname=%s, host=%s, port=%s", dbUser, dbName, dbHost, dbPort)
	log.Printf("gRPC порт: %s, HTTP порт: %s", grpcPort, httpPort)
	log.Printf("Redis: addr=%s, db=%d", redisAddr, redisDB)
	log.Printf("Metrics: addr=%s", metricsAddr)
	log.Printf("Tracing: url=%s, service=%s", tracingURL, serviceName)

	return &Config{
		KafkaBrokers: kafkaBrokers,
		KafkaGroupID: kafkaGroupID,
		KafkaTopic:   kafkaTopic,
		DBUser:       dbUser,
		DBPassword:   dbPassword,
		DBName:       dbName,
		DBHost:       dbHost,
		DBPort:       dbPort,
		GrpcPort:     grpcPort,
		HttpPort:     httpPort,
		RedisAddr:    redisAddr,
		RedisDB:      redisDB,
		MetricsAddr:  metricsAddr,
		TracingURL:   tracingURL,
		ServiceName:  serviceName,
	}
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// getEnvAsInt возвращает значение переменной окружения как целое число или значение по умолчанию
func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		} else {
			log.Printf("Ошибка при преобразовании переменной окружения %s: %v", key, err)
		}
	}
	return fallback
}

// getEnvAsSlice возвращает значение переменной окружения как срез строк или значение по умолчанию
func getEnvAsSlice(key string, fallback []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return splitAndTrim(value, ",")
	}
	return fallback
}

// splitAndTrim разбивает строку по разделителю и удаляет пробелы
func splitAndTrim(str, sep string) []string {
	parts := strings.Split(str, sep)
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
