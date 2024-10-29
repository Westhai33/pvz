package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
	"time"
)

// GetUser получает пользователя по ID, используя контроллер
func (s *APIServiceServer) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код ошибки InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Получение пользователя через контроллер
	user, err := controller.GetUserByID(ctx, s.userService, int(req.UserId))
	if err != nil {
		// Возвращаем код ошибки NotFound, если пользователь не найден
		return nil, status.Errorf(codes.NotFound, "ошибка получения пользователя: %v", err)
	}

	// Возвращаем успешный ответ с данными пользователя
	return &v1.GetUserResponse{
		UserId:    int32(user.UserID),
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}
