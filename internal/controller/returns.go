package controller

import (
	"context"
	"fmt"
	"homework1/internal/metrics"
	"homework1/internal/model"
	"homework1/internal/service"
	"strconv"
)

func CreateReturn(ctx context.Context, returnService *service.ReturnService, orderIDStr string) error {
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		return fmt.Errorf("неверный формат идентификатора заказа")
	}

	err = returnService.CreateReturn(ctx, orderID)
	if err != nil {
		return fmt.Errorf("ошибка создания возврата: %w", err)
	}

	metrics.IncrementCreatedReturns("created")

	return nil
}

// GetReturns возвращает все возвраты
func GetReturns(ctx context.Context, returnService *service.ReturnService) ([]model.Return, error) {
	return returnService.GetReturns(ctx)
}

// GetReturnByOrderID возвращает возврат по идентификатору заказа
func GetReturnByOrderID(ctx context.Context, returnService *service.ReturnService, orderIDStr string) (*model.Return, error) {
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат идентификатора заказа")
	}

	ret, err := returnService.GetReturnByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения возврата для заказа с ID %d: %w", orderID, err)
	}

	return ret, nil
}

// DeleteReturn удаляет возврат
func DeleteReturn(ctx context.Context, returnService *service.ReturnService, returnIDStr string) error {
	returnID, err := strconv.Atoi(returnIDStr)
	if err != nil {
		return fmt.Errorf("неверный формат идентификатора возврата")
	}

	err = returnService.DeleteReturn(ctx, returnID)
	if err != nil {
		return fmt.Errorf("ошибка удаления возврата с ID %d: %w", returnID, err)
	}

	return nil
}

// GetReturnsByUserID возвращает все возвраты для пользователя
func GetReturnsByUserID(ctx context.Context, returnService *service.ReturnService, userIDStr string) ([]model.Return, error) {
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат идентификатора пользователя")
	}

	return returnService.GetReturnsByUserID(ctx, userID)
}

// ProcessReturn обрабатывает возврат для заказа
func ProcessReturn(ctx context.Context, returnService *service.ReturnService, statusService *service.StatusService, orderID int) error {
	returnRecord, err := returnService.GetReturnByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("ошибка поиска возврата для заказа с ID %d: %w", orderID, err)
	}

	if returnRecord == nil {
		return fmt.Errorf("возврат для заказа с ID %d не найден", orderID)
	}

	statusID, err := statusService.GetStatusByName(ctx, "Передан курьеру")
	if err != nil {
		return fmt.Errorf("ошибка получения статуса 'Передан курьеру': %w", err)
	}

	returnRecord.StatusID = statusID

	if err := returnService.UpdateReturn(ctx, returnRecord.ReturnID, returnRecord.OrderID, returnRecord.UserID,
		returnRecord.ReasonID, returnRecord.BaseCost, returnRecord.PackagingCost,
		returnRecord.TotalCost, returnRecord.PackagingID, returnRecord.StatusID); err != nil {
		return fmt.Errorf("ошибка обновления статуса возврата с ID %d: %w", returnRecord.ReturnID, err)
	}

	return nil
}
