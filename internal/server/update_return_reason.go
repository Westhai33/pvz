package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
)

// UpdateReturnReason обновляет причину возврата через gRPC
func (s *APIServiceServer) UpdateReturnReason(ctx context.Context, req *v1.UpdateReturnReasonRequest) (*v1.UpdateReturnReasonResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код ошибки InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Обновление причины возврата через контроллер
	err := controller.UpdateReturnReason(ctx, s.returnReasonService, fmt.Sprintf("%d", req.ReasonId), req.Reason)
	if err != nil {
		// Возвращаем код ошибки Internal, если возникла проблема при обновлении причины возврата
		return nil, status.Errorf(codes.Internal, "ошибка обновления причины возврата: %v", err)
	}

	// Возвращаем успешный ответ
	return &v1.UpdateReturnReasonResponse{
		Message: "Причина возврата успешно обновлена",
	}, nil
}
