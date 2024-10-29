package dao_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"homework1/internal/dao"
	"homework1/internal/model"
	"testing"
)

// Test CreateReturnReason function
func TestCreateReturnReason(t *testing.T) {
	ctx := context.Background()

	reason := model.ReturnReason{
		Reason: "Поврежденный товар",
	}

	// Вставка новой причины возврата
	reasonID, err := dao.CreateReturnReason(ctx, reason, testDB)
	assert.NoError(t, err, "ошибка при создании причины возврата")
	assert.NotZero(t, reasonID, "ID причины возврата должен быть больше нуля")
}

// Test GetReturnReasonByID function
func TestGetReturnReasonByID(t *testing.T) {
	ctx := context.Background()

	// Создание тестовой причины возврата
	reason := model.ReturnReason{
		Reason: "Неправильный товар",
	}
	reasonID, err := dao.CreateReturnReason(ctx, reason, testDB)
	assert.NoError(t, err, "ошибка при создании причины возврата")

	// Получение причины возврата по ID
	fetchedReason, err := dao.GetReturnReasonByID(ctx, reasonID, testDB)
	assert.NoError(t, err, "ошибка при получении причины возврата по ID")
	assert.Equal(t, reasonID, fetchedReason.ReasonID, "ID причины возврата должен совпадать")
	assert.Equal(t, reason.Reason, fetchedReason.Reason, "Текст причины возврата должен совпадать")
}

// Test UpdateReturnReason function
func TestUpdateReturnReason(t *testing.T) {
	ctx := context.Background()

	// Создание тестовой причины возврата
	reason := model.ReturnReason{
		Reason: "Неверный заказ",
	}
	reasonID, err := dao.CreateReturnReason(ctx, reason, testDB)
	assert.NoError(t, err, "ошибка при создании причины возврата")

	// Обновление причины возврата
	reason.ReasonID = reasonID
	reason.Reason = "Товар не подошел"
	err = dao.UpdateReturnReason(ctx, reason, testDB)
	assert.NoError(t, err, "ошибка при обновлении причины возврата")

	// Получение обновленной причины возврата
	updatedReason, err := dao.GetReturnReasonByID(ctx, reasonID, testDB)
	assert.NoError(t, err, "ошибка при получении обновленной причины возврата")
	assert.Equal(t, "Товар не подошел", updatedReason.Reason, "Текст причины возврата должен был измениться")
}

// Test DeleteReturnReason function
func TestDeleteReturnReason(t *testing.T) {
	ctx := context.Background()

	// Создание тестовой причины возврата
	reason := model.ReturnReason{
		Reason: "Упаковка повреждена",
	}
	reasonID, err := dao.CreateReturnReason(ctx, reason, testDB)
	assert.NoError(t, err, "ошибка при создании причины возврата")

	// Удаление причины возврата
	err = dao.DeleteReturnReason(ctx, reasonID, testDB)
	assert.NoError(t, err, "ошибка при удалении причины возврата")

	// Проверка, что причина возврата удалена
	fetchedReason, err := dao.GetReturnReasonByID(ctx, reasonID, testDB)
	assert.Error(t, err, "должна возникнуть ошибка при попытке получить удаленную причину возврата")
	assert.Nil(t, fetchedReason, "причина возврата должна быть nil после удаления")
}

// Test GetAllReturnReasons function
func TestGetAllReturnReasons(t *testing.T) {
	ctx := context.Background()

	// Создание нескольких причин возврата
	reason1 := model.ReturnReason{
		Reason: "Не понравился товар",
	}
	reason2 := model.ReturnReason{
		Reason: "Дефектный товар",
	}
	_, err := dao.CreateReturnReason(ctx, reason1, testDB)
	assert.NoError(t, err, "ошибка при создании причины возврата 1")

	_, err = dao.CreateReturnReason(ctx, reason2, testDB)
	assert.NoError(t, err, "ошибка при создании причины возврата 2")

	// Получение всех причин возврата
	returnReasons, err := dao.GetAllReturnReasons(ctx, testDB)
	assert.NoError(t, err, "ошибка при получении всех причин возвратов")
	assert.GreaterOrEqual(t, len(returnReasons), 2, "Должно быть как минимум 2 причины возврата")
}
