package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
)

// ProcessReturn обрабатывает процесс возврата
func (s *APIServiceServer) ProcessReturn(ctx context.Context, req *v1.ProcessReturnRequest) (*emptypb.Empty, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}

	// Логика обработки возврата
	if err := controller.ProcessReturn(ctx, s.returnService, s.statusService, int(req.OrderId)); err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка обработки возврата: %v", err)
	}

	// Возвращаем пустой ответ
	return &emptypb.Empty{}, nil
}
