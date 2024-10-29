package dao_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"homework1/internal/dao"
	"homework1/internal/model"
	"testing"
)

// Тест создания статуса
func TestCreateStatus(t *testing.T) {
	ctx := context.Background()

	newStatus := model.Status{
		StatusName: "Создано",
	}

	statusID, err := dao.CreateStatus(ctx, newStatus, testDB)
	assert.NoError(t, err, "ошибка при создании статуса")
	assert.Greater(t, statusID, 0, "ID статуса должен быть больше 0")
}

// Тест получения статуса по ID
func TestGetStatusByID(t *testing.T) {
	ctx := context.Background()

	// Создание нового статуса
	newStatus := model.Status{
		StatusName: "В процессе",
	}
	statusID, err := dao.CreateStatus(ctx, newStatus, testDB)
	assert.NoError(t, err, "ошибка при создании статуса")

	// Получение статуса по ID
	status, err := dao.GetStatusByID(ctx, statusID, testDB)
	assert.NoError(t, err, "ошибка при получении статуса")
	assert.NotNil(t, status, "статус не должен быть nil")
	assert.Equal(t, newStatus.StatusName, status.StatusName, "имя статуса должно совпадать")
}

// Тест обновления статуса
func TestUpdateStatus(t *testing.T) {
	ctx := context.Background()

	// Создание нового статуса
	newStatus := model.Status{
		StatusName: "Завершено",
	}
	statusID, err := dao.CreateStatus(ctx, newStatus, testDB)
	assert.NoError(t, err, "ошибка при создании статуса")

	// Обновление статуса
	updatedStatus := model.Status{
		StatusID:   statusID,
		StatusName: "Завершено успешно",
	}
	err = dao.UpdateStatus(ctx, updatedStatus, testDB)
	assert.NoError(t, err, "ошибка при обновлении статуса")

	// Проверка обновленного статуса
	status, err := dao.GetStatusByID(ctx, statusID, testDB)
	assert.NoError(t, err, "ошибка при получении статуса")
	assert.Equal(t, updatedStatus.StatusName, status.StatusName, "имя статуса должно быть обновлено")
}

// Тест удаления статуса
func TestDeleteStatus(t *testing.T) {
	ctx := context.Background()

	// Создание нового статуса
	newStatus := model.Status{
		StatusName: "Удалить",
	}
	statusID, err := dao.CreateStatus(ctx, newStatus, testDB)
	assert.NoError(t, err, "ошибка при создании статуса")

	// Удаление статуса
	err = dao.DeleteStatus(ctx, statusID, testDB)
	assert.NoError(t, err, "ошибка при удалении статуса")

	// Проверка, что статус удален
	exists, err := dao.CheckStatusExists(ctx, statusID, testDB)
	assert.NoError(t, err, "ошибка при проверке существования статуса")
	assert.False(t, exists, "статус не должен существовать")
}

// Тест получения всех статусов
func TestGetAllStatuses(t *testing.T) {
	ctx := context.Background()

	// Создание нескольких статусов
	_, err := dao.CreateStatus(ctx, model.Status{StatusName: "Ожидание"}, testDB)
	assert.NoError(t, err, "ошибка при создании статуса 'Ожидание'")
	_, err = dao.CreateStatus(ctx, model.Status{StatusName: "Доставка"}, testDB)
	assert.NoError(t, err, "ошибка при создании статуса 'Доставка'")

	// Получение всех статусов
	statuses, err := dao.GetAllStatuses(ctx, testDB)
	assert.NoError(t, err, "ошибка при получении всех статусов")
	assert.GreaterOrEqual(t, len(statuses), 2, "должно быть хотя бы два статуса")
}

// Тест получения статуса по имени
func TestGetStatusByName(t *testing.T) {
	ctx := context.Background()

	// Создание нового статуса
	statusName := "Ожидание"
	newStatus := model.Status{
		StatusName: statusName,
	}
	_, err := dao.CreateStatus(ctx, newStatus, testDB)
	assert.NoError(t, err, "ошибка при создании статуса")

	// Получение статуса по имени
	status, err := dao.GetStatusByName(ctx, statusName, testDB)
	assert.NoError(t, err, "ошибка при получении статуса по имени")
	assert.NotNil(t, status, "статус не должен быть nil")
	assert.Equal(t, statusName, status.StatusName, "имя статуса должно совпадать")
}

// Тест проверки существования статуса
func TestCheckStatusExists(t *testing.T) {
	ctx := context.Background()

	// Создание нового статуса
	newStatus := model.Status{
		StatusName: "Для проверки",
	}
	statusID, err := dao.CreateStatus(ctx, newStatus, testDB)
	assert.NoError(t, err, "ошибка при создании статуса")

	// Проверка существования статуса
	exists, err := dao.CheckStatusExists(ctx, statusID, testDB)
	assert.NoError(t, err, "ошибка при проверке существования статуса")
	assert.True(t, exists, "статус должен существовать")
}
