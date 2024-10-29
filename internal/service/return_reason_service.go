package service

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"homework1/internal/cache"
	"homework1/internal/dao"
	"homework1/internal/model"
	"homework1/internal/pool"
	"log"
	"time"

	"go.opentelemetry.io/otel/trace"
	"homework1/internal/tracing"
)

// ReturnReasonService представляет сервис для работы с причинами возврата
type ReturnReasonService struct {
	pool   *pgxpool.Pool
	wp     *pool.WorkerPool
	cache  *cache.RedisCache[string, model.ReturnReason]
	tracer trace.Tracer
}

// NewReturnReasonService создает новый сервис для работы с причинами возврата
func NewReturnReasonService(dbPool *pgxpool.Pool, workerPool *pool.WorkerPool, cache *cache.RedisCache[string, model.ReturnReason]) *ReturnReasonService {
	return &ReturnReasonService{
		pool:   dbPool,
		wp:     workerPool,
		cache:  cache,
		tracer: tracing.GetTracer(),
	}
}

// CreateReturnReason создает новую причину возврата через общий worker pool
func (s *ReturnReasonService) CreateReturnReason(ctx context.Context, reason string) (int, error) {
	ctx, span := s.tracer.Start(ctx, "CreateReturnReason")
	defer span.End()

	var reasonID int
	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		newReason := model.ReturnReason{
			Reason: reason,
		}

		var err error
		reasonID, err = dao.CreateReturnReason(ctx, newReason, s.pool)
		if err != nil {
			errCh <- fmt.Errorf("ошибка создания причины возврата: %w", err)
			return
		}

		cacheKey := fmt.Sprintf("return_reason_%d", reasonID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка удаления кэша причины возврата с ключом %s: %v", cacheKey, err)
		}

		if err := s.cache.Delete(ctx, "all_return_reasons"); err != nil {
			log.Printf("Ошибка удаления кэша всех причин возврата: %v", err)
		}

		errCh <- nil
	})

	err := <-errCh
	return reasonID, err
}

// GetReturnReasonByID возвращает причину возврата по её ID с использованием кэша
func (s *ReturnReasonService) GetReturnReasonByID(ctx context.Context, reasonID int) (*model.ReturnReason, error) {
	ctx, span := s.tracer.Start(ctx, "GetReturnReasonByID")
	defer span.End()

	cacheKey := fmt.Sprintf("return_reason_%d", reasonID)

	var cachedReason model.ReturnReason
	err := s.cache.Get(ctx, cacheKey, &cachedReason)
	if err == nil {
		log.Printf("Причина возврата с ID %d получена из кэша", reasonID)
		return &cachedReason, nil
	} else {
		log.Printf("Причина возврата с ID %d не найдена в кэше, получаем из базы данных", reasonID)
	}

	reason, err := dao.GetReturnReasonByID(ctx, reasonID, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения причины возврата с ID %d: %w", reasonID, err)
	}

	err = s.cache.Set(ctx, cacheKey, *reason, 10*time.Minute)
	if err != nil {
		log.Printf("Ошибка сохранения причины возврата в кэш: %v", err)
	}

	return reason, nil
}

// GetAllReturnReasons возвращает все причины возвратов с использованием кэша
func (s *ReturnReasonService) GetAllReturnReasons(ctx context.Context) ([]model.ReturnReason, error) {
	ctx, span := s.tracer.Start(ctx, "GetAllReturnReasons")
	defer span.End()

	cacheKey := "all_return_reasons"

	var cachedReasons []model.ReturnReason
	err := s.cache.GetSlice(ctx, cacheKey, &cachedReasons)
	if err == nil {
		log.Println("Все причины возврата получены из кэша")
		return cachedReasons, nil
	} else {
		log.Println("Причины возврата не найдены в кэше, получаем из базы данных")
	}

	reasons, err := dao.GetAllReturnReasons(ctx, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения всех причин возвратов: %w", err)
	}

	err = s.cache.SetSlice(ctx, cacheKey, reasons, 10*time.Minute)
	if err != nil {
		log.Printf("Ошибка сохранения всех причин возврата в кэш: %v", err)
	}

	log.Printf("Получены все причины возвратов: %+v", reasons)
	return reasons, nil
}

// UpdateReturnReason обновляет причину возврата и сбрасывает кэш
func (s *ReturnReasonService) UpdateReturnReason(ctx context.Context, reasonID int, reason string) error {
	ctx, span := s.tracer.Start(ctx, "UpdateReturnReason")
	defer span.End()

	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		exists, err := dao.CheckReturnReasonExists(ctx, reasonID, s.pool)
		if err != nil {
			errCh <- fmt.Errorf("ошибка проверки существования причины возврата с ID %d: %w", reasonID, err)
			return
		}
		if !exists {
			errCh <- fmt.Errorf("причина возврата с ID %d не найдена", reasonID)
			return
		}

		updatedReason := model.ReturnReason{
			ReasonID: reasonID,
			Reason:   reason,
		}

		if err := dao.UpdateReturnReason(ctx, updatedReason, s.pool); err != nil {
			errCh <- fmt.Errorf("ошибка обновления причины возврата с ID %d: %w", reasonID, err)
			return
		}

		cacheKey := fmt.Sprintf("return_reason_%d", reasonID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка удаления кэша причины возврата с ключом %s: %v", cacheKey, err)
		}

		if err := s.cache.Delete(ctx, "all_return_reasons"); err != nil {
			log.Printf("Ошибка удаления кэша всех причин возврата: %v", err)
		}

		errCh <- nil
	})

	return <-errCh
}

// DeleteReturnReason удаляет причину возврата и сбрасывает кэш
func (s *ReturnReasonService) DeleteReturnReason(ctx context.Context, reasonID int) error {
	ctx, span := s.tracer.Start(ctx, "DeleteReturnReason")
	defer span.End()

	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		exists, err := dao.CheckReturnReasonExists(ctx, reasonID, s.pool)
		if err != nil {
			errCh <- fmt.Errorf("ошибка проверки существования причины возврата с ID %d: %w", reasonID, err)
			return
		}
		if !exists {
			errCh <- fmt.Errorf("причина возврата с ID %d не найдена", reasonID)
			return
		}

		if err := dao.DeleteReturnReason(ctx, reasonID, s.pool); err != nil {
			errCh <- fmt.Errorf("ошибка удаления причины возврата с ID %d: %w", reasonID, err)
			return
		}

		cacheKey := fmt.Sprintf("return_reason_%d", reasonID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка удаления кэша причины возврата с ключом %s: %v", cacheKey, err)
		}

		if err := s.cache.Delete(ctx, "all_return_reasons"); err != nil {
			log.Printf("Ошибка удаления кэша всех причин возврата: %v", err)
		}

		errCh <- nil
	})

	return <-errCh
}

// CheckReturnReasonExists проверяет существование причины возврата через worker pool
func (s *ReturnReasonService) CheckReturnReasonExists(ctx context.Context, reasonID int) (bool, error) {
	ctx, span := s.tracer.Start(ctx, "CheckReturnReasonExists")
	defer span.End()

	exists, err := dao.CheckReturnReasonExists(ctx, reasonID, s.pool)
	if err != nil {
		return false, fmt.Errorf("ошибка проверки существования причины возврата с ID %d: %w", reasonID, err)
	}

	log.Printf("Причина возврата с ID %d существует: %v", reasonID, exists)
	return exists, nil
}
