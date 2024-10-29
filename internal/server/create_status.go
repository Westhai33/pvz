package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework1/internal/api/v1"
	"homework1/internal/controller"
)

// CreateStatus создает новый статус через gRPC, используя контроллер
func (s *APIServiceServer) CreateStatus(ctx context.Context, req *v1.CreateStatusRequest) (*v1.CreateStatusResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код ошибки InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Создание статуса через контроллер
	statusID, err := controller.CreateStatus(ctx, s.statusService, req.StatusName)
	if err != nil {
		// Возвращаем код ошибки Internal, если возникли внутренние проблемы при создании статуса
		return nil, status.Errorf(codes.Internal, "ошибка создания статуса: %v", err)
	}

	// Возвращаем успешный ответ
	return &v1.CreateStatusResponse{
		StatusId: int32(statusID),
		Message:  "Статус успешно создан",
	}, nil
}
