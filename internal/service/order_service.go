package service

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"homework1/internal/tracing"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.opentelemetry.io/otel/trace"
	"homework1/internal/cache"
	"homework1/internal/dao"
	"homework1/internal/kafka"
	"homework1/internal/model"
	"homework1/internal/pool"
)

// OrderService представляет сервис для работы с заказами
type OrderService struct {
	pool     *pgxpool.Pool
	wp       *pool.WorkerPool
	producer *kafka.Producer
	cache    *cache.RedisCache[string, model.Order]
	tracer   trace.Tracer
}

// NewOrderService создает новый сервис для работы с заказами
func NewOrderService(dbPool *pgxpool.Pool, workerPool *pool.WorkerPool, producer *kafka.Producer, cache *cache.RedisCache[string, model.Order]) *OrderService {
	return &OrderService{
		pool:     dbPool,
		wp:       workerPool,
		producer: producer,
		cache:    cache,
		tracer:   tracing.GetTracer(),
	}
}

// CreateOrder создает новый заказ и отправляет событие в Kafka
func (s *OrderService) CreateOrder(ctx context.Context, userID, packagingID, statusID int, expirationDate time.Time, weight, baseCost, packagingCost, totalCost float64, withFilm bool) (int, error) {
	ctx, span := s.tracer.Start(ctx, "CreateOrder")
	defer span.End()

	resultChan := make(chan int, 1)
	errChan := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(1)
	s.wp.SubmitTask(func() {
		defer wg.Done()

		newOrder := s.buildOrder(userID, packagingID, statusID, expirationDate, weight, baseCost, packagingCost, totalCost, withFilm)

		orderID, err := s.saveOrder(ctx, newOrder)
		if err != nil {
			errChan <- err
			return
		}

		if err := s.cache.Delete(ctx, "all_orders"); err != nil {
			log.Printf("Ошибка инвалидации кэша всех заказов: %v", err)
		}

		orderCacheKey := fmt.Sprintf("order_%d", orderID)
		if err := s.cache.Delete(ctx, orderCacheKey); err != nil {
			log.Printf("Ошибка инвалидации кэша для заказа %d: %v", orderID, err)
		}

		userOrdersCacheKey := fmt.Sprintf("user_orders_%d", userID)
		if err := s.cache.Delete(ctx, userOrdersCacheKey); err != nil {
			log.Printf("Ошибка инвалидации кэша для заказов пользователя %d: %v", userID, err)
		}

		if err := s.notifyOrderCreation(orderID); err != nil {
			errChan <- fmt.Errorf("Ошибка отправки сообщения в Kafka: %v", err)
		} else {
			resultChan <- orderID
		}
	})

	wg.Wait()
	return s.getResultOrError(resultChan, errChan)
}

// buildOrder создает новый объект заказа
func (s *OrderService) buildOrder(userID, packagingID, statusID int, expirationDate time.Time, weight, baseCost, packagingCost, totalCost float64, withFilm bool) model.Order {
	return model.Order{
		UserID:         userID,
		PackagingID:    packagingID,
		StatusID:       statusID,
		AcceptanceDate: time.Now(),
		ExpirationDate: expirationDate,
		Weight:         weight,
		BaseCost:       baseCost,
		PackagingCost:  packagingCost,
		TotalCost:      totalCost,
		WithFilm:       withFilm,
	}
}

// saveOrder сохраняет заказ в базе данных
func (s *OrderService) saveOrder(ctx context.Context, order model.Order) (int, error) {
	orderID, err := dao.CreateOrder(ctx, order, s.pool)
	if err != nil {
		log.Printf("Ошибка создания заказа: %v", err)
		return 0, err
	}
	log.Printf("Заказ создан с ID %d", orderID)
	return orderID, nil
}

// notifyOrderCreation отправляет уведомление в Kafka о создании заказа
func (s *OrderService) notifyOrderCreation(orderID int) error {
	return s.sendKafkaMessage("create", orderID, fmt.Sprintf("Order %d created", orderID))
}

