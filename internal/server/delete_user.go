package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework1/internal/api/v1"
	"homework1/internal/controller"
)

// DeleteUser удаляет пользователя по ID, используя контроллер
func (s *APIServiceServer) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*emptypb.Empty, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код ошибки InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Удаление пользователя через контроллер
	err := controller.DeleteUser(ctx, s.userService, int(req.UserId))
	if err != nil {
		// Возвращаем код ошибки Internal, если возникли проблемы при удалении пользователя
		return nil, status.Errorf(codes.Internal, "ошибка удаления пользователя: %v", err)
	}

	// Возвращаем успешный ответ
	return &emptypb.Empty{}, nil
}
