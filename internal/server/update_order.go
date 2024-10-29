package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
	"homework1/internal/model"
	"time"
)

// UpdateOrder обновляет информацию о заказе
func (s *APIServiceServer) UpdateOrder(ctx context.Context, req *v1.UpdateOrderRequest) (*v1.UpdateOrderResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код ошибки InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}

	// Парсинг даты истечения
	expirationDate, err := time.Parse("2006-01-02", req.ExpirationDate)
	if err != nil {
		// Возвращаем код ошибки InvalidArgument, если формат даты истечения некорректен
		return nil, status.Errorf(codes.InvalidArgument, "invalid expiration date format: %v", err)
	}

	// Парсинг даты выдачи
	issueDate, err := time.Parse("2006-01-02", req.IssueDate)
	if err != nil {
		// Возвращаем код ошибки InvalidArgument, если формат даты выдачи некорректен
		return nil, status.Errorf(codes.InvalidArgument, "invalid issue date format: %v", err)
	}

	// Обновление данных заказа
	order := model.Order{
		OrderID:        int(req.OrderId),
		UserID:         int(req.UserId),
		PackagingID:    int(req.PackagingId),
		StatusID:       int(req.StatusId),
		AcceptanceDate: time.Now(),
		ExpirationDate: expirationDate,
		Weight:         req.Weight,
		BaseCost:       req.BaseCost,
		PackagingCost:  req.PackagingCost,
		TotalCost:      req.TotalCost,
		WithFilm:       req.WithFilm,
		IssueDate:      issueDate,
	}

	// Обновление заказа через контроллер
	if err := controller.UpdateOrder(ctx, s.orderService, order); err != nil {
		// Возвращаем код ошибки Internal, если возникла проблема при обновлении заказа
		return nil, status.Errorf(codes.Internal, "ошибка обновления заказа: %v", err)
	}

	// Возвращаем успешный ответ
	return &v1.UpdateOrderResponse{
		Message: "Order updated successfully",
	}, nil
}