// getResultOrError возвращает результат или ошибку
func (s *OrderService) getResultOrError(resultChan chan int, errChan chan error) (int, error) {
	select {
	case orderID := <-resultChan:
		return orderID, nil
	case err := <-errChan:
		return 0, err
	}
}

// updateOrderFields обновляет поля существующего заказа на основе нового
func updateOrderFields(existingOrder, updatedOrder *model.Order) {
	existingOrder.UserID = updatedOrder.UserID
	existingOrder.AcceptanceDate = updatedOrder.AcceptanceDate
	existingOrder.ExpirationDate = updatedOrder.ExpirationDate
	existingOrder.Weight = updatedOrder.Weight
	existingOrder.BaseCost = updatedOrder.BaseCost
	existingOrder.PackagingCost = updatedOrder.PackagingCost
	existingOrder.TotalCost = updatedOrder.TotalCost
	existingOrder.PackagingID = updatedOrder.PackagingID
	existingOrder.StatusID = updatedOrder.StatusID
	existingOrder.IssueDate = updatedOrder.IssueDate
	existingOrder.WithFilm = updatedOrder.WithFilm
}

// sendKafkaMessage отправляет сообщение в Kafka
func (s *OrderService) sendKafkaMessage(operation string, orderID int, description string) error {
	orderMessage := kafka.OrderMessage{
		TimeStamp:   time.Now(),
		Operation:   operation,
		OrderID:     orderID,
		Description: description,
	}
	return s.producer.SendOrderMessage(orderMessage)
}

// GetOrderByID возвращает заказ по ID
func (s *OrderService) GetOrderByID(ctx context.Context, orderID int) (*model.Order, error) {
	ctx, span := s.tracer.Start(ctx, "GetOrderByID")
	defer span.End()

	// Проверка кэша
	cacheKey := fmt.Sprintf("order_%d", orderID)
	var cachedOrder model.Order
	if err := s.cache.Get(ctx, cacheKey, &cachedOrder); err == nil {
		return &cachedOrder, nil
	}

	order, err := dao.GetOrderByID(ctx, orderID, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения заказа с ID %d: %w", orderID, err)
	}

	if err := s.cache.Set(ctx, cacheKey, *order); err != nil {
		log.Printf("Ошибка сохранения заказа в кэше: %v", err)
	}

	return order, nil
}

// UpdateOrder обновляет заказ и отправляет событие в Kafka
func (s *OrderService) UpdateOrder(ctx context.Context, order model.Order) {
	ctx, span := s.tracer.Start(ctx, "UpdateOrder")
	defer span.End()

	var wg sync.WaitGroup
	wg.Add(1)
	s.wp.SubmitTask(func() {
		defer wg.Done()

		existingOrder, err := dao.GetOrderByID(ctx, order.OrderID, s.pool)
		if err != nil {
			log.Printf("Ошибка поиска заказа с ID %d: %v", order.OrderID, err)
			return
		}

		updateOrderFields(existingOrder, &order)

		if err := dao.UpdateOrder(ctx, *existingOrder, s.pool); err != nil {
			log.Printf("Ошибка обновления заказа с ID %d: %v", order.OrderID, err)
			return
		}

		orderCacheKey := fmt.Sprintf("order_%d", order.OrderID)
		if err := s.cache.Delete(ctx, orderCacheKey); err != nil {
			log.Printf("Ошибка инвалидации кэша для заказа %d: %v", order.OrderID, err)
		}

		if err := s.cache.Delete(ctx, "all_orders"); err != nil {
			log.Printf("Ошибка инвалидации кэша всех заказов: %v", err)
		}

		userOrdersCacheKey := fmt.Sprintf("user_orders_%d", order.UserID)
		if err := s.cache.Delete(ctx, userOrdersCacheKey); err != nil {
			log.Printf("Ошибка инвалидации кэша для заказов пользователя %d: %v", order.UserID, err)
		}

		if err := s.sendKafkaMessage("update", order.OrderID, fmt.Sprintf("Order %d updated", order.OrderID)); err != nil {
			log.Printf("Ошибка отправки сообщения в Kafka: %v", err)
		}

		log.Printf("Заказ с ID %d успешно обновлен", order.OrderID)
	})

	wg.Wait()
}

