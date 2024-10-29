package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework1/internal/api/v1"
	"homework1/internal/controller"
)

// GetAllPackaging возвращает все упаковки
func (s *APIServiceServer) GetAllPackaging(ctx context.Context, req *emptypb.Empty) (*v1.GetAllPackagingResponse, error) {
	// Получение всех упаковок через контроллер
	packagingOptions, err := controller.GetAllPackaging(ctx, s.packagingService)
	if err != nil {
		// Возвращаем код ошибки Internal, если возникли проблемы при получении данных
		return nil, status.Errorf(codes.Internal, "ошибка получения всех упаковок: %v", err)
	}

	// Формирование списка упаковок для ответа
	var options []*v1.Packaging
	for _, p := range packagingOptions {
		options = append(options, &v1.Packaging{
			PackagingId:   int32(p.PackagingID),
			PackagingType: p.Type,
			Cost:          p.Cost,
			MaxWeight:     p.MaxWeight,
		})
	}

	// Возвращаем успешный ответ
	return &v1.GetAllPackagingResponse{
		PackagingOptions: options,
	}, nil
}
