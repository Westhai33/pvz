package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework1/internal/api/v1"
	"homework1/internal/controller"
)

// GetAllReturnReasons возвращает список всех причин возврата через gRPC
func (s *APIServiceServer) GetAllReturnReasons(ctx context.Context, req *emptypb.Empty) (*v1.GetAllReturnReasonsResponse, error) {
	// Получение всех причин возврата через контроллер
	reasons, err := controller.GetAllReturnReasons(ctx, s.returnReasonService)
	if err != nil {
		// Возвращаем код ошибки Internal, если возникли проблемы при получении данных
		return nil, status.Errorf(codes.Internal, "ошибка получения всех причин возврата: %v", err)
	}

	// Формирование ответа для gRPC
	var reasonResponses []*v1.GetReturnReasonResponse
	for _, reason := range reasons {
		reasonResponses = append(reasonResponses, &v1.GetReturnReasonResponse{
			ReasonId: int32(reason.ReasonID),
			Reason:   reason.Reason,
		})
	}

	// Возвращаем успешный ответ
	return &v1.GetAllReturnReasonsResponse{
		Reasons: reasonResponses,
	}, nil
}
