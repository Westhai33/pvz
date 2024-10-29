package controller

import (
	"context"
	"fmt"
	"homework1/internal/model"
	"homework1/internal/service"
	"strconv"
)

// CreateReturnReason создает новую причину возврата
func CreateReturnReason(ctx context.Context, returnReasonService *service.ReturnReasonService, reason string) (int, error) {
	reasonID, err := returnReasonService.CreateReturnReason(ctx, reason)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания причины возврата: %w", err)
	}

	return reasonID, nil
}

// GetReturnReasonByID возвращает причину возврата по её ID
func GetReturnReasonByID(ctx context.Context, returnReasonService *service.ReturnReasonService, reasonID int) (*model.ReturnReason, error) {
	reason, err := returnReasonService.GetReturnReasonByID(ctx, reasonID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения причины возврата с ID %d: %w", reasonID, err)
	}
	return reason, nil
}

// GetAllReturnReasons возвращает все причины возвратов
func GetAllReturnReasons(ctx context.Context, returnReasonService *service.ReturnReasonService) ([]model.ReturnReason, error) {
	reasons, err := returnReasonService.GetAllReturnReasons(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка причин возвратов: %w", err)
	}

	return reasons, nil
}

// UpdateReturnReason обновляет существующую причину возврата
func UpdateReturnReason(ctx context.Context, returnReasonService *service.ReturnReasonService, reasonIDStr, reason string) error {
	reasonID, err := strconv.Atoi(reasonIDStr)
	if err != nil {
		return fmt.Errorf("неверный формат идентификатора причины возврата")
	}

	err = returnReasonService.UpdateReturnReason(ctx, reasonID, reason)
	if err != nil {
		return fmt.Errorf("ошибка обновления причины возврата с ID %d: %w", reasonID, err)
	}

	return nil
}

// DeleteReturnReason удаляет причину возврата по её ID
func DeleteReturnReason(ctx context.Context, returnReasonService *service.ReturnReasonService, reasonIDStr string) error {
	reasonID, err := strconv.Atoi(reasonIDStr)
	if err != nil {
		return fmt.Errorf("неверный формат идентификатора причины возврата")
	}

	err = returnReasonService.DeleteReturnReason(ctx, reasonID)
	if err != nil {
		return fmt.Errorf("ошибка удаления причины возврата с ID %d: %w", reasonID, err)
	}

	return nil
}
