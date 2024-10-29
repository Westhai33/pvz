package service_test

import (
	"context"
	"homework1/internal/dao"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"homework1/internal/kafka"
	"homework1/internal/model"
	"homework1/internal/pool"
	"homework1/internal/service"
)

var testDB *pgxpool.Pool

// TestMain инициализирует тестовую среду перед запуском тестов
func TestMain(m *testing.M) {
	dbUrl := "postgres://postgres:postgres@localhost:5433/test_db?sslmode=disable"
	var err error
	testDB, err = pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		panic("Не удалось подключиться к тестовой базе данных: " + err.Error())
	}
	defer testDB.Close()

	os.Exit(m.Run())
}

// TestCreateOrder проверяет создание нового заказа и отправку сообщения в Kafka
func TestCreateOrder(t *testing.T) {
	brokers := []string{"localhost:9092"}
	producer, err := kafka.NewProducer(brokers, "test-topic")
	assert.NoError(t, err)
	defer producer.Close()

	wp := pool.NewWorkerPool(2)
	orderService := service.NewOrderService(testDB, wp, producer)

	ctx := context.Background()
	orderID, err := orderService.CreateOrder(ctx, 1, 1, 1, time.Now().Add(24*time.Hour), 1.0, 100.0, 10.0, 110.0, false)
	assert.NoError(t, err)
	assert.NotZero(t, orderID)

}

// TestUpdateOrder проверяет обновление заказа и отправку сообщения в Kafka
func TestUpdateOrder(t *testing.T) {
	brokers := []string{"localhost:9092"}
	producer, err := kafka.NewProducer(brokers, "test-topic")
	assert.NoError(t, err)
	defer producer.Close()

	wp := pool.NewWorkerPool(2)
	orderService := service.NewOrderService(testDB, wp, producer)

	ctx := context.Background()
	orderID, err := orderService.CreateOrder(ctx, 1, 1, 1, time.Now().Add(24*time.Hour), 1.0, 100.0, 10.0, 110.0, false)
	assert.NoError(t, err)

	order := model.Order{
		OrderID:        orderID,
		UserID:         1,
		PackagingID:    1,
		ExpirationDate: time.Now().Add(48 * time.Hour),
		Weight:         2.0,
		BaseCost:       200.0,
		PackagingCost:  20.0,
		TotalCost:      220.0,
	}
	orderService.UpdateOrder(ctx, order)

}

// TestDeleteOrder проверяет удаление заказа и отправку сообщения в Kafka
func TestDeleteOrder(t *testing.T) {
	brokers := []string{"localhost:9092"}
	producer, err := kafka.NewProducer(brokers, "test-topic")
	assert.NoError(t, err)
	defer producer.Close()

	wp := pool.NewWorkerPool(2)
	orderService := service.NewOrderService(testDB, wp, producer)

	ctx := context.Background()
	orderID, err := orderService.CreateOrder(ctx, 1, 1, 1, time.Now().Add(24*time.Hour), 1.0, 100.0, 10.0, 110.0, false)
	assert.NoError(t, err)

	orderService.DeleteOrder(ctx, orderID)

}

// TestHandleExpiredOrder проверяет обработку просроченного заказа и отправку сообщения в Kafka
func TestHandleExpiredOrder(t *testing.T) {
	brokers := []string{"localhost:9092"}
	producer, err := kafka.NewProducer(brokers, "test-topic")
	assert.NoError(t, err)
	defer producer.Close()

	wp := pool.NewWorkerPool(2)
	orderService := service.NewOrderService(testDB, wp, producer)

	ctx := context.Background()
	orderID, err := orderService.CreateOrder(ctx, 1, 1, 1, time.Now().Add(-24*time.Hour), 1.0, 100.0, 10.0, 110.0, false) // Устанавливаем прошедшую дату
	assert.NoError(t, err)

	assert.NotZero(t, orderID, "Order ID should not be zero")

	err = orderService.CheckExpiredOrders(ctx)
	assert.NoError(t, err)

}

// TestCreateReturn проверяет создание возврата и отправку сообщения в Kafka
func TestCreateReturn(t *testing.T) {
	brokers := []string{"localhost:9092"}
	producer, err := kafka.NewProducer(brokers, "test-topic")
	assert.NoError(t, err)
	defer producer.Close()

	wp := pool.NewWorkerPool(2)
	returnService := service.NewReturnService(testDB, wp, producer)

	now := time.Now()

	order := model.Order{
		UserID:         1,
		PackagingID:    1,
		StatusID:       1,
		AcceptanceDate: now.Add(-1 * time.Hour),
		ExpirationDate: now.Add(1 * time.Hour),
		Weight:         1.0,
		BaseCost:       100.0,
		PackagingCost:  10.0,
		TotalCost:      110.0,
		WithFilm:       false,
		IssueDate:      now,
	}

	orderID, err := dao.CreateOrder(context.Background(), order, testDB)
	assert.NoError(t, err)

	err = returnService.CreateReturn(context.Background(), orderID)
	assert.NoError(t, err, "Ошибка при создании возврата")

}

// TestUpdateReturn проверяет обновление возврата и отправку сообщения в Kafka
func TestUpdateReturn(t *testing.T) {
	brokers := []string{"localhost:9092"}
	producer, err := kafka.NewProducer(brokers, "test-topic")
	assert.NoError(t, err)
	defer producer.Close()

	wp := pool.NewWorkerPool(2)
	returnService := service.NewReturnService(testDB, wp, producer)

	now := time.Now()

	order := model.Order{
		UserID:         1,
		PackagingID:    1,
		StatusID:       1,
		AcceptanceDate: now.Add(-1 * time.Hour),
		ExpirationDate: now.Add(1 * time.Hour),
		Weight:         1.0,
		BaseCost:       100.0,
		PackagingCost:  10.0,
		TotalCost:      110.0,
		WithFilm:       false,
		IssueDate:      now,
	}
	orderID, err := dao.CreateOrder(context.Background(), order, testDB)
	assert.NoError(t, err)

	err = returnService.CreateReturn(context.Background(), orderID)
	assert.NoError(t, err)

	returns, err := returnService.GetReturns(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, returns)

	returnID := returns[0].ReturnID

	err = returnService.UpdateReturn(context.Background(), returnID, orderID, 1, 1, 100.0, 10.0, 110.0, 1, 1)
	assert.NoError(t, err, "Ошибка при обновлении возврата")
}
