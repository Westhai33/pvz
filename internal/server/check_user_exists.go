package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework1/internal/api/v1"
)

// CheckUserExists проверяет, существует ли пользователь по ID, используя контроллер
func (s *APIServiceServer) CheckUserExists(ctx context.Context, req *v1.CheckUserExistsRequest) (*v1.CheckUserExistsResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код ошибки InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Проверяем существование пользователя через userService
	exists, err := s.userService.CheckUserExists(ctx, int(req.UserId))
	if err != nil {
		// Возвращаем код ошибки Internal, если возникли проблемы при проверке пользователя
		return nil, status.Errorf(codes.Internal, "ошибка проверки существования пользователя: %v", err)
	}

	// Если пользователь не существует, можно вернуть NotFound
	if !exists {
		return nil, status.Errorf(codes.NotFound, "пользователь с ID %d не найден", req.UserId)
	}

	// Возвращаем успешный ответ, если пользователь существует
	return &v1.CheckUserExistsResponse{
		Exists: exists,
	}, nil
}
