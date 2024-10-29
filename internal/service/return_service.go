package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"homework1/internal/cache"
	"homework1/internal/dao"
	"homework1/internal/kafka"
	"homework1/internal/model"
	"homework1/internal/pool"
	"log"
	"time"

	"go.opentelemetry.io/otel/trace"
	"homework1/internal/tracing"
)

// ReturnService представляет сервис для работы с возвратами
type ReturnService struct {
	pool     *pgxpool.Pool
	wp       *pool.WorkerPool
	producer *kafka.Producer
	cache    *cache.RedisCache[string, model.Return]
	tracer   trace.Tracer
}

// NewReturnService создает новый сервис для работы с возвратами
func NewReturnService(dbPool *pgxpool.Pool, workerPool *pool.WorkerPool, producer *kafka.Producer, cache *cache.RedisCache[string, model.Return]) *ReturnService {
	return &ReturnService{
		pool:     dbPool,
		wp:       workerPool,
		producer: producer,
		cache:    cache,
		tracer:   tracing.GetTracer(),
	}
}

// CreateReturn создает новый возврат через общий worker pool
func (s *ReturnService) CreateReturn(ctx context.Context, orderID int) error {
	ctx, span := s.tracer.Start(ctx, "CreateReturn")
	defer span.End()

	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		order, err := dao.GetOrderByID(ctx, orderID, s.pool)
		if err != nil {
			s.handleKafkaError("create_return", orderID, fmt.Sprintf("ошибка получения заказа с ID %d: %v", orderID, err))
			errCh <- fmt.Errorf("ошибка получения заказа с ID %d: %v", orderID, err)
			return
		}

		if time.Since(order.IssueDate).Hours() > 48 {
			s.handleKafkaError("create_return", orderID, fmt.Sprintf("возврат заказа с ID %d невозможен, так как прошло более двух дней с момента выдачи", orderID))
			errCh <- fmt.Errorf("возврат заказа с ID %d невозможен, так как прошло более двух дней с момента выдачи", orderID)
			return
		}

		reason, err := dao.GetReturnReasonByName(ctx, "Вернул покупатель", s.pool)
		if err != nil {
			s.handleKafkaError("create_return", orderID, fmt.Sprintf("ошибка получения причины возврата 'Вернул покупатель': %v", err))
			errCh <- fmt.Errorf("ошибка получения причины возврата 'Вернул покупатель': %v", err)
			return
		}
		reasonID := reason.ReasonID

		status, err := dao.GetStatusByName(ctx, "Возврат", s.pool)
		if err != nil {
			s.handleKafkaError("create_return", orderID, fmt.Sprintf("ошибка получения статуса 'возврат': %v", err))
			errCh <- fmt.Errorf("ошибка получения статуса 'возврат': %v", err)
			return
		}

		newReturn := model.Return{
			OrderID:       orderID,
			UserID:        order.UserID,
			ReturnDate:    time.Now().UTC(),
			ReasonID:      reasonID,
			BaseCost:      order.BaseCost,
			PackagingCost: order.PackagingCost,
			PackagingID:   order.PackagingID,
			TotalCost:     order.TotalCost,
			StatusID:      status.StatusID,
		}

		if err := dao.WriteReturns(ctx, []model.Return{newReturn}, s.pool); err != nil {
			s.handleKafkaError("create_return", orderID, fmt.Sprintf("ошибка создания возврата: %v", err))
			errCh <- fmt.Errorf("ошибка создания возврата: %v", err)
			return
		}

		// Сбрасываем кэш для данного возврата
		cacheKey := fmt.Sprintf("return_%d", orderID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка удаления кэша возврата с ключом %s: %v", cacheKey, err)
		}

		message := kafka.OrderMessage{
			TimeStamp:   time.Now(),
			Operation:   "Create Return",
			OrderID:     orderID,
			Description: fmt.Sprintf("Возврат создан для заказа %d", orderID),
		}

		if err := s.producer.SendOrderMessage(message); err != nil {
			s.handleKafkaError("create_return", orderID, fmt.Sprintf("ошибка отправки сообщения в Kafka: %v", err))
			errCh <- fmt.Errorf("ошибка отправки сообщения в Kafka: %v", err)
			return
		}

		errCh <- nil
	})

	return <-errCh
}

// handleKafkaError отправляет сообщение об ошибке в Kafka
func (s *ReturnService) handleKafkaError(operation string, orderID int, errMsg string) {
	if kafkaErr := s.producer.SendKafkaErrorMessage(operation, orderID, errMsg); kafkaErr != nil {
		fmt.Printf("Ошибка отправки сообщения в Kafka: %v\n", kafkaErr)
	}
}

