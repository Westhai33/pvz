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

// UserService представляет сервис для работы с пользователями
type UserService struct {
	pool   *pgxpool.Pool
	wp     *pool.WorkerPool
	cache  *cache.RedisCache[string, model.User]
	tracer trace.Tracer
}

// NewUserService создает новый сервис для работы с пользователями
func NewUserService(dbPool *pgxpool.Pool, workerPool *pool.WorkerPool, cache *cache.RedisCache[string, model.User]) *UserService {
	return &UserService{
		pool:   dbPool,
		wp:     workerPool,
		cache:  cache,
		tracer: tracing.GetTracer(),
	}
}

// CreateUser создает нового пользователя через общий worker pool
func (s *UserService) CreateUser(ctx context.Context, username string) (int, error) {
	ctx, span := s.tracer.Start(ctx, "CreateUser")
	defer span.End()

	var userID int
	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		newUser := model.User{
			Username:  username,
			CreatedAt: time.Now().UTC(),
		}

		var err error
		userID, err = dao.CreateUser(ctx, newUser, s.pool)
		if err != nil {
			errCh <- fmt.Errorf("ошибка создания пользователя: %w", err)
			return
		}

		cacheKey := fmt.Sprintf("user_%d", userID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка удаления кэша пользователя с ключом %s: %v", cacheKey, err)
		}

		errCh <- nil
	})

	err := <-errCh
	return userID, err
}

// GetUserByID возвращает пользователя по его ID с использованием кэша
func (s *UserService) GetUserByID(ctx context.Context, userID int) (*model.User, error) {
	ctx, span := s.tracer.Start(ctx, "GetUserByID")
	defer span.End()

	cacheKey := fmt.Sprintf("user_%d", userID)

	var cachedUser model.User
	err := s.cache.Get(ctx, cacheKey, &cachedUser)
	if err == nil {
		log.Printf("Пользователь с ID %d получен из кэша", userID)
		return &cachedUser, nil
	} else {
		log.Printf("Пользователь с ID %d не найден в кэше, получаем из базы данных", userID)
	}

	user, err := dao.GetUserByID(ctx, userID, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователя с ID %d: %w", userID, err)
	}

	err = s.cache.Set(ctx, cacheKey, *user, 10*time.Minute) // Разыменовываем указатель
	if err != nil {
		log.Printf("Ошибка сохранения пользователя в кэш: %v", err)
	}

	return user, nil
}

// GetAllUsers возвращает всех пользователей через worker pool
func (s *UserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	ctx, span := s.tracer.Start(ctx, "GetAllUsers")
	defer span.End()

	users, err := dao.GetAllUsers(ctx, s.pool)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения всех пользователей: %w", err)
	}

	log.Printf("Получены все пользователи: %+v", users)
	return users, nil
}

// UpdateUser обновляет данные пользователя и сбрасывает кэш
func (s *UserService) UpdateUser(ctx context.Context, userID int, username string) error {
	ctx, span := s.tracer.Start(ctx, "UpdateUser")
	defer span.End()

	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		exists, err := dao.CheckUserExists(ctx, userID, s.pool)
		if err != nil {
			errCh <- fmt.Errorf("ошибка проверки существования пользователя с ID %d: %w", userID, err)
			return
		}
		if !exists {
			errCh <- fmt.Errorf("пользователь с ID %d не найден", userID)
			return
		}

		updatedUser := model.User{
			UserID:    userID,
			Username:  username,
			CreatedAt: time.Now().UTC(),
		}

		if err := dao.UpdateUser(ctx, updatedUser, s.pool); err != nil {
			errCh <- fmt.Errorf("ошибка обновления данных пользователя с ID %d: %w", userID, err)
			return
		}

		cacheKey := fmt.Sprintf("user_%d", userID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка удаления кэша пользователя с ключом %s: %v", cacheKey, err)
		}

		errCh <- nil
	})

	return <-errCh
}

// DeleteUser удаляет пользователя и сбрасывает кэш
func (s *UserService) DeleteUser(ctx context.Context, userID int) error {
	ctx, span := s.tracer.Start(ctx, "DeleteUser")
	defer span.End()

	errCh := make(chan error, 1)

	s.wp.SubmitTask(func() {
		exists, err := dao.CheckUserExists(ctx, userID, s.pool)
		if err != nil {
			errCh <- fmt.Errorf("ошибка проверки существования пользователя с ID %d: %w", userID, err)
			return
		}
		if !exists {
			errCh <- fmt.Errorf("пользователь с ID %d не найден", userID)
			return
		}

		if err := dao.DeleteUser(ctx, userID, s.pool); err != nil {
			errCh <- fmt.Errorf("ошибка удаления пользователя с ID %d: %w", userID, err)
			return
		}

		cacheKey := fmt.Sprintf("user_%d", userID)
		if err := s.cache.Delete(ctx, cacheKey); err != nil {
			log.Printf("Ошибка удаления кэша пользователя с ключом %s: %v", cacheKey, err)
		}

		errCh <- nil
	})

	return <-errCh
}

// GetUserNameByID возвращает имя пользователя по его ID через worker pool
func (s *UserService) GetUserNameByID(ctx context.Context, userID int) (string, error) {
	ctx, span := s.tracer.Start(ctx, "GetUserNameByID")
	defer span.End()

	cacheKey := fmt.Sprintf("user_name_%d", userID)

	var cachedUsername string
	err := s.cache.Get(ctx, cacheKey, &cachedUsername)
	if err == nil {
		log.Printf("Имя пользователя с ID %d получено из кэша", userID)
		return cachedUsername, nil
	} else {
		log.Printf("Имя пользователя с ID %d не найдено в кэше, получаем из базы данных", userID)
	}

	username, err := dao.GetUserNameByID(ctx, userID, s.pool)
	if err != nil {
		return "", fmt.Errorf("ошибка получения имени пользователя с ID %d: %w", userID, err)
	}

	err = s.cache.SetString(ctx, cacheKey, username, 10*time.Minute)
	if err != nil {
		log.Printf("Ошибка сохранения имени пользователя в кэш: %v", err)
	}

	return username, nil
}

// CheckUserExists проверяет существование пользователя через worker pool
func (s *UserService) CheckUserExists(ctx context.Context, userID int) (bool, error) {
	ctx, span := s.tracer.Start(ctx, "CheckUserExists")
	defer span.End()

	exists, err := dao.CheckUserExists(ctx, userID, s.pool)
	if err != nil {
		return false, fmt.Errorf("ошибка проверки существования пользователя с ID %d: %w", userID, err)
	}

	log.Printf("Пользователь с ID %d существует: %v", userID, exists)
	return exists, nil
}
