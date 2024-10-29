package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
)

// GetReturnReason получает причину возврата по её ID через gRPC
func (s *APIServiceServer) GetReturnReason(ctx context.Context, req *v1.GetReturnReasonRequest) (*v1.GetReturnReasonResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код ошибки InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Получение причины возврата по ID через контроллер
	reason, err := controller.GetReturnReasonByID(ctx, s.returnReasonService, int(req.ReasonId))
	if err != nil {
		// Возвращаем код ошибки Internal, если возникла проблема при получении причины возврата
		return nil, status.Errorf(codes.Internal, "ошибка получения причины возврата: %v", err)
	}

	// Возвращаем успешный ответ с данными причины возврата
	return &v1.GetReturnReasonResponse{
		ReasonId: int32(reason.ReasonID),
		Reason:   reason.Reason,
	}, nil
}
