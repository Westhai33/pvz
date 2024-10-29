package main

import (
	"context"
	"fmt"
	"homework1/internal/model"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"homework1/internal/api/v1"
	"homework1/internal/cache"
	"homework1/internal/client"
	"homework1/internal/config"
	"homework1/internal/dao"
	"homework1/internal/gateway"
	"homework1/internal/kafka"
	"homework1/internal/metrics"
	"homework1/internal/pool"
	"homework1/internal/server"
	"homework1/internal/service"
	"homework1/internal/tracing" // Импортируем пакет для трейсинга
	"homework1/internal/view"
	"homework1/utility/logger"
	"homework1/utility/validation"
)

func main() {
	err := addServersSectionToOpenAPI()
	if err != nil {
		log.Fatalf("Ошибка при добавлении секции servers в OpenAPI: %v\n", err)
	}

	initLoggerAndValidator()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gracefulShutdown(cancel)

	cfg := config.LoadConfig()

	shutdownTracer := tracing.InitTracer(cfg.ServiceName, cfg.TracingURL)
	defer func() {
		if err := shutdownTracer(ctx); err != nil {
			log.Fatalf("Ошибка завершения трейсинга: %v", err)
		}
	}()
	log.Println("Трейсинг инициализирован:", cfg.TracingURL)

	initDatabase(cfg)
	defer dao.Closedb()

	dbPool := dao.GetPool()

	redisClient := initRedis(cfg)
	defer redisClient.Close()

	go metrics.StartMetricsServer(cfg.MetricsAddr)
	log.Println("Приложение запущено. Адрес метрик:", cfg.MetricsAddr)

	kafkaProducer, err := kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
	if err != nil {
		log.Fatalf("Ошибка при инициализации Kafka Producer: %v", err)
	}
	defer kafkaProducer.Close()

	wp := pool.NewWorkerPool(2)

	orderService, userService, packagingService, returnService, returnReasonService, statusService := initServices(dbPool, wp, kafkaProducer, redisClient)

	startServers(ctx, cfg, orderService, userService, packagingService, returnService, returnReasonService, statusService, wp)

	startBackgroundTask(ctx, orderService, wp)

	grpcClients := setupGRPCClients(cfg.GrpcPort)
	view.RunInteractiveMode(ctx, grpcClients, wp)
}

// Функция для инициализации Redis-клиента
func initRedis(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
		DB:   cfg.RedisDB,
	})
	return rdb
}

// Функция для инициализации сервисов с Redis-кэшем и Kafka Producer
func initServices(dbPool *pgxpool.Pool, wp *pool.WorkerPool, kafkaProd *kafka.Producer, redisClient *redis.Client) (
	*service.OrderService, *service.UserService, *service.PackagingService, *service.ReturnService, *service.ReturnReasonService, *service.StatusService) {

	cacheConfig := cache.CacheConfig{
		DefaultTTL: 10 * time.Minute, // Установите значение по умолчанию для TTL
	}

	orderCache := cache.NewRedisCache[string, model.Order](redisClient, cacheConfig)
	userCache := cache.NewRedisCache[string, model.User](redisClient, cacheConfig)
	packagingCache := cache.NewRedisCache[string, model.PackagingOption](redisClient, cacheConfig)
	returnCache := cache.NewRedisCache[string, model.Return](redisClient, cacheConfig)
	returnReasonCache := cache.NewRedisCache[string, model.ReturnReason](redisClient, cacheConfig)
	statusCache := cache.NewRedisCache[string, model.Status](redisClient, cacheConfig)

	orderService := service.NewOrderService(dbPool, wp, kafkaProd, orderCache)
	userService := service.NewUserService(dbPool, wp, userCache)
	packagingService := service.NewPackagingService(dbPool, wp, packagingCache)
	returnService := service.NewReturnService(dbPool, wp, kafkaProd, returnCache)
	returnReasonService := service.NewReturnReasonService(dbPool, wp, returnReasonCache)
	statusService := service.NewStatusService(dbPool, wp, statusCache)

	return orderService, userService, packagingService, returnService, returnReasonService, statusService
}

