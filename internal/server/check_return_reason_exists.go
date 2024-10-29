package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework1/internal/api/v1"
)

// CheckReturnReasonExists проверяет существование причины возврата через gRPC
func (s *APIServiceServer) CheckReturnReasonExists(ctx context.Context, req *v1.CheckReturnReasonExistsRequest) (*v1.CheckReturnReasonExistsResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Проверка существования причины возврата через service
	exists, err := s.returnReasonService.CheckReturnReasonExists(ctx, int(req.ReasonId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка проверки существования причины возврата: %v", err)
	}

	// Возвращаем успешный ответ
	return &v1.CheckReturnReasonExistsResponse{
		Exists: exists,
	}, nil
}