// DeleteOrder удаляет заказ и отправляет событие в Kafka
func (s *OrderService) DeleteOrder(ctx context.Context, orderID int) {
	ctx, span := s.tracer.Start(ctx, "DeleteOrder")
	defer span.End()

	var wg sync.WaitGroup
	wg.Add(1)
	s.wp.SubmitTask(func() {
		defer wg.Done()

		order, err := dao.GetOrderByID(ctx, orderID, s.pool)
		if err != nil {
			log.Printf("Заказ с ID %d не найден: %v", orderID, err)
			return
		}

		if err := dao.DeleteOrder(ctx, orderID, s.pool); err != nil {
			log.Printf("Ошибка удаления заказа с ID %d: %v", orderID, err)
			return
		}

		orderCacheKey := fmt.Sprintf("order_%d", orderID)
		if err := s.cache.Delete(ctx, orderCacheKey); err != nil {
			log.Printf("Ошибка инвалидации кэша для заказа %d: %v", orderID, err)
		}

		if err := s.cache.Delete(ctx, "all_orders"); err != nil {
			log.Printf("Ошибка инвалидации кэша всех заказов: %v", err)
		}

		userOrdersCacheKey := fmt.Sprintf("user_orders_%d", order.UserID)
		if err := s.cache.Delete(ctx, userOrdersCacheKey); err != nil {
			log.Printf("Ошибка инвалидации кэша для заказов пользователя %d: %v", order.UserID, err)
		}

		orderMessage := kafka.OrderMessage{
			TimeStamp:   time.Now(),
			Operation:   "delete",
			OrderID:     orderID,
			Description: fmt.Sprintf("Order %d deleted", orderID),
		}

		if err := s.producer.SendOrderMessage(orderMessage); err != nil {
			log.Printf("Ошибка отправки сообщения в Kafka: %v", err)
		}

		log.Printf("Заказ с ID %d успешно удален", orderID)
	})

	wg.Wait()
}

// CheckExpiredOrders проверяет и обрабатывает просроченные заказы через worker pool
func (s *OrderService) CheckExpiredOrders(ctx context.Context) error {
	ctx, span := s.tracer.Start(ctx, "CheckExpiredOrders")
	defer span.End()

	expiredOrders, err := dao.GetExpiredOrders(ctx, s.pool)
	if err != nil {
		errMsg := fmt.Sprintf("ошибка при получении просроченных заказов: %v", err)
		if kafkaErr := s.producer.SendKafkaErrorMessage("check_expired_orders", 0, errMsg); kafkaErr != nil {
			log.Printf("Ошибка при отправке сообщения в Kafka: %v", kafkaErr)
		}
		return fmt.Errorf("%s", errMsg)
	}

	reasonID, err := s.getReturnReasonID(ctx, "Истек срок хранения")
	if err != nil {
		if kafkaErr := s.producer.SendKafkaErrorMessage("check_expired_orders", 0, fmt.Sprintf("ошибка получения причины возврата: %v", err)); kafkaErr != nil {
			log.Printf("Ошибка при отправке сообщения в Kafka: %v", kafkaErr)
		}
		return err
	}

	returnStatusID, err := s.getStatusID(ctx, "Возврат")
	if err != nil {
		if kafkaErr := s.producer.SendKafkaErrorMessage("check_expired_orders", 0, fmt.Sprintf("ошибка получения статуса 'Возврат': %v", err)); kafkaErr != nil {
			log.Printf("Ошибка при отправке сообщения в Kafka: %v", kafkaErr)
		}
		return err
	}

	createdStatusID, err := s.getStatusID(ctx, "Создан")
	if err != nil {
		if kafkaErr := s.producer.SendKafkaErrorMessage("check_expired_orders", 0, fmt.Sprintf("ошибка получения статуса 'Создан': %v", err)); kafkaErr != nil {
			log.Printf("Ошибка при отправке сообщения в Kafka: %v", kafkaErr)
		}
		return err
	}

	return s.processExpiredOrders(ctx, expiredOrders, createdStatusID, returnStatusID, reasonID)
}

