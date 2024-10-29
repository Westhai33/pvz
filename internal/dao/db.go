package dao

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

var pool *pgxpool.Pool

// InitDB инициализирует пул соединений к базе данных PostgreSQL
func Initdb(user, password, dbname, host string, port int) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Ошибка парсинга конфигурации подключения: %v\n", err)
	}

	config.MaxConns = 20
	config.MaxConnIdleTime = 5 * time.Minute
	config.ConnConfig.ConnectTimeout = 10 * time.Second

	pool, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v\n", err)
	} else {
		log.Println("Успешное подключение к базе данных с использованием пула.")
	}
}

// CloseDB закрывает пул соединений
func Closedb() {
	if pool != nil {
		pool.Close()
		log.Println("Пул соединений закрыт.")
	}
}

// GetPool возвращает пул соединений для использования в других местах программы
func GetPool() *pgxpool.Pool {
	return pool
}
