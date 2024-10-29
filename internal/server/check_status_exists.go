package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework1/internal/api/v1"
	"homework1/internal/controller"
)

// CheckStatusExists проверяет существование статуса через gRPC, используя контроллер
func (s *APIServiceServer) CheckStatusExists(ctx context.Context, req *v1.CheckStatusExistsRequest) (*v1.CheckStatusExistsResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код InvalidArgument для ошибок валидации
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Получаем имя статуса по ID через statusService
	statusName, err := controller.GetStatusNameByID(ctx, s.statusService, int(req.StatusId))
	if err != nil {
		// Возвращаем код Internal для ошибок при проверке статуса
		return nil, status.Errorf(codes.Internal, "ошибка при проверке существования статуса: %v", err)
	}

	// Проверяем, существует ли статус
	exists := statusName != ""
	return &v1.CheckStatusExistsResponse{
		Exists: exists,
	}, nil
}
