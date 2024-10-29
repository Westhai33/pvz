package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework1/internal/api/v1"
	"homework1/internal/controller"
)

// CreatePackaging создает новую упаковку
func (s *APIServiceServer) CreatePackaging(ctx context.Context, req *v1.CreatePackagingRequest) (*v1.CreatePackagingResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем ошибку InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Создание упаковки с помощью контроллера
	packagingID, err := controller.CreatePackaging(ctx, s.packagingService, req.PackagingType,
		fmt.Sprintf("%f", req.Cost), fmt.Sprintf("%f", req.MaxWeight))
	if err != nil {
		// Возвращаем ошибку Internal в случае ошибки во время создания упаковки
		return nil, status.Errorf(codes.Internal, "ошибка создания упаковки: %v", err)
	}

	// Возвращаем успешный ответ с ID созданной упаковки
	return &v1.CreatePackagingResponse{
		PackagingId: int32(packagingID),
	}, nil
}
