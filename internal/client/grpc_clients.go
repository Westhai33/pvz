package client

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "homework1/internal/api/v1"
)

// APIServiceClientWrapper обертка для APIServiceClient
type APIServiceClientWrapper struct {
	client v1.APIServiceClient
	conn   *grpc.ClientConn
}

// NewAPIServiceClientWrapper создает новый экземпляр обертки APIServiceClientWrapper
func NewAPIServiceClientWrapper(conn *grpc.ClientConn) (*APIServiceClientWrapper, error) {
	if conn == nil {
		log.Println("gRPC подключение не может быть nil")
		return nil, status.Error(codes.InvalidArgument, "gRPC подключение отсутствует")
	}

	client := v1.NewAPIServiceClient(conn)
	return &APIServiceClientWrapper{client: client, conn: conn}, nil
}

// SetWorkerCount проксирует запрос к SetWorkerCount gRPC методу
func (w *APIServiceClientWrapper) SetWorkerCount(ctx context.Context, req *v1.SetWorkerCountRequest) error {
	_, err := w.client.SetWorkerCount(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова SetWorkerCount: %v", err)
		return err
	}
	return nil
}

// SeedOrders проксирует запрос к SeedOrders gRPC методу
func (w *APIServiceClientWrapper) SeedOrders(ctx context.Context, req *v1.SeedOrdersRequest) error {
	_, err := w.client.SeedOrders(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова SeedOrders: %v", err)
		return err
	}
	return nil
}

