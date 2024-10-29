package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework1/internal/api/v1"
	"homework1/internal/controller"
)

// CreateReturn создает новый возврат, используя контроллер
func (s *APIServiceServer) CreateReturn(ctx context.Context, req *v1.CreateReturnRequest) (*v1.CreateReturnResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Используем код ошибки InvalidArgument для некорректных данных запроса
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Создание возврата через контроллер
	err := controller.CreateReturn(ctx, s.returnService, fmt.Sprint(req.OrderId))
	if err != nil {
		// Возвращаем код ошибки Internal для внутренних ошибок
		return nil, status.Errorf(codes.Internal, "ошибка создания возврата: %v", err)
	}

	// Возвращаем успешный ответ
	return &v1.CreateReturnResponse{
		Message: "Возврат успешно создан",
	}, nil
}
