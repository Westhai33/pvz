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

// StatusService представляет сервис для работы со статусами
type StatusService struct {
	pool   *pgxpool.Pool
	wp     *pool.WorkerPool
	cache  *cache.RedisCache[string, model.Status]
	tracer trace.Tracer
}

// NewStatusService создает новый сервис для работы со статусами
func NewStatusService(dbPool *pgxpool.Pool, workerPool *pool.WorkerPool, cache *cache.RedisCache[string, model.Status]) *StatusService {
	return &StatusService{
		pool:   dbPool,
		wp:     workerPool,
		cache:  cache,
		tracer: tracing.GetTracer(),
	}
}

// CreateStatus создает новый статус через общий worker pool
func (s *StatusService) CreateStatus(ctx context.Context, statusName string) (int, error) {
	ctx, span := s.tracer.Start(ctx, "CreateStatus")
	defer span.End()

	var statusID int
	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		newStatus := model.Status{
			StatusName: statusName,
		}

		var err error
		statusID, err = dao.CreateStatus(ctx, newStatus, s.pool)
		if err != nil {
			errCh <- fmt.Errorf("ошибка создания статуса: %w", err)
			return
		}

		cacheKey := fmt.Sprintf("status_%d", statusID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка удаления кэша статуса с ключом %s: %v", cacheKey, err)
		}

		errCh <- nil
	})

	err := <-errCh
	return statusID, err
}

// GetStatusByID возвращает статус по его ID с использованием кэша
func (s *StatusService) GetStatusByID(ctx context.Context, statusID int) (*model.Status, error) {
	ctx, span := s.tracer.Start(ctx, "GetStatusByID")
	defer span.End()

	cacheKey := fmt.Sprintf("status_%d", statusID)

	var cachedStatus model.Status
	err := s.cache.Get(ctx, cacheKey, &cachedStatus)
	if err == nil {
		log.Printf("Статус с ID %d получен из кэша", statusID)
		return &cachedStatus, nil
	} else {
		log.Printf("Статус с ID %d не найден в кэше, получаем из базы данных", statusID)
	}

	status, err := dao.GetStatusByID(ctx, statusID, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения статуса с ID %d: %w", statusID, err)
	}

	// Сохраняем в кэш на 10 минут
	err = s.cache.Set(ctx, cacheKey, *status, 10*time.Minute)
	if err != nil {
		log.Printf("Ошибка сохранения статуса в кэш: %v", err)
	}

	return status, nil
}

// GetAllStatuses возвращает все статусы через worker pool
func (s *StatusService) GetAllStatuses(ctx context.Context) ([]model.Status, error) {
	ctx, span := s.tracer.Start(ctx, "GetAllStatuses")
	defer span.End()

	statuses, err := dao.GetAllStatuses(ctx, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения всех статусов: %w", err)
	}

	log.Printf("Получены все статусы: %+v", statuses)
	return statuses, nil
}

// UpdateStatus обновляет статус и сбрасывает кэш
func (s *StatusService) UpdateStatus(ctx context.Context, statusID int, statusName string) error {
	ctx, span := s.tracer.Start(ctx, "UpdateStatus")
	defer span.End()

	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		exists, err := dao.CheckStatusExists(ctx, statusID, s.pool)
		if err != nil {
			errCh <- fmt.Errorf("ошибка проверки существования статуса с ID %d: %w", statusID, err)
			return
		}
		if !exists {
			errCh <- fmt.Errorf("статус с ID %d не найден", statusID)
			return
		}

		updatedStatus := model.Status{
			StatusID:   statusID,
			StatusName: statusName,
		}

		if err := dao.UpdateStatus(ctx, updatedStatus, s.pool); err != nil {
			errCh <- fmt.Errorf("ошибка обновления статуса с ID %d: %w", statusID, err)
			return
		}

		cacheKey := fmt.Sprintf("status_%d", statusID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка удаления кэша статуса с ключом %s: %v", cacheKey, err)
		}

		errCh <- nil
	})

	return <-errCh
}

// DeleteStatus удаляет статус и сбрасывает кэш
func (s *StatusService) DeleteStatus(ctx context.Context, statusID int) error {
	ctx, span := s.tracer.Start(ctx, "DeleteStatus")
	defer span.End()

	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		exists, err := dao.CheckStatusExists(ctx, statusID, s.pool)
		if err != nil {
			errCh <- fmt.Errorf("ошибка проверки существования статуса с ID %d: %w", statusID, err)
			return
		}
		if !exists {
			errCh <- fmt.Errorf("статус с ID %d не найден", statusID)
			return
		}

		if err := dao.DeleteStatus(ctx, statusID, s.pool); err != nil {
			errCh <- fmt.Errorf("ошибка удаления статуса с ID %d: %w", statusID, err)
			return
		}

		cacheKey := fmt.Sprintf("status_%d", statusID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка удаления кэша статуса с ключом %s: %v", cacheKey, err)
		}

		errCh <- nil
	})

	return <-errCh
}

// GetStatusByName возвращает статус по его имени через worker pool
func (s *StatusService) GetStatusByName(ctx context.Context, statusName string) (int, error) {
	ctx, span := s.tracer.Start(ctx, "GetStatusByName")
	defer span.End()

	status, err := dao.GetStatusByName(ctx, statusName, s.pool)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения статуса с именем %s: %w", statusName, err)
	}

	log.Printf("Статус с именем %s имеет ID %d", statusName, status.StatusID)
	return status.StatusID, nil
}

// CheckStatusExists проверяет существование статуса через worker pool
func (s *StatusService) CheckStatusExists(ctx context.Context, statusID int) (bool, error) {
	ctx, span := s.tracer.Start(ctx, "CheckStatusExists")
	defer span.End()

	exists, err := dao.CheckStatusExists(ctx, statusID, s.pool)
	if err != nil {
		return false, fmt.Errorf("ошибка проверки существования статуса с ID %d: %w", statusID, err)
	}

	log.Printf("Статус с ID %d существует: %v", statusID, exists)
	return exists, nil
}

// GetStatusNameByID возвращает имя статуса по ID через worker pool
func (s *StatusService) GetStatusNameByID(ctx context.Context, statusID int) (string, error) {
	ctx, span := s.tracer.Start(ctx, "GetStatusNameByID")
	defer span.End()

	cacheKey := fmt.Sprintf("status_name_%d", statusID)

	var cachedStatusName string
	err := s.cache.Get(ctx, cacheKey, &cachedStatusName)
	if err == nil {
		log.Printf("Имя статуса с ID %d получено из кэша", statusID)
		return cachedStatusName, nil
	} else {
		log.Printf("Имя статуса с ID %d не найдено в кэше, получаем из базы данных", statusID)
	}

	statusName, err := dao.GetStatusNameByID(ctx, statusID, s.pool)
	if err != nil {
		return "", fmt.Errorf("ошибка получения имени статуса с ID %d: %w", statusID, err)
	}

	err = s.cache.SetString(ctx, cacheKey, statusName, 10*time.Minute)
	if err != nil {
		log.Printf("Ошибка сохранения имени статуса в кэш: %v", err)
	}

	return statusName, nil
}
