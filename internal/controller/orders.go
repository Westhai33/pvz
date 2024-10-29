package controller

import (
	"context"
	"fmt"
	"homework1/internal/metrics"
	"homework1/internal/model"
	"homework1/internal/service"
	"strconv"
	"time"
)

// CreateOrder создает новый заказ для пользователя с использованием worker pool
func CreateOrder(ctx context.Context, orderService *service.OrderService, userService *service.UserService, packagingService *service.PackagingService, userIDStr, packagingIDStr, expirationDateStr, weightStr, baseCostStr, withFilmStr string) (int, error) {
	userID, err := parseUserID(ctx, userIDStr, userService)
	if err != nil {
		return 0, err
	}

	packagingID, err := parsePackagingID(ctx, packagingIDStr, packagingService)
	if err != nil {
		return 0, err
	}

	expirationDate, err := parseExpirationDate(expirationDateStr)
	if err != nil {
		return 0, err
	}

	weight, err := parseWeight(weightStr)
	if err != nil {
		return 0, err
	}

	baseCost, err := parseBaseCost(baseCostStr)
	if err != nil {
		return 0, err
	}

	withFilm, err := parseWithFilm(withFilmStr)
	if err != nil {
		return 0, err
	}

	packaging, err := getPackaging(ctx, packagingID, packagingService)
	if err != nil {
		return 0, err
	}

	if err := validateWeight(packaging, weight, packagingService); err != nil {
		return 0, err
	}

	totalCost, packagingCost := calculateTotalCost(baseCost, packaging.Cost, withFilm)

	orderID, err := orderService.CreateOrder(ctx, userID, packagingID, 1, expirationDate, weight, baseCost, packagingCost, totalCost, withFilm)
	if err != nil {
		return 0, fmt.Errorf("ошибка при создании заказа: %v", err)
	}

	return orderID, nil
}

// GetOrders возвращает все заказы через сервис с использованием worker pool
func GetOrders(ctx context.Context, orderService *service.OrderService) ([]model.Order, error) {
	orders, err := orderService.GetAllOrders(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения всех заказов: %w", err)
	}
	return orders, nil
}

// GetOrderByID возвращает заказ по его идентификатору через сервис с использованием worker pool
func GetOrderByID(ctx context.Context, orderService *service.OrderService, orderIDStr string) (*model.Order, error) {
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат идентификатора заказа")
	}

	order, err := orderService.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения заказа с ID %d: %w", orderID, err)
	}

	return order, nil
}

// DeleteOrder удаляет заказ через сервис с использованием worker pool
func DeleteOrder(ctx context.Context, orderService *service.OrderService, orderIDStr string) error {
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		return fmt.Errorf("неверный формат идентификатора заказа")
	}

	orderService.DeleteOrder(ctx, orderID)
	return nil
}

// IssueOrder помечает заказ как выданный через сервис с использованием worker pool
func IssueOrder(ctx context.Context, orderService *service.OrderService, statusService *service.StatusService, orderIDStr string) error {
	orderID, err := parseOrderID(orderIDStr)
	if err != nil {
		return err
	}

	order, err := getOrder(ctx, orderID, orderService)
	if err != nil {
		return err
	}

	if err := checkIfIssued(order); err != nil {
		return err
	}

	if err := checkExpiration(order); err != nil {
		return err
	}

	order.IssueDate = time.Now()
	if err := updateOrderStatus(ctx, order, "Выдан", statusService); err != nil {
		return err
	}

	orderService.UpdateOrder(ctx, *order)

	metrics.IncrementIssuedOrders("issued")

	return nil
}

// parseOrderID парсит строку в идентификатор заказа
func parseOrderID(orderIDStr string) (int, error) {
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		return 0, fmt.Errorf("неверный формат идентификатора заказа")
	}
	return orderID, nil
}

// getOrder получает заказ по идентификатору
func getOrder(ctx context.Context, orderID int, orderService *service.OrderService) (*model.Order, error) {
	order, err := orderService.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения заказа с ID %d: %w", orderID, err)
	}
	return order, nil
}

// checkIfIssued проверяет, был ли заказ уже выдан
func checkIfIssued(order *model.Order) error {
	if order.IssueDate != (time.Time{}) {
		return fmt.Errorf("заказ с ID %d уже был выдан", order.OrderID)
	}
	return nil
}

// checkExpiration проверяет, не истек ли срок хранения заказа
func checkExpiration(order *model.Order) error {
	today := time.Now().Truncate(24 * time.Hour)
	expirationDate := order.ExpirationDate.Truncate(24 * time.Hour)
	if expirationDate.Before(today) {
		return fmt.Errorf("срок хранения заказа с ID %d истёк", order.OrderID)
	}
	return nil
}

