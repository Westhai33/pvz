package dao_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"homework1/internal/dao"
	"homework1/internal/model"
)

// Test CreatePackaging function
func TestCreatePackaging(t *testing.T) {
	ctx := context.Background()

	packaging := model.PackagingOption{
		Type:      "Коробка",
		Cost:      10.5,
		MaxWeight: 20,
	}

	// Вставка новой упаковки
	packagingID, err := dao.CreatePackaging(ctx, packaging, testDB)
	assert.NoError(t, err, "ошибка при создании упаковки")
	assert.NotZero(t, packagingID, "ID упаковки должен быть больше нуля")
}

// Test GetAllPackaging function
func TestGetAllPackaging(t *testing.T) {
	ctx := context.Background()

	// Получение всех упаковок
	packagingOptions, err := dao.GetAllPackaging(ctx, testDB)
	assert.NoError(t, err, "ошибка при получении упаковок")
	assert.NotEmpty(t, packagingOptions, "Список упаковок не должен быть пустым")
}

// Test GetPackagingByID function
func TestGetPackagingByID(t *testing.T) {
	ctx := context.Background()

	// Создание тестовой упаковки
	packaging := model.PackagingOption{
		Type:      "Коробка",
		Cost:      12.0,
		MaxWeight: 15,
	}
	packagingID, err := dao.CreatePackaging(ctx, packaging, testDB)
	assert.NoError(t, err, "ошибка при создании упаковки")

	// Получение упаковки по ID
	fetchedPackaging, err := dao.GetPackagingByID(ctx, packagingID, testDB)
	assert.NoError(t, err, "ошибка при получении упаковки по ID")
	assert.Equal(t, packagingID, fetchedPackaging.PackagingID, "ID упаковки должен совпадать")
	assert.Equal(t, packaging.Type, fetchedPackaging.Type, "Тип упаковки должен совпадать")
}

// Test UpdatePackaging function
func TestUpdatePackaging(t *testing.T) {
	ctx := context.Background()

	// Создание тестовой упаковки
	packaging := model.PackagingOption{
		Type:      "Пакет",
		Cost:      5.0,
		MaxWeight: 10,
	}
	packagingID, err := dao.CreatePackaging(ctx, packaging, testDB)
	assert.NoError(t, err, "ошибка при создании упаковки")

	// Обновление упаковки
	packaging.PackagingID = packagingID
	packaging.Type = "Большой Пакет"
	packaging.Cost = 7.0
	err = dao.UpdatePackaging(ctx, packaging, testDB)
	assert.NoError(t, err, "ошибка при обновлении упаковки")

	// Получение обновленной упаковки
	updatedPackaging, err := dao.GetPackagingByID(ctx, packagingID, testDB)
	assert.NoError(t, err, "ошибка при получении обновленной упаковки")
	assert.Equal(t, "Большой Пакет", updatedPackaging.Type, "Тип упаковки должен был измениться")
	assert.Equal(t, 7.0, updatedPackaging.Cost, "Стоимость упаковки должна была измениться")
}

// Test DeletePackaging function
func TestDeletePackaging(t *testing.T) {
	ctx := context.Background()

	// Создание тестовой упаковки
	packaging := model.PackagingOption{
		Type:      "Бумажный пакет",
		Cost:      3.0,
		MaxWeight: 5,
	}
	packagingID, err := dao.CreatePackaging(ctx, packaging, testDB)
	assert.NoError(t, err, "ошибка при создании упаковки")

	// Удаление упаковки
	err = dao.DeletePackaging(ctx, packagingID, testDB)
	assert.NoError(t, err, "ошибка при удалении упаковки")

	// Проверка, что упаковка удалена
	fetchedPackaging, err := dao.GetPackagingByID(ctx, packagingID, testDB)
	assert.Error(t, err, "должна возникнуть ошибка при попытке получить удаленную упаковку")
	assert.Nil(t, fetchedPackaging, "упаковка должна быть nil после удаления")
}
