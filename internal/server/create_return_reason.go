package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework1/internal/api/v1"
	"homework1/internal/controller"
)

// CreateReturnReason создает новую причину возврата через gRPC
func (s *APIServiceServer) CreateReturnReason(ctx context.Context, req *v1.CreateReturnReasonRequest) (*v1.CreateReturnReasonResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код ошибки InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Создание причины возврата через контроллер
	reasonID, err := controller.CreateReturnReason(ctx, s.returnReasonService, req.Reason)
	if err != nil {
		// Возвращаем код ошибки Internal, если возникли внутренние проблемы при создании причины
		return nil, status.Errorf(codes.Internal, "ошибка создания причины возврата: %v", err)
	}

	// Возвращаем успешный ответ
	return &v1.CreateReturnReasonResponse{
		ReasonId: int32(reasonID),
		Message:  "Причина возврата успешно создана",
	}, nil
}
