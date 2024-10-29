package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "homework1/internal/api/v1"
)

// SeedOrders создает фейковые заказы
func (s *APIServiceServer) SeedOrders(ctx context.Context, req *v1.SeedOrdersRequest) (*emptypb.Empty, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}

	// Логика создания фейковых заказов
	if err := s.orderService.SeedOrders(ctx, int(req.Count)); err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка создания фейковых заказов: %v", err)
	}

	// Возвращаем пустой ответ
	return &emptypb.Empty{}, nil
}
