package gateway

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	v1 "homework1/internal/api/v1"
	"log"
	"net/http"
)

// RunGateway запускает HTTP-gateway, который работает как прокси для gRPC сервера.
func RunGateway(ctx context.Context, grpcEndpoint, httpEndpoint string) error {
	// Создание нового gRPC-Gateway мультиплексора.
	mux := runtime.NewServeMux()

	// Параметры подключения к gRPC серверу.
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()), // Использование небезопасных данных для тестов.
	}

	// Регистрация всех сервисов.
	if err := registerServices(ctx, mux, grpcEndpoint, opts); err != nil {
		log.Fatalf("Не удалось зарегистрировать сервисы HTTP-gateway: %v", err)
	}

	// Добавление CORS middleware для обработки запросов с других доменов.
	handler := corsMiddleware(mux)

	// Логирование и запуск HTTP-сервера.
	log.Printf("HTTP Gateway запущен на %s, проксирует к gRPC на %s", httpEndpoint, grpcEndpoint)
	return http.ListenAndServe(httpEndpoint, handler)
}

// registerServices регистрирует единый APIService для HTTP-Gateway.
func registerServices(ctx context.Context, mux *runtime.ServeMux, grpcEndpoint string, opts []grpc.DialOption) error {
	// Регистрируем единый APIService для всех сервисов.
	if err := v1.RegisterAPIServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts); err != nil {
		return err
	}
	return nil
}

// corsMiddleware добавляет поддержку CORS, чтобы позволить запросы с других доменов.
func corsMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Установка заголовков для поддержки CORS.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Если запрос метода OPTIONS (предварительный запрос), отвечаем сразу.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Прокси запросы к следующему обработчику.
		handler.ServeHTTP(w, r)
	})
}
