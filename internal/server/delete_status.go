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

// DeleteStatus удаляет статус по ID
func (s *APIServiceServer) DeleteStatus(ctx context.Context, req *v1.DeleteStatusRequest) (*emptypb.Empty, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}

	// Удаление статуса через контроллер
	if err := controller.DeleteStatus(ctx, s.statusService, fmt.Sprintf("%d", req.StatusId)); err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка удаления статуса: %v", err)
	}

	// Возвращаем пустой ответ
	return &emptypb.Empty{}, nil
}