// updateOrderStatus обновляет статус заказа
func updateOrderStatus(ctx context.Context, order *model.Order, statusName string, statusService *service.StatusService) error {
	issuedStatusID, err := statusService.GetStatusByName(ctx, statusName)
	if err != nil {
		return fmt.Errorf("ошибка получения ID статуса '%s': %w", statusName, err)
	}
	order.StatusID = issuedStatusID
	return nil
}

// GetOrdersByUserID возвращает все заказы по UserID
func GetOrdersByUserID(ctx context.Context, orderService *service.OrderService, userID int) ([]model.Order, error) {
	return orderService.GetOrdersByUserID(ctx, userID)
}

// parseUserID парсит строку в идентификатор пользователя
func parseUserID(ctx context.Context, userIDStr string, userService *service.UserService) (int, error) {
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, fmt.Errorf("неверный формат идентификатора пользователя")
	}

	userExists, err := userService.CheckUserExists(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("ошибка проверки существования пользователя: %w", err)
	}
	if !userExists {
		return 0, fmt.Errorf("пользователь с ID %d не найден", userID)
	}
	return userID, nil
}

// parsePackagingID парсит строку в идентификатор упаковки
func parsePackagingID(ctx context.Context, packagingIDStr string, packagingService *service.PackagingService) (int, error) {
	packagingID, err := strconv.Atoi(packagingIDStr)
	if err != nil {
		return 0, fmt.Errorf("неверный формат идентификатора упаковки")
	}

	packagingExists, err := packagingService.CheckPackagingExists(ctx, packagingID)
	if err != nil {
		return 0, fmt.Errorf("ошибка проверки существования упаковки: %w", err)
	}
	if !packagingExists {
		return 0, fmt.Errorf("упаковка с ID %d не найдена", packagingID)
	}
	return packagingID, nil
}

// parseExpirationDate парсит строку в дату окончания срока хранения
func parseExpirationDate(expirationDateStr string) (time.Time, error) {
	expirationDate, err := time.Parse("2006-01-02", expirationDateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("неверный формат даты. Используйте YYYY-MM-DD")
	}

	today := time.Now().Truncate(24 * time.Hour)
	if expirationDate.Before(today) {
		return time.Time{}, fmt.Errorf("ошибка: дата окончания срока хранения не может быть просроченной")
	}

	return expirationDate, nil
}

// parseWeight парсит строку в вес заказа
func parseWeight(weightStr string) (float64, error) {
	weight, err := strconv.ParseFloat(weightStr, 64)
	if err != nil {
		return 0, fmt.Errorf("неверный формат веса")
	}
	return weight, nil
}

// parseBaseCost парсит строку в базовую стоимость
func parseBaseCost(baseCostStr string) (float64, error) {
	baseCost, err := strconv.ParseFloat(baseCostStr, 64)
	if err != nil {
		return 0, fmt.Errorf("неверный формат базовой стоимости")
	}
	return baseCost, nil
}

// parseWithFilm парсит строку в булевое значение для наличия плёнки
func parseWithFilm(withFilmStr string) (bool, error) {
	if withFilmStr == "y" {
		return true, nil
	} else if withFilmStr == "n" {
		return false, nil
	}
	return false, fmt.Errorf("неверный формат для withFilm, используйте 'y' или 'n'")
}

// getPackaging получает информацию об упаковке по идентификатору
func getPackaging(ctx context.Context, packagingID int, packagingService *service.PackagingService) (*model.PackagingOption, error) {
	packaging, err := packagingService.GetPackagingByID(ctx, packagingID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения информации об упаковке: %w", err)
	}
	return packaging, nil
}

// validateWeight проверяет вес заказа в зависимости от типа упаковки
func validateWeight(packaging *model.PackagingOption, weight float64, packagingService *service.PackagingService) error {
	filmPackagingID, err := packagingService.GetPackagingIDByName(context.Background(), "Пленка")
	if err != nil {
		return fmt.Errorf("ошибка при получении ID упаковки 'Плёнка': %w", err)
	}
	if packaging.PackagingID != filmPackagingID && weight > packaging.MaxWeight {
		return fmt.Errorf("вес заказа превышает максимальный допустимый вес для данной упаковки: %.2f кг", packaging.MaxWeight)
	}
	return nil
}

// calculateTotalCost вычисляет общую стоимость заказа
func calculateTotalCost(baseCost, packagingCost float64, withFilm bool) (float64, float64) {
	if withFilm {
		packagingCost += 1 // добавляем дополнительную стоимость за плёнку
	}
	totalCost := baseCost + packagingCost
	return totalCost, packagingCost
}

// UpdateOrder обновляет заказ через сервис
func UpdateOrder(ctx context.Context, orderService *service.OrderService, order model.Order) error {
	orderService.UpdateOrder(ctx, order)
	return nil
}
