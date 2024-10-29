package dao_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"homework1/internal/dao"
	"homework1/internal/model"
	"testing"
	"time"
)

// Тест создания возврата
func TestCreateReturn(t *testing.T) {
	ctx := context.Background()

	newReturn := model.Return{
		OrderID:       1,
		UserID:        1,
		ReturnDate:    time.Now(),
		ReasonID:      1,
		BaseCost:      10.50,
		PackagingCost: 2.50,
		PackagingID:   1,
		TotalCost:     13.00,
		StatusID:      1,
	}

	err := dao.CreateReturn(ctx, newReturn, testDB)
	assert.NoError(t, err, "ошибка при создании возврата")
}

// Тест чтения всех возвратов
func TestReadReturns(t *testing.T) {
	ctx := context.Background()

	returns, err := dao.ReadReturns(ctx, testDB)
	assert.NoError(t, err, "ошибка при чтении возвратов")
	assert.GreaterOrEqual(t, len(returns), 1, "должен быть хотя бы один возврат")
}

// Тест поиска возврата по ID заказа
func TestFindReturnByOrderID(t *testing.T) {
	ctx := context.Background()

	// Поиск возврата по ID заказа
	orderID := 1
	ret, err := dao.FindReturnByOrderID(ctx, orderID, testDB)
	assert.NoError(t, err, "ошибка при поиске возврата по ID заказа")
	assert.NotNil(t, ret, "возврат не должен быть nil")
	assert.Equal(t, orderID, ret.OrderID, "ID заказа должен совпадать")
}

// Тест обновления возврата
func TestUpdateReturn(t *testing.T) {
	ctx := context.Background()

	// Поиск существующего возврата
	returns, err := dao.ReadReturns(ctx, testDB)
	assert.NoError(t, err, "ошибка при чтении возвратов")
	assert.Greater(t, len(returns), 0, "должен быть хотя бы один возврат")

	// Обновление возврата
	ret := returns[0]
	ret.BaseCost = 12.00
	err = dao.UpdateReturn(ctx, ret, testDB)
	assert.NoError(t, err, "ошибка при обновлении возврата")
}

// Тест удаления возврата
func TestDeleteReturn(t *testing.T) {
	ctx := context.Background()

	// Поиск существующего возврата
	returns, err := dao.ReadReturns(ctx, testDB)
	assert.NoError(t, err, "ошибка при чтении возвратов")
	assert.Greater(t, len(returns), 0, "должен быть хотя бы один возврат")

	// Удаление возврата
	err = dao.DeleteReturn(ctx, returns[0].ReturnID, testDB)
	assert.NoError(t, err, "ошибка при удалении возврата")
}

// Тест поиска возвратов по ID пользователя
func TestFindReturnsByUserID(t *testing.T) {
	// Создаем тестовый возврат для пользователя с ID 1
	ret := model.Return{
		OrderID:       1, // используйте правильный ID существующего заказа
		UserID:        1, // используйте правильный ID пользователя
		ReturnDate:    time.Now(),
		ReasonID:      1,
		BaseCost:      100,
		PackagingCost: 10,
		TotalCost:     110,
		StatusID:      1,
	}

	// Создаем возврат в базе данных
	err := dao.CreateReturn(context.Background(), ret, testDB)
	if err != nil {
		t.Fatalf("Не удалось создать возврат: %v", err)
	}

	// Теперь проверяем, что возврат успешно сохранен и его можно найти
	returns, err := dao.FindReturnsByUserID(context.Background(), ret.UserID, testDB)
	if err != nil {
		t.Fatalf("Ошибка поиска возвратов по ID пользователя: %v", err)
	}

	// Проверяем, что найден хотя бы один возврат
	if len(returns) == 0 {
		t.Errorf("должен быть хотя бы один возврат для пользователя")
	}
}

// Тест проверки существования возврата
func TestCheckReturnExists(t *testing.T) {
	// Создаем тестовый возврат
	ret := model.Return{
		OrderID:       1, // Используйте существующий OrderID
		UserID:        1, // Используйте существующий UserID
		ReturnDate:    time.Now(),
		ReasonID:      1,
		BaseCost:      100,
		PackagingCost: 10,
		TotalCost:     110,
		StatusID:      1,
	}

	// Создаем возврат и сохраняем его идентификатор
	err := dao.CreateReturn(context.Background(), ret, testDB)
	if err != nil {
		t.Fatalf("Не удалось создать возврат: %v", err)
	}

	// Получаем ID созданного возврата
	returns, err := dao.ReadReturns(context.Background(), testDB)
	if err != nil {
		t.Fatalf("Ошибка при получении возвратов: %v", err)
	}
	if len(returns) == 0 {
		t.Fatalf("Не удалось получить список возвратов")
	}
	createdReturnID := returns[0].ReturnID

	// Проверяем, существует ли возврат
	exists, err := dao.CheckReturnExists(context.Background(), createdReturnID, testDB)
	if err != nil {
		t.Fatalf("Ошибка при проверке существования возврата: %v", err)
	}

	if !exists {
		t.Errorf("должен быть хотя бы один возврат")
	}
}
