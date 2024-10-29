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

// PackagingService представляет сервис для работы с упаковками
type PackagingService struct {
	pool   *pgxpool.Pool
	wp     *pool.WorkerPool
	cache  *cache.RedisCache[string, model.PackagingOption]
	tracer trace.Tracer
}

// NewPackagingService создает новый сервис для работы с упаковками
func NewPackagingService(dbPool *pgxpool.Pool, workerPool *pool.WorkerPool, cache *cache.RedisCache[string, model.PackagingOption]) *PackagingService {
	return &PackagingService{
		pool:   dbPool,
		wp:     workerPool,
		cache:  cache,
		tracer: tracing.GetTracer(),
	}
}

// CreatePackaging создает новую упаковку через общий worker pool
func (s *PackagingService) CreatePackaging(ctx context.Context, packagingType string, cost, maxWeight float64) (int, error) {
	ctx, span := s.tracer.Start(ctx, "CreatePackaging")
	defer span.End()

	var packagingID int
	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		newPackaging := model.PackagingOption{
			Type:      packagingType,
			Cost:      cost,
			MaxWeight: maxWeight,
		}

		var err error
		packagingID, err = dao.CreatePackaging(ctx, newPackaging, s.pool)
		if err != nil {
			errCh <- fmt.Errorf("ошибка создания упаковки: %w", err)
			return
		}

		cacheKey := fmt.Sprintf("packaging_%d", packagingID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка удаления кэша упаковки с ключом %s: %v", cacheKey, err)
		}

		if err := s.cache.Delete(ctx, "all_packaging"); err != nil {
			log.Printf("Ошибка удаления кэша всех упаковок: %v", err)
		}

		errCh <- nil
	})

	err := <-errCh
	return packagingID, err
}

// GetPackagingByID возвращает упаковку по её ID с использованием кэша
func (s *PackagingService) GetPackagingByID(ctx context.Context, packagingID int) (*model.PackagingOption, error) {
	ctx, span := s.tracer.Start(ctx, "GetPackagingByID")
	defer span.End()

	cacheKey := fmt.Sprintf("packaging_%d", packagingID)

	var cachedPackaging model.PackagingOption
	err := s.cache.Get(ctx, cacheKey, &cachedPackaging)
	if err == nil {
		log.Printf("Упаковка с ID %d получена из кэша", packagingID)
		return &cachedPackaging, nil
	} else {
		log.Printf("Упаковка с ID %d не найдена в кэше, получаем из базы данных", packagingID)
	}

	packaging, err := dao.GetPackagingByID(ctx, packagingID, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения упаковки с ID %d: %w", packagingID, err)
	}

	err = s.cache.Set(ctx, cacheKey, *packaging, 10*time.Minute) // Разыменование указателя
	if err != nil {
		log.Printf("Ошибка сохранения упаковки в кэш: %v", err)
	}

	return packaging, nil
}

// GetAllPackaging возвращает все упаковки с использованием кэша
func (s *PackagingService) GetAllPackaging(ctx context.Context) ([]model.PackagingOption, error) {
	ctx, span := s.tracer.Start(ctx, "GetAllPackaging")
	defer span.End()

	cacheKey := "all_packaging"

	var cachedPackagingOptions []model.PackagingOption
	err := s.cache.GetSlice(ctx, cacheKey, &cachedPackagingOptions) // Используем новый метод для получения среза
	if err == nil {
		log.Println("Все упаковки получены из кэша")
		return cachedPackagingOptions, nil
	} else {
		log.Println("Упаковки не найдены в кэше, получаем из базы данных")
	}

	packagingOptions, err := dao.GetAllPackaging(ctx, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения всех упаковок: %w", err)
	}

	err = s.cache.SetSlice(ctx, cacheKey, packagingOptions, 10*time.Minute) // Используем новый метод для сохранения среза
	if err != nil {
		log.Printf("Ошибка сохранения всех упаковок в кэш: %v", err)
	}

	log.Printf("Получены все упаковки: %+v", packagingOptions)
	return packagingOptions, nil
}

