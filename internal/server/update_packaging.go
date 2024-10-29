package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "homework1/internal/api/v1"
	"homework1/internal/model"
)

// UpdatePackaging обновляет информацию об упаковке
func (s *APIServiceServer) UpdatePackaging(ctx context.Context, req *v1.UpdatePackagingRequest) (*v1.UpdatePackagingResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код ошибки InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Создание объекта PackagingOption для обновления данных
	packaging := model.PackagingOption{
		PackagingID: int(req.PackagingId),
		Type:        req.PackagingType,
		Cost:        req.Cost,
		MaxWeight:   req.MaxWeight,
	}

	// Обновление упаковки через сервис
	s.packagingService.UpdatePackaging(ctx, packaging)

	// Возвращаем успешный ответ
	return &v1.UpdatePackagingResponse{
		Message: "Упаковка успешно обновлена",
	}, nil
}
