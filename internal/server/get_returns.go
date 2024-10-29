package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
)

// GetReturns получает все возвраты, используя контроллер
func (s *APIServiceServer) GetReturns(ctx context.Context, req *emptypb.Empty) (*v1.GetReturnsResponse, error) {
	// Получение всех возвратов через контроллер
	returns, err := controller.GetReturns(ctx, s.returnService)
	if err != nil {
		// Возвращаем код ошибки Internal, если возникла проблема при получении возвратов
		return nil, status.Errorf(codes.Internal, "ошибка получения всех возвратов: %v", err)
	}

	// Формирование списка возвратов для ответа
	var returnResponses []*v1.ReturnResponse
	for _, ret := range returns {
		returnResponses = append(returnResponses, &v1.ReturnResponse{
			ReturnId:      int32(ret.ReturnID),
			OrderId:       int32(ret.OrderID),
			UserId:        int32(ret.UserID),
			ReasonId:      int32(ret.ReasonID),
			BaseCost:      float32(ret.BaseCost),
			PackagingCost: float32(ret.PackagingCost),
			TotalCost:     float32(ret.TotalCost),
			PackagingId:   int32(ret.PackagingID),
			StatusId:      int32(ret.StatusID),
			ReturnDate:    ret.ReturnDate.Format("2006-01-02"),
		})
	}

	// Возвращаем успешный ответ
	return &v1.GetReturnsResponse{
		Returns: returnResponses,
	}, nil
}
