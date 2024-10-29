package dao_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/assert"
	"homework1/internal/dao"
	"homework1/internal/model"
)

var testDB *pgxpool.Pool

// Setup function to initialize test database connection.
func TestMain(m *testing.M) {
	// Подключение к тестовой базе данных
	dbUrl := "postgres://postgres:postgres@localhost:5433/test_db?sslmode=disable"
	var err error
	testDB, err = pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		panic("Не удалось подключиться к тестовой базе данных: " + err.Error())
	}
	defer testDB.Close()

	// Запуск всех тестов
	os.Exit(m.Run())
}

// Test CreateOrder function
func TestCreateOrder(t *testing.T) {
	ctx := context.Background()

	order := model.Order{
		UserID:         1,
		AcceptanceDate: time.Now(),
		ExpirationDate: time.Now().Add(24 * time.Hour),
		Weight:         5.0,
		BaseCost:       100.0,
		PackagingCost:  10.0,
		TotalCost:      110.0,
		PackagingID:    1,
		StatusID:       1,
		WithFilm:       true,
	}

	// Вставка нового заказа
	orderID, err := dao.CreateOrder(ctx, order, testDB)
	assert.NoError(t, err, "ошибка при создании заказа")
	assert.NotZero(t, orderID, "ID заказа должен быть больше нуля")
}

// Test GetOrderByID function
func TestGetOrderByID(t *testing.T) {
	ctx := context.Background()

	// Создание тестового заказа
	order := model.Order{
		UserID:         2,
		AcceptanceDate: time.Now(),
		ExpirationDate: time.Now().Add(24 * time.Hour),
		Weight:         10.0,
		BaseCost:       150.0,
		PackagingCost:  20.0,
		TotalCost:      170.0,
		PackagingID:    2,
		StatusID:       2,
		WithFilm:       false,
	}
	orderID, err := dao.CreateOrder(ctx, order, testDB)
	assert.NoError(t, err, "ошибка при создании заказа")

	// Получение заказа по ID
	fetchedOrder, err := dao.GetOrderByID(ctx, orderID, testDB)
	assert.NoError(t, err, "ошибка при получении заказа по ID")
	assert.Equal(t, orderID, fetchedOrder.OrderID, "ID заказа должен совпадать")
	assert.Equal(t, order.UserID, fetchedOrder.UserID, "ID пользователя должен совпадать")
}

// Test UpdateOrder function
func TestUpdateOrder(t *testing.T) {
	ctx := context.Background()

	// Создание тестового заказа
	order := model.Order{
		UserID:         3,
		AcceptanceDate: time.Now(),
		ExpirationDate: time.Now().Add(24 * time.Hour),
		Weight:         8.0,
		BaseCost:       200.0,
		PackagingCost:  30.0,
		TotalCost:      230.0,
		PackagingID:    3,
		StatusID:       1,
		WithFilm:       false,
	}
	orderID, err := dao.CreateOrder(ctx, order, testDB)
	assert.NoError(t, err, "ошибка при создании заказа")

	// Обновление заказа
	order.OrderID = orderID
	order.Weight = 12.0
	order.TotalCost = 260.0
	err = dao.UpdateOrder(ctx, order, testDB)
	assert.NoError(t, err, "ошибка при обновлении заказа")

	// Получение обновленного заказа
	updatedOrder, err := dao.GetOrderByID(ctx, orderID, testDB)
	assert.NoError(t, err, "ошибка при получении обновленного заказа")
	assert.Equal(t, 12.0, updatedOrder.Weight, "Вес должен был измениться")
	assert.Equal(t, 260.0, updatedOrder.TotalCost, "Общая стоимость должна была измениться")
}

// Test DeleteOrder function
func TestDeleteOrder(t *testing.T) {
	ctx := context.Background()

	// Создание тестового заказа
	order := model.Order{
		UserID:         4,
		AcceptanceDate: time.Now(),
		ExpirationDate: time.Now().Add(24 * time.Hour),
		Weight:         15.0,
		BaseCost:       300.0,
		PackagingCost:  40.0,
		TotalCost:      340.0,
		PackagingID:    4,
		StatusID:       1,
		WithFilm:       true,
	}
	orderID, err := dao.CreateOrder(ctx, order, testDB)
	assert.NoError(t, err, "ошибка при создании заказа")

	// Удаление заказа
	err = dao.DeleteOrder(ctx, orderID, testDB)
	assert.NoError(t, err, "ошибка при удалении заказа")

	// Проверка, что заказ удален
	fetchedOrder, err := dao.GetOrderByID(ctx, orderID, testDB)
	assert.Error(t, err, "должна возникнуть ошибка при попытке получить удаленный заказ")
	assert.Nil(t, fetchedOrder, "заказ должен быть nil после удаления")
}
