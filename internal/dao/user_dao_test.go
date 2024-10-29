package dao_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"homework1/internal/dao"
	"homework1/internal/model"
	"testing"
	"time"
)

// Тест создания пользователя
func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	newUser := model.User{
		Username:  "Иван",
		CreatedAt: time.Now(),
	}

	userID, err := dao.CreateUser(ctx, newUser, testDB)
	assert.NoError(t, err, "ошибка при создании пользователя")
	assert.Greater(t, userID, 0, "ID пользователя должен быть больше 0")
}

// Тест получения пользователя по ID
func TestGetUserByID(t *testing.T) {
	ctx := context.Background()

	// Создание нового пользователя
	newUser := model.User{
		Username:  "Алексей",
		CreatedAt: time.Now(),
	}
	userID, err := dao.CreateUser(ctx, newUser, testDB)
	assert.NoError(t, err, "ошибка при создании пользователя")

	// Получение пользователя по ID
	user, err := dao.GetUserByID(ctx, userID, testDB)
	assert.NoError(t, err, "ошибка при получении пользователя")
	assert.NotNil(t, user, "пользователь не должен быть nil")
	assert.Equal(t, newUser.Username, user.Username, "имя пользователя должно совпадать")
}

// Тест обновления пользователя
func TestUpdateUser(t *testing.T) {
	ctx := context.Background()

	// Создание нового пользователя
	newUser := model.User{
		Username:  "Обновляемый пользователь",
		CreatedAt: time.Now(),
	}
	userID, err := dao.CreateUser(ctx, newUser, testDB)
	assert.NoError(t, err, "ошибка при создании пользователя")

	// Обновление пользователя
	updatedUser := model.User{
		UserID:    userID,
		Username:  "Обновлено",
		CreatedAt: time.Now(),
	}
	err = dao.UpdateUser(ctx, updatedUser, testDB)
	assert.NoError(t, err, "ошибка при обновлении пользователя")

	// Проверка обновленного пользователя
	user, err := dao.GetUserByID(ctx, userID, testDB)
	assert.NoError(t, err, "ошибка при получении пользователя")
	assert.Equal(t, updatedUser.Username, user.Username, "имя пользователя должно быть обновлено")
}

// Тест удаления пользователя
func TestDeleteUser(t *testing.T) {
	ctx := context.Background()

	// Создание нового пользователя
	newUser := model.User{
		Username:  "Удаляемый пользователь",
		CreatedAt: time.Now(),
	}
	userID, err := dao.CreateUser(ctx, newUser, testDB)
	assert.NoError(t, err, "ошибка при создании пользователя")

	// Удаление пользователя
	err = dao.DeleteUser(ctx, userID, testDB)
	assert.NoError(t, err, "ошибка при удалении пользователя")

	// Проверка, что пользователь удален
	exists, err := dao.CheckUserExists(ctx, userID, testDB)
	assert.NoError(t, err, "ошибка при проверке существования пользователя")
	assert.False(t, exists, "пользователь не должен существовать")
}

// Тест получения всех пользователей
func TestGetAllUsers(t *testing.T) {
	ctx := context.Background()

	// Создание нескольких пользователей
	_, err := dao.CreateUser(ctx, model.User{Username: "Пользователь1", CreatedAt: time.Now()}, testDB)
	assert.NoError(t, err, "ошибка при создании пользователя 'Пользователь1'")
	_, err = dao.CreateUser(ctx, model.User{Username: "Пользователь2", CreatedAt: time.Now()}, testDB)
	assert.NoError(t, err, "ошибка при создании пользователя 'Пользователь2'")

	// Получение всех пользователей
	users, err := dao.GetAllUsers(ctx, testDB)
	assert.NoError(t, err, "ошибка при получении всех пользователей")
	assert.GreaterOrEqual(t, len(users), 2, "должно быть хотя бы два пользователя")
}

// Тест получения имени пользователя по ID
func TestGetUserNameByID(t *testing.T) {
	ctx := context.Background()

	// Создание нового пользователя
	newUser := model.User{
		Username:  "Иван",
		CreatedAt: time.Now(),
	}
	userID, err := dao.CreateUser(ctx, newUser, testDB)
	assert.NoError(t, err, "ошибка при создании пользователя")

	// Получение имени пользователя по ID
	username, err := dao.GetUserNameByID(ctx, userID, testDB)
	assert.NoError(t, err, "ошибка при получении имени пользователя")
	assert.Equal(t, newUser.Username, username, "имя пользователя должно совпадать")
}

// Тест проверки существования пользователя
func TestCheckUserExists(t *testing.T) {
	ctx := context.Background()

	// Создание нового пользователя
	newUser := model.User{
		Username:  "Пользователь для проверки",
		CreatedAt: time.Now(),
	}
	userID, err := dao.CreateUser(ctx, newUser, testDB)
	assert.NoError(t, err, "ошибка при создании пользователя")

	// Проверка существования пользователя
	exists, err := dao.CheckUserExists(ctx, userID, testDB)
	assert.NoError(t, err, "ошибка при проверке существования пользователя")
	assert.True(t, exists, "пользователь должен существовать")
}
