package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "homework1/internal/api/v1"
)

// SetWorkerCount изменяет количество воркеров
func (s *APIServiceServer) SetWorkerCount(ctx context.Context, req *v1.SetWorkerCountRequest) (*emptypb.Empty, error) {
	// Валидация входных данных
	if req.Count < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Количество воркеров не может быть отрицательным")
	}

	// Устанавливаем новое количество воркеров в пуле
	s.workerPool.SetWorkerCount(int(req.Count))

	// Возвращаем пустой ответ
	return &emptypb.Empty{}, nil
}
