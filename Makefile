# Название бинарного файла после сборки
BINARY_NAME = homework.exe

# Директории для .proto файлов и сгенерированных файлов
PROTO_DIR = ./proto
GEN_DIR = ./
SWAGGER_DIR = ./swagger

# Go-пакеты в проекте
GOPACKAGES := $(shell go list ./...)

# Таргет для сборки приложения
build: lint generate swagger
	@echo "Building the application..."
	go build -o $(BINARY_NAME) ./main.go
	@echo "The application has been successfully built: $(BINARY_NAME)"

# Таргет для обновления и установки зависимостей
deps:
	@echo "Updating and installing dependencies..."
	go mod tidy
	go mod download
	@echo "Dependencies are updated and installed."

# Таргет для запуска приложения
run: build
	@echo "Launching the application..."
	./$(BINARY_NAME)

# Линтинг с использованием golangci-lint
lint: install-linters
	@echo "Launching golangci-lint..."
	golangci-lint run --config ./manifest/config/.golangci.yml ./...

# Установка golangci-lint
install-linters:
	@echo "Installing golangci-lint..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Таргет для генерации кода из .proto файлов
generate: install-proto
	@echo "Generating Go code from .proto files..."
	@if not exist "$(GEN_DIR)" mkdir "$(GEN_DIR)"
	@protoc -I=$(PROTO_DIR) \
		--go_out=$(GEN_DIR) \
		--go-grpc_out=$(GEN_DIR) \
		--grpc-gateway_out=$(GEN_DIR) \
		--validate_out="lang=go:$(GEN_DIR)" \
		$(PROTO_DIR)/*.proto
	@echo "Proto files have been successfully compiled."

# Установка protoc-gen-go и protoc-gen-go-grpc
install-proto:
	@echo "Installing protoc-gen-go and protoc-gen-go-grpc..."
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest

# Таргет для чистки скомпилированных файлов
clean:
	@echo "Cleaning compiled files and generated code..."
	rm -f $(BINARY_NAME)
	rm -rf $(GEN_DIR)
	rm -rf $(SWAGGER_DIR)
	@echo "The cleanup is complete."

# Таргет для расчета покрытия тестов
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Test coverage report generated: coverage.html"

# Таргет для запуска тестов
test:
	@echo "Running tests..."
	go test ./...

# Таргет для проверки статуса тестов
test-status:
	@echo "Checking test status..."
	go test -v ./...

# Таргет для применения миграций базы данных
migrate:
	@echo "Applying database migrations..."
	goose -dir db/migrations postgres "user=postgres password=postgres dbname=postgres sslmode=disable" up
	@echo "Database migrations applied successfully."

# Таргет для сброса миграций базы данных
reset:
	@echo "Resetting database migrations..."
	goose -dir db/migrations postgres "user=postgres password=postgres dbname=postgres sslmode=disable" reset
	@echo "Database migrations reset successfully."

# Таргет для отката миграций базы данных
down:
	@echo "Rolling back the last migration..."
	goose -dir db/migrations postgres "user=postgres password=postgres dbname=postgres sslmode=disable" down
	@echo "Last migration rolled back successfully."

# Таргет для запуска контейнеров Docker
docker-up:
	@echo "Starting Docker containers..."
	docker-compose -f docker-compose.yml up -d
	@echo "Docker containers started successfully."

# Таргет для доступа к базе данных в контейнере
db-shell:
	@echo "Accessing PostgreSQL shell in the Docker container..."
	docker exec -it route256_db psql -U postgres -d postgres
	@echo "Exited PostgreSQL shell."
swagger: generate
	@echo "Generating Swagger specifications..."
	if not exist ./internal/api/v1 mkdir ./internal/api/v1
	protoc -I=$(PROTO_DIR) \
		--openapi_out=./internal/api/v1/ \
		$(PROTO_DIR)/*.proto
	@echo "Swagger specifications generated and saved to swagger.json."
# Основной таргет по умолчанию
all: deps generate build run

# Определяем .PHONY для целей
.PHONY: build deps run lint clean all install-linters generate install-proto swagger install-swagger coverage test test-status migrate reset down docker-up db-shell
