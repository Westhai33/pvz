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

// DeleteReturn удаляет возврат, используя контроллер
func (s *APIServiceServer) DeleteReturn(ctx context.Context, req *v1.DeleteReturnRequest) (*emptypb.Empty, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}

	// Удаление возврата через контроллер
	if err := controller.DeleteReturn(ctx, s.returnService, fmt.Sprintf("%d", req.ReturnId)); err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка удаления возврата: %v", err)
	}

	// Возвращаем пустой ответ
	return &emptypb.Empty{}, nil
}
