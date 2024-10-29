package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
)

// DeleteReturnReason удаляет причину возврата по ID
func (s *APIServiceServer) DeleteReturnReason(ctx context.Context, req *v1.DeleteReturnReasonRequest) (*emptypb.Empty, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}

	// Удаление причины возврата через контроллер
	if err := controller.DeleteReturnReason(ctx, s.returnReasonService, fmt.Sprintf("%d", req.ReasonId)); err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка удаления причины возврата: %v", err)
	}

	// Возвращаем пустой ответ
	return &emptypb.Empty{}, nil
}
