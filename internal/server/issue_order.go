package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
)

// IssueOrder обрабатывает выдачу заказа
func (s *APIServiceServer) IssueOrder(ctx context.Context, req *v1.IssueOrderRequest) (*emptypb.Empty, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}

	// Выдача заказа через контроллер, теперь с передачей statusService
	if err := controller.IssueOrder(ctx, s.orderService, s.statusService, fmt.Sprintf("%d", req.OrderId)); err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка выдачи заказа: %v", err)
	}

	// Возвращаем пустой ответ
	return &emptypb.Empty{}, nil
}
