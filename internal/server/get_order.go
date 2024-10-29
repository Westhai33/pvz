package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
)

// GetOrder возвращает информацию о заказе по его ID
func (s *APIServiceServer) GetOrder(ctx context.Context, req *v1.GetOrderRequest) (*v1.GetOrderResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код ошибки InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}

	// Получение заказа через контроллер
	order, err := controller.GetOrderByID(ctx, s.orderService, fmt.Sprintf("%d", req.OrderId))
	if err != nil {
		// Возвращаем код ошибки Internal, если возникла внутренняя ошибка при получении заказа
		return nil, status.Errorf(codes.Internal, "ошибка при получении заказа: %v", err)
	}

	// Возвращаем успешный ответ с данными заказа
	return &v1.GetOrderResponse{
		OrderId:        int32(order.OrderID),
		UserId:         int32(order.UserID),
		PackagingId:    int32(order.PackagingID),
		StatusId:       int32(order.StatusID),
		AcceptanceDate: order.AcceptanceDate.Format("2006-01-02"),
		ExpirationDate: order.ExpirationDate.Format("2006-01-02"),
		Weight:         order.Weight,
		BaseCost:       order.BaseCost,
		PackagingCost:  order.PackagingCost,
		TotalCost:      order.TotalCost,
		WithFilm:       order.WithFilm,
		IssueDate:      order.IssueDate.Format("2006-01-02"),
	}, nil
}