// UpdateReturn обновляет возврат и сбрасывает кэш
func (s *ReturnService) UpdateReturn(ctx context.Context, returnID, orderID, userID, reasonID int, baseCost, packagingCost, totalCost float64, packagingID, statusID int) error {
	ctx, span := s.tracer.Start(ctx, "UpdateReturn")
	defer span.End()

	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		ret, err := dao.FindReturnByOrderID(ctx, orderID, s.pool)
		if err != nil {
			s.handleKafkaError("update_return", orderID, fmt.Sprintf("ошибка поиска возврата с ID %d: %v", returnID, err))
			errCh <- fmt.Errorf("ошибка поиска возврата с ID %d: %v", returnID, err)
			return
		}

		ret.UserID = userID
		ret.ReasonID = reasonID
		ret.BaseCost = baseCost
		ret.PackagingCost = packagingCost
		ret.TotalCost = totalCost
		ret.PackagingID = packagingID
		ret.StatusID = statusID
		ret.ReturnDate = time.Now().UTC()

		if err := dao.UpdateReturn(ctx, *ret, s.pool); err != nil {
			s.handleKafkaError("update_return", orderID, fmt.Sprintf("ошибка обновления возврата с ID %d: %v", returnID, err))
			errCh <- fmt.Errorf("ошибка обновления возврата с ID %d: %v", returnID, err)
			return
		}

		cacheKey := fmt.Sprintf("return_%d", orderID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка удаления кэша возврата с ключом %s: %v", cacheKey, err)
		}

		message := kafka.OrderMessage{
			TimeStamp:   time.Now(),
			Operation:   "Update Return",
			OrderID:     orderID,
			Description: fmt.Sprintf("Возврат обновлен для заказа %d", orderID),
		}

		if err := s.producer.SendOrderMessage(message); err != nil {
			s.handleKafkaError("update_return", orderID, fmt.Sprintf("ошибка отправки сообщения в Kafka: %v", err))
			errCh <- fmt.Errorf("ошибка отправки сообщения в Kafka: %v", err)
			return
		}

		errCh <- nil
	})

	return <-errCh
}

// DeleteReturn удаляет возврат через общий worker pool
func (s *ReturnService) DeleteReturn(ctx context.Context, returnID int) error {
	ctx, span := s.tracer.Start(ctx, "DeleteReturn")
	defer span.End()

	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		exists, err := dao.CheckReturnExists(ctx, returnID, s.pool)
		if err != nil {
			s.handleKafkaError("delete_return", returnID, fmt.Sprintf("ошибка проверки существования возврата с ID %d: %v", returnID, err))
			errCh <- fmt.Errorf("ошибка проверки существования возврата с ID %d: %v", returnID, err)
			return
		}
		if !exists {
			s.handleKafkaError("delete_return", returnID, fmt.Sprintf("возврат с ID %d не найден", returnID))
			errCh <- fmt.Errorf("возврат с ID %d не найден", returnID)
			return
		}

		if err := dao.DeleteReturn(ctx, returnID, s.pool); err != nil {
			s.handleKafkaError("delete_return", returnID, fmt.Sprintf("ошибка удаления возврата с ID %d: %v", returnID, err))
			errCh <- fmt.Errorf("ошибка удаления возврата с ID %d: %v", returnID, err)
			return
		}

		cacheKey := fmt.Sprintf("return_%d", returnID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка удаления кэша возврата с ключом %s: %v", cacheKey, err)
		}

		message := kafka.OrderMessage{
			TimeStamp:   time.Now(),
			Operation:   "Delete Return",
			OrderID:     returnID,
			Description: fmt.Sprintf("Возврат с ID %d удален", returnID),
		}

		if err := s.producer.SendOrderMessage(message); err != nil {
			s.handleKafkaError("delete_return", returnID, fmt.Sprintf("ошибка отправки сообщения в Kafka: %v", err))
			errCh <- fmt.Errorf("ошибка отправки сообщения в Kafka: %v", err)
			return
		}

		errCh <- nil
	})

	return <-errCh
}

// GetReturns возвращает все возвраты через worker pool
func (s *ReturnService) GetReturns(ctx context.Context) ([]model.Return, error) {
	ctx, span := s.tracer.Start(ctx, "GetReturns")
	defer span.End()

	returns, err := dao.ReadReturns(ctx, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения возвратов: %w", err)
	}
	return returns, nil
}

// GetReturnByOrderID возвращает возврат по ID заказа через worker pool
func (s *ReturnService) GetReturnByOrderID(ctx context.Context, orderID int) (*model.Return, error) {
	ctx, span := s.tracer.Start(ctx, "GetReturnByOrderID")
	defer span.End()

	cacheKey := fmt.Sprintf("return_%d", orderID)

	var cachedReturn model.Return
	err := s.cache.Get(ctx, cacheKey, &cachedReturn)
	if err == nil {
		log.Printf("Возврат с ID заказа %d получен из кэша", orderID)
		return &cachedReturn, nil
	} else {
		log.Printf("Возврат с ID заказа %d не найден в кэше, получаем из базы данных", orderID)
	}

	ret, err := dao.FindReturnByOrderID(ctx, orderID, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска возврата для заказа с ID %d: %w", orderID, err)
	}

	err = s.cache.Set(ctx, cacheKey, *ret, 10*time.Minute)
	if err != nil {
		log.Printf("Ошибка сохранения возврата в кэш: %v", err)
	}

	return ret, nil
}

// GetReturnsByUserID возвращает возвраты по ID пользователя через worker pool
func (s *ReturnService) GetReturnsByUserID(ctx context.Context, userID int) ([]model.Return, error) {
	ctx, span := s.tracer.Start(ctx, "GetReturnsByUserID")
	defer span.End()

	returns, err := dao.FindReturnsByUserID(ctx, userID, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска возвратов для пользователя с ID %d: %w", userID, err)
	}
	return returns, nil
}
