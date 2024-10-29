package dao

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// TransactionManager управляет транзакциями через пул соединений.
type TransactionManager struct {
	pool *pgxpool.Pool
}

// NewTransactionManager создает новый менеджер транзакций с использованием пула соединений.
func NewTransactionManager(pool *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{pool: pool}
}

// BeginTransaction начинает новую транзакцию с заданным уровнем изоляции.
func (tm *TransactionManager) BeginTransaction(ctx context.Context, isoLevel pgx.TxIsoLevel) (pgx.Tx, *pgxpool.Conn, error) {
	conn, err := tm.pool.Acquire(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("ошибка получения соединения из пула: %w", err)
	}

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   isoLevel,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		conn.Release()
		return nil, nil, fmt.Errorf("ошибка начала транзакции: %w", err)
	}

	return tx, conn, nil
}

// CommitTransaction подтверждает транзакцию и освобождает соединение обратно в пул.
func (tm *TransactionManager) CommitTransaction(ctx context.Context, tx pgx.Tx, conn *pgxpool.Conn) error {
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}
	conn.Release()
	return nil
}

// RollbackTransaction откатывает транзакцию и освобождает соединение обратно в пул.
func (tm *TransactionManager) RollbackTransaction(ctx context.Context, tx pgx.Tx, conn *pgxpool.Conn) error {
	if err := tx.Rollback(ctx); err != nil {
		return fmt.Errorf("ошибка отката транзакции: %w", err)
	}
	conn.Release()
	return nil
}