// UpdatePackaging обновляет данные упаковки и сбрасывает кэш
func (s *PackagingService) UpdatePackaging(ctx context.Context, packaging model.PackagingOption) {
	ctx, span := s.tracer.Start(ctx, "UpdatePackaging")
	defer span.End()

	s.wp.SubmitTask(func() {
		existingPackaging, err := dao.GetPackagingByID(ctx, packaging.PackagingID, s.pool)
		if err != nil {
			log.Printf("Ошибка поиска упаковки с ID %d: %v", packaging.PackagingID, err)
			return
		}

		existingPackaging.Type = packaging.Type
		existingPackaging.Cost = packaging.Cost
		existingPackaging.MaxWeight = packaging.MaxWeight

		if err := dao.UpdatePackaging(ctx, *existingPackaging, s.pool); err != nil {
			log.Printf("Ошибка обновления упаковки с ID %d: %v", packaging.PackagingID, err)
			return
		}

		cacheKey := fmt.Sprintf("packaging_%d", packaging.PackagingID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка удаления кэша упаковки с ключом %s: %v", cacheKey, err)
		}

		// Очищаем общий кэш всех упаковок
		if err := s.cache.Delete(ctx, "all_packaging"); err != nil {
			log.Printf("Ошибка удаления кэша всех упаковок: %v", err)
		}

		log.Printf("Упаковка с ID %d обновлена", packaging.PackagingID)
	})
}

// DeletePackaging удаляет упаковку и сбрасывает кэш
func (s *PackagingService) DeletePackaging(ctx context.Context, packagingID int) error {
	ctx, span := s.tracer.Start(ctx, "DeletePackaging")
	defer span.End()

	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		exists, err := dao.CheckPackagingExists(ctx, packagingID, s.pool)
		if err != nil {
			errCh <- fmt.Errorf("ошибка проверки существования упаковки с ID %d: %w", packagingID, err)
			return
		}
		if !exists {
			errCh <- fmt.Errorf("упаковка с ID %d не найдена", packagingID)
			return
		}

		if err := dao.DeletePackaging(ctx, packagingID, s.pool); err != nil {
			errCh <- fmt.Errorf("ошибка удаления упаковки с ID %d: %w", packagingID, err)
			return
		}

		cacheKey := fmt.Sprintf("packaging_%d", packagingID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка удаления кэша упаковки с ключом %s: %v", cacheKey, err)
		}

		if err := s.cache.Delete(ctx, "all_packaging"); err != nil {
			log.Printf("Ошибка удаления кэша всех упаковок: %v", err)
		}

		errCh <- nil
	})

	return <-errCh
}

// GetPackagingIDByName возвращает ID упаковки по её имени через worker pool
func (s *PackagingService) GetPackagingIDByName(ctx context.Context, packagingType string) (int, error) {
	ctx, span := s.tracer.Start(ctx, "GetPackagingIDByName")
	defer span.End()

	packaging, err := dao.GetPackagingByName(ctx, packagingType, s.pool)
	if err != nil {
		return 0, fmt.Errorf("ошибка получения ID упаковки по названию: %w", err)
	}

	log.Printf("Упаковка с именем %s имеет ID %d", packagingType, packaging.PackagingID)
	return packaging.PackagingID, nil
}

// CheckPackagingExists проверяет существование упаковки через worker pool
func (s *PackagingService) CheckPackagingExists(ctx context.Context, packagingID int) (bool, error) {
	ctx, span := s.tracer.Start(ctx, "CheckPackagingExists")
	defer span.End()

	exists, err := dao.CheckPackagingExists(ctx, packagingID, s.pool)
	if err != nil {
		return false, fmt.Errorf("ошибка проверки существования упаковки с ID %d: %w", packagingID, err)
	}

	log.Printf("Упаковка с ID %d существует: %v", packagingID, exists)
	return exists, nil
}
