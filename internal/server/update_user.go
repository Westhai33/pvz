package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
)

// UpdateUser обновляет информацию о существующем пользователе, используя контроллер
func (s *APIServiceServer) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код ошибки InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Обновление данных пользователя через контроллер
	err := controller.UpdateUser(ctx, s.userService, int(req.UserId), req.Username)
	if err != nil {
		// Возвращаем код ошибки Internal, если возникла проблема при обновлении данных пользователя
		return nil, status.Errorf(codes.Internal, "ошибка обновления пользователя: %v", err)
	}

	// Возвращаем успешный ответ
	return &v1.UpdateUserResponse{
		Message: "Пользователь успешно обновлен",
	}, nil
}
