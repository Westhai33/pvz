package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb" // Ensure this import is included
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
)

// DeletePackaging удаляет упаковку по ID
func (s *APIServiceServer) DeletePackaging(ctx context.Context, req *v1.DeletePackagingRequest) (*emptypb.Empty, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}

	// Удаление упаковки через контроллер
	if err := controller.DeletePackaging(ctx, s.packagingService, fmt.Sprintf("%d", req.PackagingId)); err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка удаления упаковки: %v", err)
	}

	// Возвращаем пустой ответ
	return &emptypb.Empty{}, nil
}
