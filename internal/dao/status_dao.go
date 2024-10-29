package dao

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"homework1/internal/model"
)

// CreateStatus создает новый статус с уровнем изоляции Read Committed
func CreateStatus(ctx context.Context, status model.Status, pool *pgxpool.Pool) (int, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	var statusID int
	err = tx.QueryRow(ctx,
		`INSERT INTO statuses (status_name) VALUES ($1) RETURNING status_id`,
		status.StatusName).Scan(&statusID)
	if err != nil {
		_ = tm.RollbackTransaction(ctx, tx, conn)
		return 0, fmt.Errorf("ошибка создания статуса: %w", err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return 0, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return statusID, nil
}

// GetStatusByID возвращает статус по его ID с уровнем изоляции Repeatable Read
func GetStatusByID(ctx context.Context, statusID int, pool *pgxpool.Pool) (*model.Status, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.RepeatableRead)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var status model.Status
	err = tx.QueryRow(ctx,
		`SELECT status_id, status_name FROM statuses WHERE status_id = $1`, statusID).
		Scan(&status.StatusID, &status.StatusName)
	if err != nil {
		_ = tm.RollbackTransaction(ctx, tx, conn)
		return nil, fmt.Errorf("ошибка получения статуса с ID %d: %w", statusID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return &status, nil
}

// GetAllStatuses возвращает все статусы с уровнем изоляции Serializable
func GetAllStatuses(ctx context.Context, pool *pgxpool.Pool) ([]model.Status, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := tx.Query(ctx, `SELECT status_id, status_name FROM statuses`)
	if err != nil {
		_ = tm.RollbackTransaction(ctx, tx, conn)
		return nil, fmt.Errorf("ошибка получения статусов: %w", err)
	}
	defer rows.Close()

	var statuses []model.Status
	for rows.Next() {
		var status model.Status
		if err := rows.Scan(&status.StatusID, &status.StatusName); err != nil {
			_ = tm.RollbackTransaction(ctx, tx, conn)
			return nil, fmt.Errorf("ошибка сканирования статуса: %w", err)
		}
		statuses = append(statuses, status)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return statuses, nil
}

// UpdateStatus обновляет статус с уровнем изоляции Read Committed
func UpdateStatus(ctx context.Context, status model.Status, pool *pgxpool.Pool) error {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = tx.Exec(ctx,
		`UPDATE statuses SET status_name = $2 WHERE status_id = $1`,
		status.StatusID, status.StatusName)
	if err != nil {
		_ = tm.RollbackTransaction(ctx, tx, conn)
		return fmt.Errorf("ошибка обновления статуса с ID %d: %w", status.StatusID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}

// DeleteStatus удаляет статус по его ID с уровнем изоляции Serializable
func DeleteStatus(ctx context.Context, statusID int, pool *pgxpool.Pool) error {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = tx.Exec(ctx, `DELETE FROM statuses WHERE status_id = $1`, statusID)
	if err != nil {
		_ = tm.RollbackTransaction(ctx, tx, conn)
		return fmt.Errorf("ошибка удаления статуса с ID %d: %w", statusID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}

// CheckStatusExists проверяет, существует ли статус по его ID с уровнем изоляции Read Committed
func CheckStatusExists(ctx context.Context, statusID int, pool *pgxpool.Pool) (bool, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	var exists bool
	err = tx.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM statuses WHERE status_id = $1)`, statusID).
		Scan(&exists)
	if err != nil {
		_ = tm.RollbackTransaction(ctx, tx, conn)
		return false, fmt.Errorf("ошибка проверки существования статуса: %w", err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return false, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return exists, nil
}

// GetStatusByName возвращает статус по его имени с уровнем изоляции Serializable
func GetStatusByName(ctx context.Context, statusName string, pool *pgxpool.Pool) (*model.Status, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var status model.Status
	err = tx.QueryRow(ctx,
		`SELECT status_id, status_name FROM statuses WHERE status_name = $1`, statusName).
		Scan(&status.StatusID, &status.StatusName)
	if err != nil {
		_ = tm.RollbackTransaction(ctx, tx, conn)
		return nil, fmt.Errorf("ошибка получения статуса с именем %s: %w", statusName, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return &status, nil
}

// GetStatusNameByID возвращает имя статуса по его ID с уровнем изоляции Repeatable Read
func GetStatusNameByID(ctx context.Context, statusID int, pool *pgxpool.Pool) (string, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.RepeatableRead)
	if err != nil {
		return "", err
	}
	defer conn.Release()

	var statusName string
	err = tx.QueryRow(ctx,
		`SELECT status_name FROM statuses WHERE status_id = $1`, statusID).
		Scan(&statusName)
	if err != nil {
		_ = tm.RollbackTransaction(ctx, tx, conn)
		return "", fmt.Errorf("ошибка получения имени статуса с ID %d: %w", statusID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return "", fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return statusName, nil
}
