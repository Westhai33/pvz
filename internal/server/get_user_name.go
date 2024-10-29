package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
)

// GetUserName получает имя пользователя по ID, используя контроллер
func (s *APIServiceServer) GetUserName(ctx context.Context, req *v1.GetUserNameRequest) (*v1.GetUserNameResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код ошибки InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Получение имени пользователя через контроллер
	username, err := controller.GetUserNameByID(ctx, s.userService, int(req.UserId))
	if err != nil {
		// Возвращаем код ошибки NotFound, если пользователь не найден
		return nil, status.Errorf(codes.NotFound, "ошибка получения имени пользователя: %v", err)
	}

	// Возвращаем успешный ответ с именем пользователя
	return &v1.GetUserNameResponse{
		Username: username,
	}, nil
}
