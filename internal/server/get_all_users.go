package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework1/internal/api/v1"
	"homework1/internal/controller"
	"time"
)

// GetAllUsers получает всех пользователей, используя контроллер
func (s *APIServiceServer) GetAllUsers(ctx context.Context, req *emptypb.Empty) (*v1.GetAllUsersResponse, error) {
	// Получение всех пользователей через контроллер
	users, err := controller.GetAllUsers(ctx, s.userService)
	if err != nil {
		// Возвращаем код ошибки Internal, если возникли проблемы при получении пользователей
		return nil, status.Errorf(codes.Internal, "ошибка получения всех пользователей: %v", err)
	}

	// Формирование ответа для gRPC
	var userList []*v1.User
	for _, user := range users {
		userList = append(userList, &v1.User{
			UserId:    int32(user.UserID),
			Username:  user.Username,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		})
	}

	// Возвращаем успешный ответ
	return &v1.GetAllUsersResponse{
		Users: userList,
	}, nil
}