// processExpiredOrders обрабатывает просроченные заказы
func (s *OrderService) processExpiredOrders(ctx context.Context, expiredOrders []model.Order, createdStatusID, returnStatusID, reasonID int) error {
	for _, order := range expiredOrders {
		if order.StatusID == createdStatusID {
			if err := s.handleExpiredOrder(ctx, order, reasonID, returnStatusID); err != nil {
				if kafkaErr := s.producer.SendKafkaErrorMessage("process_expired_orders", order.OrderID, fmt.Sprintf("ошибка при обработке заказа с ID %d: %v", order.OrderID, err)); kafkaErr != nil {
					log.Printf("Ошибка при отправке сообщения в Kafka: %v", kafkaErr)
				}
				return fmt.Errorf("ошибка при обработке заказа с ID %d: %w", order.OrderID, err)
			}
		} else if order.StatusID == returnStatusID {
			continue
		}
	}

	return nil
}

func (s *OrderService) handleExpiredOrder(ctx context.Context, order model.Order, reasonID, returnStatusID int) error {
	var wg sync.WaitGroup

	wg.Add(1)
	s.wp.SubmitTask(func() {
		defer wg.Done()

		if err := s.createReturn(ctx, order, reasonID, returnStatusID); err != nil {
			log.Printf("Ошибка создания возврата: %v", err)
			return
		}

		if err := s.updateOrderStatus(ctx, order, returnStatusID); err != nil {
			log.Printf("Ошибка обновления статуса заказа: %v", err)
			return
		}

		if err := s.notifyKafkaReturn(order.OrderID); err != nil {
			log.Printf("Ошибка отправки сообщения в Kafka: %v", err)
		}
	})

	wg.Wait()
	return nil
}

// createReturn создает новый возврат для заказа
func (s *OrderService) createReturn(ctx context.Context, order model.Order, reasonID, returnStatusID int) error {
	newReturn := model.Return{
		OrderID:       order.OrderID,
		UserID:        order.UserID,
		ReturnDate:    time.Now().UTC(),
		ReasonID:      reasonID,
		BaseCost:      order.BaseCost,
		PackagingCost: order.PackagingCost,
		TotalCost:     order.TotalCost,
		PackagingID:   order.PackagingID,
		StatusID:      returnStatusID,
	}
	err := dao.CreateReturn(ctx, newReturn, s.pool)
	if err != nil {
		log.Printf("Ошибка создания возврата: %v", err)
		return err
	}

	// Инвалидация кэша для возвратов
	if err := s.cache.Delete(ctx, fmt.Sprintf("return_%d", newReturn.OrderID)); err != nil {
		log.Printf("Ошибка инвалидации кэша для возврата заказа %d: %v", newReturn.OrderID, err)
	}

	return nil
}

// updateOrderStatus обновляет статус заказа
func (s *OrderService) updateOrderStatus(ctx context.Context, order model.Order, returnStatusID int) error {
	order.StatusID = returnStatusID
	return dao.UpdateOrder(ctx, order, s.pool)
}

// notifyKafkaReturn отправляет сообщение в Kafka о возврате
func (s *OrderService) notifyKafkaReturn(orderID int) error {
	orderMessage := kafka.OrderMessage{
		TimeStamp:   time.Now(),
		Operation:   "return",
		OrderID:     orderID,
		Description: fmt.Sprintf("Возврат создан для заказа с ID %d", orderID),
	}
	return s.producer.SendOrderMessage(orderMessage)
}

