package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
)

// GetReturnByOrderID получает возврат по ID заказа, используя контроллер
func (s *APIServiceServer) GetReturnByOrderID(ctx context.Context, req *v1.GetReturnByOrderIDRequest) (*v1.GetReturnByOrderIDResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код ошибки InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Получение возврата по ID заказа через контроллер
	ret, err := controller.GetReturnByOrderID(ctx, s.returnService, fmt.Sprint(req.OrderId))
	if err != nil {
		// Возвращаем код ошибки Internal, если возникла проблема при получении возврата
		return nil, status.Errorf(codes.Internal, "ошибка получения возврата: %v", err)
	}

	// Возвращаем успешный ответ с данными возврата
	return &v1.GetReturnByOrderIDResponse{
		ReturnInfo: &v1.ReturnResponse{
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
		},
	}, nil
}
