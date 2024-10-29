package controller

import (
	"context"
	"fmt"
	"homework1/internal/model"
	"homework1/internal/service"
)

// CreateUser создает нового пользователя
func CreateUser(ctx context.Context, userService *service.UserService, username string) (int, error) {
	userID, err := userService.CreateUser(ctx, username)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	return userID, nil
}

// GetUserByID возвращает пользователя по его ID
func GetUserByID(ctx context.Context, userService *service.UserService, userID int) (*model.User, error) {
	user, err := userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователя с ID %d: %w", userID, err)
	}
	return user, nil
}

// GetAllUsers возвращает всех пользователей
func GetAllUsers(ctx context.Context, userService *service.UserService) ([]model.User, error) {
	users, err := userService.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка пользователей: %w", err)
	}

	return users, nil
}

// UpdateUser обновляет существующего пользователя
func UpdateUser(ctx context.Context, userService *service.UserService, userID int, username string) error {
	err := userService.UpdateUser(ctx, userID, username)
	if err != nil {
		return fmt.Errorf("ошибка обновления пользователя с ID %d: %w", userID, err)
	}

	return nil
}

// DeleteUser удаляет пользователя по его ID
func DeleteUser(ctx context.Context, userService *service.UserService, userID int) error {
	err := userService.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("ошибка удаления пользователя с ID %d: %w", userID, err)
	}

	return nil
}

// GetUserNameByID возвращает имя пользователя по его ID
func GetUserNameByID(ctx context.Context, userService *service.UserService, userID int) (string, error) {
	return userService.GetUserNameByID(ctx, userID)
}
