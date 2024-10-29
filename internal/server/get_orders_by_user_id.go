package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "homework1/internal/api/v1"
	"homework1/internal/controller"
)

// GetOrdersByUserID возвращает заказы по ID пользователя
func (s *APIServiceServer) GetOrdersByUserID(ctx context.Context, req *v1.GetOrdersByUserIDRequest) (*v1.GetOrdersByUserIDResponse, error) {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		// Возвращаем код ошибки InvalidArgument, если данные запроса некорректны
		return nil, status.Errorf(codes.InvalidArgument, "ошибка валидации: %v", err)
	}

	// Получение заказов по ID пользователя через контроллер
	orders, err := controller.GetOrdersByUserID(ctx, s.orderService, int(req.UserId))
	if err != nil {
		// Возвращаем код ошибки Internal, если возникла проблема при получении заказов
		return nil, status.Errorf(codes.Internal, "ошибка получения заказов для пользователя с ID %d: %v", req.UserId, err)
	}

	// Формирование списка заказов для ответа
	var grpcOrders []*v1.GetOrderResponse
	for _, order := range orders {
		grpcOrder := &v1.GetOrderResponse{
			OrderId:        int32(order.OrderID),
			UserId:         int32(order.UserID),
			PackagingId:    int32(order.PackagingID),
			StatusId:       int32(order.StatusID),
			AcceptanceDate: order.AcceptanceDate.Format("2006-01-02"),
			ExpirationDate: order.ExpirationDate.Format("2006-01-02"),
			Weight:         order.Weight,
			BaseCost:       order.BaseCost,
			PackagingCost:  order.PackagingCost,
			TotalCost:      order.TotalCost,
			WithFilm:       order.WithFilm,
			IssueDate:      order.IssueDate.Format("2006-01-02"),
		}
		grpcOrders = append(grpcOrders, grpcOrder)
	}

	// Возвращаем успешный ответ
	return &v1.GetOrdersByUserIDResponse{
		Orders: grpcOrders,
	}, nil
}
