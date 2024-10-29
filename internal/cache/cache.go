package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// CacheConfig содержит параметры конфигурации для кэша
type CacheConfig struct {
	DefaultTTL time.Duration
}

// RedisCache представляет структуру для работы с Redis с универсальными типами ключей и значений
type RedisCache[K comparable, V any] struct {
	client *redis.Client
	config CacheConfig
}

// NewRedisCache создает новый экземпляр кэша с Redis-клиентом
func NewRedisCache[K comparable, V any](client *redis.Client, config CacheConfig) *RedisCache[K, V] {
	return &RedisCache[K, V]{
		client: client,
		config: config,
	}
}

// Set сохраняет объект в Redis с временем жизни (TTL)
func (c *RedisCache[K, V]) Set(ctx context.Context, key K, value V, ttl ...time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных для кэша: %w", err)
	}

	expiration := c.config.DefaultTTL
	if len(ttl) > 0 {
		expiration = ttl[0]
	}

	err = c.client.Set(ctx, fmt.Sprintf("%v", key), data, expiration).Err()
	if err != nil {
		return fmt.Errorf("ошибка сохранения данных в Redis: %w", err)
	}

	return nil
}

// Get возвращает объект из Redis и десериализует его
func (c *RedisCache[K, V]) Get(ctx context.Context, key K, dest interface{}) error {
	val, err := c.client.Get(ctx, fmt.Sprintf("%v", key)).Result()
	if err == redis.Nil {
		return fmt.Errorf("данные не найдены в кэше по ключу: %v", key)
	} else if err != nil {
		return fmt.Errorf("ошибка при получении данных из Redis: %w", err)
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return fmt.Errorf("ошибка при десериализации данных из Redis: %w", err)
	}

	return nil
}

// Delete удаляет объект из Redis по ключу
func (c *RedisCache[K, V]) Delete(ctx context.Context, key K) error {
	err := c.client.Del(ctx, fmt.Sprintf("%v", key)).Err()
	if err != nil {
		return fmt.Errorf("ошибка удаления данных из Redis: %w", err)
	}

	return nil
}

// Exists проверяет наличие объекта в Redis по ключу
func (c *RedisCache[K, V]) Exists(ctx context.Context, key K) (bool, error) {
	count, err := c.client.Exists(ctx, fmt.Sprintf("%v", key)).Result()
	if err != nil {
		return false, fmt.Errorf("ошибка при проверке существования ключа в Redis: %w", err)
	}

	return count > 0, nil
}

// Keys возвращает все ключи по шаблону
func (c *RedisCache[K, V]) Keys(ctx context.Context, pattern string) ([]string, error) {
	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении ключей из Redis: %w", err)
	}
	return keys, nil
}

// SetSlice сохраняет срез объектов в Redis с временем жизни (TTL)
func (c *RedisCache[K, V]) SetSlice(ctx context.Context, key K, values []V, ttl ...time.Duration) error {
	data, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("ошибка сериализации данных для кэша: %w", err)
	}

	expiration := c.config.DefaultTTL
	if len(ttl) > 0 {
		expiration = ttl[0]
	}

	err = c.client.Set(ctx, fmt.Sprintf("%v", key), data, expiration).Err()
	if err != nil {
		return fmt.Errorf("ошибка сохранения данных в Redis: %w", err)
	}

	return nil
}

// GetSlice возвращает срез объектов из Redis и десериализует его
func (c *RedisCache[K, V]) GetSlice(ctx context.Context, key K, dest *[]V) error {
	val, err := c.client.Get(ctx, fmt.Sprintf("%v", key)).Result()
	if err == redis.Nil {
		return fmt.Errorf("данные не найдены в кэше по ключу: %v", key)
	} else if err != nil {
		return fmt.Errorf("ошибка при получении данных из Redis: %w", err)
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return fmt.Errorf("ошибка при десериализации данных из Redis: %w", err)
	}

	return nil
}

// SetString сохраняет строку в Redis с временем жизни (TTL)
func (c *RedisCache[K, V]) SetString(ctx context.Context, key K, value string, ttl ...time.Duration) error {
	expiration := c.config.DefaultTTL
	if len(ttl) > 0 {
		expiration = ttl[0]
	}

	err := c.client.Set(ctx, fmt.Sprintf("%v", key), value, expiration).Err()
	if err != nil {
		return fmt.Errorf("ошибка сохранения строки в Redis: %w", err)
	}
	return nil
}
