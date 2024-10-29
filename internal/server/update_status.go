package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
)

// UpdateStatus обновляет статус через gRPC, используя контроллер
func (s *APIServiceServer) UpdateStatus(ctx context.Context, req *v1.UpdateStatusRequest) (*v1.UpdateStatusResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код ошибки InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Обновление статуса через контроллер
	err := controller.UpdateStatus(ctx, s.statusService, fmt.Sprintf("%d", req.StatusId), req.StatusName)
	if err != nil {
		// Возвращаем код ошибки Internal, если возникла проблема при обновлении статуса
		return nil, status.Errorf(codes.Internal, "ошибка обновления статуса: %v", err)
	}

	// Возвращаем успешный ответ
	return &v1.UpdateStatusResponse{
		Message: "Статус успешно обновлен",
	}, nil
}
