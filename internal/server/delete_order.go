package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework1/internal/api/v1"
	"homework1/internal/controller"
)

// DeleteOrder удаляет заказ
func (s *APIServiceServer) DeleteOrder(ctx context.Context, req *v1.DeleteOrderRequest) (*emptypb.Empty, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код ошибки InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}

	// Удаление заказа через контроллер
	if err := controller.DeleteOrder(ctx, s.orderService, fmt.Sprintf("%d", req.OrderId)); err != nil {
		// Возвращаем код ошибки Internal, если возникли проблемы при удалении заказа
		return nil, status.Errorf(codes.Internal, "ошибка удаления заказа: %v", err)
	}

	// Возвращаем успешный ответ
	return &emptypb.Empty{}, nil
}
