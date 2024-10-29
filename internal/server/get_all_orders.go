package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework1/internal/api/v1"
	"homework1/internal/controller"
)

// GetAllOrders возвращает все заказы
func (s *APIServiceServer) GetAllOrders(ctx context.Context, req *emptypb.Empty) (*v1.GetAllOrdersResponse, error) {
	// Получение всех заказов через контроллер
	orders, err := controller.GetOrders(ctx, s.orderService)
	if err != nil {
		// Возвращаем код ошибки Internal, если возникли проблемы при получении заказов
		return nil, status.Errorf(codes.Internal, "ошибка получения всех заказов: %v", err)
	}

	// Формирование ответа для gRPC
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
	return &v1.GetAllOrdersResponse{
		Orders: grpcOrders,
	}, nil
}
