package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
)

// GetReturnsByUserID получает возвраты по ID пользователя, используя контроллер
func (s *APIServiceServer) GetReturnsByUserID(ctx context.Context, req *v1.GetReturnsByUserIDRequest) (*v1.GetReturnsByUserIDResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код ошибки InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации запроса: %v", err)
	}

	// Получение возвратов по ID пользователя через контроллер
	returns, err := controller.GetReturnsByUserID(ctx, s.returnService, fmt.Sprint(req.UserId))
	if err != nil {
		// Возвращаем код ошибки Internal, если возникла проблема при получении возвратов
		return nil, status.Errorf(codes.Internal, "ошибка получения возвратов для пользователя: %v", err)
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
	return &v1.GetReturnsByUserIDResponse{
		Returns: returnResponses,
	}, nil
}