// Функция для добавления секции servers в OpenAPI YAML
func addServersSectionToOpenAPI() error {
	openAPIFile := "./internal/api/v1/openapi.yaml"

	data, err := os.ReadFile(openAPIFile)
	if err != nil {
		return fmt.Errorf("ошибка при чтении файла: %w", err)
	}

	content := string(data)

	if strings.Contains(content, "info:") {
		content = strings.Replace(content, `title: ""`, `title: "Order Management API"`, 1)
		content = strings.Replace(content, `version: 0.0.1`, `version: 1.0.0`, 1)
		content = strings.Replace(content, `description: ""`, `description: "API для управления заказами, упаковками, возвратами, статусами и пользователями"`, 1)
	} else {
		infoSection := `
info:
  title: "Order Management API"
  description: "API для управления заказами, упаковками, возвратами, статусами и пользователями"
  version: "1.0.0"
`
		content = infoSection + content
	}

	if !strings.Contains(content, "servers:") {
		serversSection := `
servers:
  - url: http://localhost:8080
    description: "Local HTTP Gateway"
`
		content += serversSection
	}

	err = os.WriteFile(openAPIFile, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("ошибка при записи в файл: %w", err)
	}

	return nil
}

func initLoggerAndValidator() {
	logger.InitLogger()
	validation.InitValidator()
}

// Функция для загрузки конфигурации базы данных
func initDatabase(cfg *config.Config) {
	port, err := strconv.Atoi(cfg.DBPort)
	if err != nil {
		log.Fatalf("Некорректный формат порта: %v", err)
	}
	dao.Initdb(cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBHost, port)
}

// startServers запускает gRPC и HTTP Gateway серверы
func startServers(
	ctx context.Context, cfg *config.Config,
	orderService *service.OrderService, userService *service.UserService,
	packagingService *service.PackagingService, returnService *service.ReturnService,
	returnReasonService *service.ReturnReasonService, statusService *service.StatusService,
	wp *pool.WorkerPool) {

	go func() {
		if err := startGRPCServer(cfg.GrpcPort, orderService, userService, packagingService, returnService, returnReasonService, statusService, wp); err != nil {
			log.Fatalf("Ошибка при запуске gRPC сервера: %v", err)
		}
		log.Println("gRPC сервер завершил работу")
	}()

	go func() {
		log.Printf("Запуск HTTP Gateway на порту %s", cfg.HttpPort)
		if err := gateway.RunGateway(ctx, "localhost:"+cfg.GrpcPort, "localhost:"+cfg.HttpPort); err != nil {
			log.Fatalf("Ошибка при запуске HTTP Gateway: %v", err)
		}
	}()
}

// startBackgroundTask запускает фоновую задачу для проверки просроченных заказов
func startBackgroundTask(ctx context.Context, orderService *service.OrderService, wp *pool.WorkerPool) {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if ctx.Err() != nil {
					log.Println("Контекст уже отменен:", ctx.Err())
					return
				}
				wp.SubmitTask(func() {
					err := orderService.CheckExpiredOrders(ctx)
					if err != nil {
						log.Printf("Ошибка при проверке просроченных заказов: %v", err)
					} else {
						log.Println("Проверка просроченных заказов завершена успешно.")
					}
				})
			case <-ctx.Done():
				log.Println("Завершение фоновой задачи проверки просроченных заказов.")
				return
			}
		}
	}()
}

func setupGRPCClients(grpcPort string) *client.APIServiceClientWrapper {
	clientOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(50*1024*1024),
			grpc.MaxCallSendMsgSize(50*1024*1024),
		),
	}

	grpcConn, err := grpc.NewClient("localhost:"+grpcPort, clientOptions...)
	if err != nil {
		log.Fatalf("Не удалось подключиться к gRPC серверу: %v", err)
	}

	grpcClientWrapper, err := client.NewAPIServiceClientWrapper(grpcConn)
	if err != nil {
		log.Fatalf("Ошибка при создании обертки gRPC клиента: %v", err)
	}

	return grpcClientWrapper
}

// startGRPCServer запускает gRPC сервер
func startGRPCServer(grpcPort string, orderService *service.OrderService, userService *service.UserService, packagingService *service.PackagingService, returnService *service.ReturnService, returnReasonService *service.ReturnReasonService, statusService *service.StatusService, workerPool *pool.WorkerPool) error {
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		return fmt.Errorf("не удалось начать слушать порт %s: %w", grpcPort, err)
	}

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(50*1024*1024),
		grpc.MaxSendMsgSize(50*1024*1024),
	)

	v1.RegisterAPIServiceServer(grpcServer, server.NewAPIServiceServer(userService, orderService, packagingService, returnService, returnReasonService, statusService, workerPool))

	log.Printf("gRPC сервер запущен на порту %s", grpcPort)
	return grpcServer.Serve(lis)
}

// gracefulShutdown добавляет обработку системных сигналов для корректного завершения работы
func gracefulShutdown(cancelFunc context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signals
		log.Println("Получен сигнал завершения. Завершаем работу...")
		cancelFunc()
		os.Exit(0)
	}()
}