// getReturnReasonID получает ID причины возврата по её имени
func (s *OrderService) getReturnReasonID(ctx context.Context, reasonName string) (int, error) {
	reason, err := dao.GetReturnReasonByName(ctx, reasonName, s.pool)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения причины возврата: %w", err)
	}
	return reason.ReasonID, nil
}

// getStatusID получает ID статуса по его имени
func (s *OrderService) getStatusID(ctx context.Context, statusName string) (int, error) {
	status, err := dao.GetStatusByName(ctx, statusName, s.pool)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения статуса: %w", err)
	}
	return status.StatusID, nil
}

// GetAllOrders получает все заказы из базы данных.
func (s *OrderService) GetAllOrders(ctx context.Context) ([]model.Order, error) {
	// Проверка кэша
	var cachedOrders []model.Order
	if err := s.cache.GetSlice(ctx, "all_orders", &cachedOrders); err == nil {
		return cachedOrders, nil
	}

	orders, err := dao.GetAllOrders(ctx, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения всех заказов: %w", err)
	}

	if err := s.cache.SetSlice(ctx, "all_orders", orders); err != nil {
		log.Printf("Ошибка сохранения всех заказов в кэше: %v", err)
	}

	return orders, nil
}

// SeedOrders создает фейковые заказы
func (s *OrderService) SeedOrders(ctx context.Context, num int) error {
	ctx, span := s.tracer.Start(ctx, "SeedOrders")
	defer span.End()

	for i := 0; i < num; i++ {
		order := model.Order{
			UserID:         gofakeit.Number(1, 3),
			AcceptanceDate: gofakeit.Date(),
			ExpirationDate: gofakeit.Date(),
			Weight:         gofakeit.Float64Range(0, 100.),
			BaseCost:       gofakeit.Float64Range(10, 500),
			PackagingCost:  gofakeit.Float64Range(1, 20),
			TotalCost:      gofakeit.Float64Range(15, 600),
			PackagingID:    gofakeit.Number(1, 3),
			StatusID:       gofakeit.Number(1, 4),
			IssueDate:      gofakeit.Date(),
			WithFilm:       gofakeit.Bool(),
		}

		orderID, err := dao.CreateOrder(ctx, order, s.pool)
		if err != nil {
			return fmt.Errorf("ошибка при создании заказа: %w", err)
		}

		cacheKey := fmt.Sprintf("order_%d", orderID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка инвалидации кэша для заказа %d: %v", orderID, err)
		}
	}

	if err := s.cache.Delete(ctx, "all_orders"); err != nil {
		log.Printf("Ошибка инвалидации кэша всех заказов после массовой генерации: %v", err)
	}

	userKeys, err := s.cache.Keys(ctx, "user_orders_*")
	if err != nil {
		log.Printf("Ошибка получения ключей кэша для заказов пользователей: %v", err)
	} else {
		for _, key := range userKeys {
			if err := s.cache.Delete(ctx, key); err != nil {
				log.Printf("Ошибка инвалидации кэша для ключа %s: %v", key, err)
			}
		}
	}

	return nil
}

// GetOrdersByUserID получает заказы по идентификатору пользователя.
func (s *OrderService) GetOrdersByUserID(ctx context.Context, userID int) ([]model.Order, error) {
	cacheKey := fmt.Sprintf("user_orders_%d", userID)
	var cachedOrders []model.Order
	if err := s.cache.GetSlice(ctx, cacheKey, &cachedOrders); err == nil {
		return cachedOrders, nil
	}

	orders, err := dao.GetOrdersByUserID(ctx, userID, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения заказов для пользователя с ID %d: %w", userID, err)
	}

	if err := s.cache.SetSlice(ctx, cacheKey, orders); err != nil {
		log.Printf("Ошибка сохранения заказов пользователя в кэше: %v", err)
	}

	return orders, nil
}