// CreateOrder проксирует запрос к CreateOrder gRPC методу
func (w *APIServiceClientWrapper) CreateOrder(ctx context.Context, req *v1.CreateOrderRequest) (*v1.CreateOrderResponse, error) {
	resp, err := w.client.CreateOrder(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова CreateOrder: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetAllOrders проксирует запрос к GetAllOrders gRPC методу
func (w *APIServiceClientWrapper) GetAllOrders(ctx context.Context) (*v1.GetAllOrdersResponse, error) {
	req := &emptypb.Empty{} // Пустой запрос вместо GetAllOrdersRequest
	resp, err := w.client.GetAllOrders(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова GetAllOrders: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetOrdersByUserID проксирует запрос к GetOrdersByUserID gRPC методу
func (w *APIServiceClientWrapper) GetOrdersByUserID(ctx context.Context, req *v1.GetOrdersByUserIDRequest) (*v1.GetOrdersByUserIDResponse, error) {
	resp, err := w.client.GetOrdersByUserID(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова GetOrdersByUserID: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetOrder проксирует запрос к GetOrder gRPC методу
func (w *APIServiceClientWrapper) GetOrder(ctx context.Context, req *v1.GetOrderRequest) (*v1.GetOrderResponse, error) {
	resp, err := w.client.GetOrder(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова GetOrder: %v", err)
		return nil, err
	}
	return resp, nil
}

// IssueOrder проксирует запрос к IssueOrder gRPC методу
func (w *APIServiceClientWrapper) IssueOrder(ctx context.Context, req *v1.IssueOrderRequest) error {
	_, err := w.client.IssueOrder(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова IssueOrder: %v", err)
		return err
	}
	return nil
}

// DeleteOrder проксирует запрос к DeleteOrder gRPC методу
func (w *APIServiceClientWrapper) DeleteOrder(ctx context.Context, req *v1.DeleteOrderRequest) error {
	_, err := w.client.DeleteOrder(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова DeleteOrder: %v", err)
		return err
	}
	return nil
}

// Close закрывает соединение gRPC
func (w *APIServiceClientWrapper) Close() error {
	if w.conn != nil {
		if err := w.conn.Close(); err != nil {
			log.Printf("Ошибка при закрытии gRPC соединения: %v", err)
			return err
		}
		log.Println("gRPC соединение закрыто успешно")
	} else {
		log.Println("gRPC соединение уже закрыто или не было установлено")
	}
	return nil
}

// CreatePackaging проксирует запрос к CreatePackaging gRPC методу
func (w *APIServiceClientWrapper) CreatePackaging(ctx context.Context, req *v1.CreatePackagingRequest) (*v1.CreatePackagingResponse, error) {
	resp, err := w.client.CreatePackaging(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова CreatePackaging: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetPackaging проксирует запрос к GetPackaging gRPC методу
func (w *APIServiceClientWrapper) GetPackaging(ctx context.Context, req *v1.GetPackagingRequest) (*v1.GetPackagingResponse, error) {
	resp, err := w.client.GetPackaging(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова GetPackaging: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetAllPackaging проксирует запрос к GetAllPackaging gRPC методу
func (w *APIServiceClientWrapper) GetAllPackaging(ctx context.Context) (*v1.GetAllPackagingResponse, error) {
	req := &emptypb.Empty{} // Пустой запрос для получения всех упаковок
	resp, err := w.client.GetAllPackaging(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова GetAllPackaging: %v", err)
		return nil, err
	}
	return resp, nil
}

// DeletePackaging проксирует запрос к DeletePackaging gRPC методу
func (w *APIServiceClientWrapper) DeletePackaging(ctx context.Context, req *v1.DeletePackagingRequest) error {
	_, err := w.client.DeletePackaging(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова DeletePackaging: %v", err)
		return err
	}
	return nil
}

// CreateReturnReason проксирует запрос к CreateReturnReason gRPC методу
func (w *APIServiceClientWrapper) CreateReturnReason(ctx context.Context, req *v1.CreateReturnReasonRequest) (*v1.CreateReturnReasonResponse, error) {
	resp, err := w.client.CreateReturnReason(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова CreateReturnReason: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetReturnReason проксирует запрос к GetReturnReason gRPC методу
func (w *APIServiceClientWrapper) GetReturnReason(ctx context.Context, req *v1.GetReturnReasonRequest) (*v1.GetReturnReasonResponse, error) {
	resp, err := w.client.GetReturnReason(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова GetReturnReason: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetAllReturnReasons проксирует запрос к GetAllReturnReasons gRPC методу
func (w *APIServiceClientWrapper) GetAllReturnReasons(ctx context.Context) (*v1.GetAllReturnReasonsResponse, error) {
	req := &emptypb.Empty{} // Пустой запрос для получения всех причин возврата
	resp, err := w.client.GetAllReturnReasons(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова GetAllReturnReasons: %v", err)
		return nil, err
	}
	return resp, nil
}

// UpdateReturnReason проксирует запрос к UpdateReturnReason gRPC методу
func (w *APIServiceClientWrapper) UpdateReturnReason(ctx context.Context, req *v1.UpdateReturnReasonRequest) error {
	_, err := w.client.UpdateReturnReason(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова UpdateReturnReason: %v", err)
		return err
	}
	return nil
}

// DeleteReturnReason проксирует запрос к DeleteReturnReason gRPC методу
func (w *APIServiceClientWrapper) DeleteReturnReason(ctx context.Context, req *v1.DeleteReturnReasonRequest) error {
	_, err := w.client.DeleteReturnReason(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова DeleteReturnReason: %v", err)
		return err
	}
	return nil
}

// CreateReturn проксирует запрос к CreateReturn gRPC методу
func (w *APIServiceClientWrapper) CreateReturn(ctx context.Context, req *v1.CreateReturnRequest) (*v1.CreateReturnResponse, error) {
	resp, err := w.client.CreateReturn(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова CreateReturn: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetReturns проксирует запрос к GetReturns gRPC методу
func (w *APIServiceClientWrapper) GetReturns(ctx context.Context, req *emptypb.Empty) (*v1.GetReturnsResponse, error) {
	resp, err := w.client.GetReturns(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова GetReturns: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetReturnByOrderID проксирует запрос к GetReturnByOrderID gRPC методу
func (w *APIServiceClientWrapper) GetReturnByOrderID(ctx context.Context, req *v1.GetReturnByOrderIDRequest) (*v1.GetReturnByOrderIDResponse, error) {
	resp, err := w.client.GetReturnByOrderID(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова GetReturnByOrderID: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetReturnsByUserID проксирует запрос к GetReturnsByUserID gRPC методу
func (w *APIServiceClientWrapper) GetReturnsByUserID(ctx context.Context, req *v1.GetReturnsByUserIDRequest) (*v1.GetReturnsByUserIDResponse, error) {
	resp, err := w.client.GetReturnsByUserID(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова GetReturnsByUserID: %v", err)
		return nil, err
	}
	return resp, nil
}

// DeleteReturn проксирует запрос к DeleteReturn gRPC методу
func (w *APIServiceClientWrapper) DeleteReturn(ctx context.Context, req *v1.DeleteReturnRequest) error {
	_, err := w.client.DeleteReturn(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова DeleteReturn: %v", err)
		return err
	}
	return nil
}

// ProcessReturn проксирует запрос к ProcessReturn gRPC методу
func (w *APIServiceClientWrapper) ProcessReturn(ctx context.Context, req *v1.ProcessReturnRequest) error {
	_, err := w.client.ProcessReturn(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова ProcessReturn: %v", err)
		return err
	}
	return nil
}

// CreateStatus проксирует запрос к CreateStatus gRPC методу
func (w *APIServiceClientWrapper) CreateStatus(ctx context.Context, req *v1.CreateStatusRequest) (*v1.CreateStatusResponse, error) {
	resp, err := w.client.CreateStatus(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова CreateStatus: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetStatusByID проксирует запрос к GetStatusByID gRPC методу
func (w *APIServiceClientWrapper) GetStatusByID(ctx context.Context, req *v1.GetStatusByIDRequest) (*v1.GetStatusByIDResponse, error) {
	resp, err := w.client.GetStatusByID(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова GetStatusByID: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetAllStatuses проксирует запрос к GetAllStatuses gRPC методу
func (w *APIServiceClientWrapper) GetAllStatuses(ctx context.Context) (*v1.GetAllStatusesResponse, error) {
	req := &emptypb.Empty{} // Пустой запрос для получения всех статусов
	resp, err := w.client.GetAllStatuses(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова GetAllStatuses: %v", err)
		return nil, err
	}
	return resp, nil
}

// UpdateStatus проксирует запрос к UpdateStatus gRPC методу
func (w *APIServiceClientWrapper) UpdateStatus(ctx context.Context, req *v1.UpdateStatusRequest) error {
	_, err := w.client.UpdateStatus(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова UpdateStatus: %v", err)
		return err
	}
	return nil
}

// DeleteStatus проксирует запрос к DeleteStatus gRPC методу
func (w *APIServiceClientWrapper) DeleteStatus(ctx context.Context, req *v1.DeleteStatusRequest) error {
	_, err := w.client.DeleteStatus(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова DeleteStatus: %v", err)
		return err
	}
	return nil
}

// CreateUser проксирует запрос к CreateUser gRPC методу
func (w *APIServiceClientWrapper) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	resp, err := w.client.CreateUser(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова CreateUser: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetUser проксирует запрос к GetUser gRPC методу
func (w *APIServiceClientWrapper) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	resp, err := w.client.GetUser(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова GetUser: %v", err)
		return nil, err
	}
	return resp, nil
}

// GetAllUsers проксирует запрос к GetAllUsers gRPC методу
func (w *APIServiceClientWrapper) GetAllUsers(ctx context.Context) (*v1.GetAllUsersResponse, error) {
	req := &emptypb.Empty{} // Пустой запрос для получения всех пользователей
	resp, err := w.client.GetAllUsers(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова GetAllUsers: %v", err)
		return nil, err
	}
	return resp, nil
}

// UpdateUser проксирует запрос к UpdateUser gRPC методу
func (w *APIServiceClientWrapper) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) error {
	_, err := w.client.UpdateUser(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова UpdateUser: %v", err)
		return err
	}
	return nil
}

// DeleteUser проксирует запрос к DeleteUser gRPC методу
func (w *APIServiceClientWrapper) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) error {
	_, err := w.client.DeleteUser(ctx, req)
	if err != nil {
		log.Printf("Ошибка вызова DeleteUser: %v", err)
		return err
	}
	return nil
}
