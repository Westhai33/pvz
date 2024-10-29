package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework1/internal/api/v1"
	"homework1/internal/controller"
)

// CreateUser создает нового пользователя, используя контроллер
func (s *APIServiceServer) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем ошибку с кодом InvalidArgument, если запрос некорректен
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Создание пользователя через контроллер
	userID, err := controller.CreateUser(ctx, s.userService, req.Username)
	if err != nil {
		// Возвращаем ошибку с кодом Internal, если произошла внутренняя ошибка при создании пользователя
		return nil, status.Errorf(codes.Internal, "ошибка создания пользователя: %v", err)
	}

	// Возвращаем успешный ответ
	return &v1.CreateUserResponse{
		UserId:  int32(userID),
		Message: "Пользователь успешно создан",
	}, nil
}
