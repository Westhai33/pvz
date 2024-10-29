package controller

import (
	"context"
	"fmt"
	"homework1/internal/model"
	"homework1/internal/service"
	"strconv"
)

// CreateStatus создает новый статус
func CreateStatus(ctx context.Context, statusService *service.StatusService, statusName string) (int, error) {
	statusID, err := statusService.CreateStatus(ctx, statusName)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания статуса: %w", err)
	}

	return statusID, nil
}

// GetStatusByID возвращает статус по его ID
func GetStatusByID(ctx context.Context, statusService *service.StatusService, statusID int) (*model.Status, error) {
	status, err := statusService.GetStatusByID(ctx, statusID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения статуса с ID %d: %w", statusID, err)
	}

	return status, nil
}

// GetAllStatuses возвращает все статусы
func GetAllStatuses(ctx context.Context, statusService *service.StatusService) ([]model.Status, error) {
	statuses, err := statusService.GetAllStatuses(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка статусов: %w", err)
	}

	return statuses, nil
}

// UpdateStatus обновляет существующий статус
func UpdateStatus(ctx context.Context, statusService *service.StatusService, statusIDStr, statusName string) error {
	statusID, err := strconv.Atoi(statusIDStr)
	if err != nil {
		return fmt.Errorf("неверный формат идентификатора статуса")
	}

	err = statusService.UpdateStatus(ctx, statusID, statusName)
	if err != nil {
		return fmt.Errorf("ошибка обновления статуса с ID %d: %w", statusID, err)
	}

	return nil
}

// DeleteStatus удаляет статус по его ID
func DeleteStatus(ctx context.Context, statusService *service.StatusService, statusIDStr string) error {
	statusID, err := strconv.Atoi(statusIDStr)
	if err != nil {
		return fmt.Errorf("неверный формат идентификатора статуса")
	}

	err = statusService.DeleteStatus(ctx, statusID)
	if err != nil {
		return fmt.Errorf("ошибка удаления статуса с ID %d: %w", statusID, err)
	}

	return nil
}

// GetStatusNameByID возвращает имя статуса по его ID
func GetStatusNameByID(ctx context.Context, statusService *service.StatusService, statusID int) (string, error) {
	return statusService.GetStatusNameByID(ctx, statusID)
}
