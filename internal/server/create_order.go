package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
)

func (s *APIServiceServer) CreateOrder(ctx context.Context, req *v1.CreateOrderRequest) (*v1.CreateOrderResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Используем код ошибки InvalidArgument для ошибок валидации
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}

	// Преобразование поля "WithFilm" в строку для обработки
	withFilmStr := "n"
	if req.WithFilm {
		withFilmStr = "y"
	}

	// Вызов контроллера для создания заказа
	orderID, err := controller.CreateOrder(ctx, s.orderService, s.userService, s.packagingService,
		fmt.Sprintf("%d", req.UserId),
		fmt.Sprintf("%d", req.PackagingId),
		req.ExpirationDate,
		fmt.Sprintf("%f", req.Weight),
		fmt.Sprintf("%f", req.BaseCost),
		withFilmStr)
	if err != nil {
		// Используем код ошибки Internal для ошибок во время создания заказа
		return nil, status.Errorf(codes.Internal, "ошибка создания заказа: %v", err)
	}

	// Возвращаем успешный ответ
	return &v1.CreateOrderResponse{
		OrderId: int32(orderID),
		Message: "Заказ успешно создан",
	}, nil
}
