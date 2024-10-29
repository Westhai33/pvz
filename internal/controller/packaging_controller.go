package controller

import (
	"context"
	"fmt"
	"homework1/internal/model"
	"homework1/internal/service"
	"strconv"
)

// CreatePackaging обрабатывает создание новой упаковки
func CreatePackaging(ctx context.Context, packagingService *service.PackagingService, packagingType, costStr, maxWeightStr string) (int, error) {
	cost, err := strconv.ParseFloat(costStr, 64)
	if err != nil {
		return 0, fmt.Errorf("неверный формат стоимости: %w", err)
	}

	maxWeight, err := strconv.ParseFloat(maxWeightStr, 64)
	if err != nil {
		return 0, fmt.Errorf("неверный формат максимального веса: %w", err)
	}

	packagingID, err := packagingService.CreatePackaging(ctx, packagingType, cost, maxWeight)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания упаковки: %w", err)
	}

	return packagingID, nil
}

// GetAllPackaging возвращает все упаковки
func GetAllPackaging(ctx context.Context, packagingService *service.PackagingService) ([]model.PackagingOption, error) {
	packagingOptions, err := packagingService.GetAllPackaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения всех упаковок: %w", err)
	}
	return packagingOptions, nil
}

// DeletePackaging удаляет упаковку по ее ID
func DeletePackaging(ctx context.Context, packagingService *service.PackagingService, packagingIDStr string) error {
	packagingID, err := strconv.Atoi(packagingIDStr)
	if err != nil {
		return fmt.Errorf("неверный формат идентификатора упаковки: %w", err)
	}

	if err := packagingService.DeletePackaging(ctx, packagingID); err != nil {
		return fmt.Errorf("ошибка удаления упаковки с ID %d: %w", packagingID, err)
	}

	return nil
}

// GetPackagingTypeByID возвращает тип упаковки по ее ID
func GetPackagingTypeByID(ctx context.Context, packagingService *service.PackagingService, packagingID int) (string, error) {
	packaging, err := packagingService.GetPackagingByID(ctx, packagingID)
	if err != nil {
		return "", fmt.Errorf("ошибка получения упаковки: %w", err)
	}
	return packaging.Type, nil
}
