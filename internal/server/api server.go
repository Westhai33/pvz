package server

import (
	"homework1/internal/api/v1"
	"homework1/internal/pool"
	"homework1/internal/service"
)

// APIServiceServer представляет реализацию интерфейса v1.APIServiceServer
type APIServiceServer struct {
	v1.UnimplementedAPIServiceServer // Встраиваем не реализованный сервер
	userService                      *service.UserService
	orderService                     *service.OrderService
	packagingService                 *service.PackagingService
	returnService                    *service.ReturnService
	returnReasonService              *service.ReturnReasonService
	statusService                    *service.StatusService
	workerPool                       *pool.WorkerPool
}

// NewAPIServiceServer создает новый APIServiceServer
func NewAPIServiceServer(
	userService *service.UserService,
	orderService *service.OrderService,
	packagingService *service.PackagingService,
	returnService *service.ReturnService,
	returnReasonService *service.ReturnReasonService,
	statusService *service.StatusService,
	workerPool *pool.WorkerPool,
) *APIServiceServer {
	return &APIServiceServer{
		userService:         userService,
		orderService:        orderService,
		packagingService:    packagingService,
		returnService:       returnService,
		returnReasonService: returnReasonService,
		statusService:       statusService,
		workerPool:          workerPool,
	}
}
