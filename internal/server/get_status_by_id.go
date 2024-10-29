package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
)

// GetStatusByID получает статус по его ID через gRPC, используя контроллер
func (s *APIServiceServer) GetStatusByID(ctx context.Context, req *v1.GetStatusByIDRequest) (*v1.GetStatusByIDResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Использование status.Errorf из правильного пакета для формирования gRPC ошибки
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Изменяем переменную с именем status, чтобы избежать конфликта с пакетом status
	stat, err := controller.GetStatusByID(ctx, s.statusService, int(req.StatusId))
	if err != nil {
		// Использование status.Errorf из правильного пакета для формирования gRPC ошибки
		return nil, status.Errorf(codes.Internal, "ошибка получения статуса: %v", err)
	}

	// Возвращаем успешный ответ с данными статуса
	return &v1.GetStatusByIDResponse{
		StatusId:   int32(stat.StatusID),
		StatusName: stat.StatusName,
	}, nil
}
