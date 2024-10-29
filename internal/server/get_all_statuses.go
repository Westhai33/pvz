package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework1/internal/api/v1"
	"homework1/internal/controller"
)

// GetAllStatuses получает все статусы через gRPC, используя контроллер
func (s *APIServiceServer) GetAllStatuses(ctx context.Context, req *emptypb.Empty) (*v1.GetAllStatusesResponse, error) {
	// Получение всех статусов через контроллер
	statuses, err := controller.GetAllStatuses(ctx, s.statusService)
	if err != nil {
		// Возвращаем код ошибки Internal, если возникла внутренняя ошибка при получении данных
		return nil, status.Errorf(codes.Internal, "ошибка получения всех статусов: %v", err)
	}

	// Формирование ответа для gRPC
	var statusResponses []*v1.GetStatusByIDResponse
	for _, status := range statuses {
		statusResponses = append(statusResponses, &v1.GetStatusByIDResponse{
			StatusId:   int32(status.StatusID),
			StatusName: status.StatusName,
		})
	}

	// Возвращаем успешный ответ
	return &v1.GetAllStatusesResponse{
		Statuses: statusResponses,
	}, nil
}
