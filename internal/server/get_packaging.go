package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
)

// GetPackaging возвращает информацию об упаковке по ее ID
func (s *APIServiceServer) GetPackaging(ctx context.Context, req *v1.GetPackagingRequest) (*v1.GetPackagingResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код ошибки InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Получение типа упаковки через контроллер
	packagingType, err := controller.GetPackagingTypeByID(ctx, s.packagingService, int(req.PackagingId))
	if err != nil {
		// Возвращаем код ошибки Internal, если возникла ошибка при получении типа упаковки
		return nil, status.Errorf(codes.Internal, "ошибка получения типа упаковки: %v", err)
	}

	// Получение информации об упаковке через сервис
	packaging, err := s.packagingService.GetPackagingByID(ctx, int(req.PackagingId))
	if err != nil {
		// Возвращаем код ошибки Internal, если возникла ошибка при получении упаковки
		return nil, status.Errorf(codes.Internal, "ошибка получения упаковки: %v", err)
	}

	// Возвращаем успешный ответ
	return &v1.GetPackagingResponse{
		PackagingId:   int32(packaging.PackagingID),
		PackagingType: packagingType,
		Cost:          packaging.Cost,
		MaxWeight:     packaging.MaxWeight,
	}, nil
}
