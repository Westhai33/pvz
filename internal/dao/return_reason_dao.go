package dao

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"homework1/internal/model"
)

// CreateReturnReason создает новую причину возврата с уровнем изоляции Read Committed
func CreateReturnReason(ctx context.Context, reason model.ReturnReason, pool *pgxpool.Pool) (int, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	var reasonID int
	query := `INSERT INTO return_reasons (reason) VALUES ($1) RETURNING reason_id`

	err = tx.QueryRow(ctx, query, reason.Reason).Scan(&reasonID)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return 0, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return 0, fmt.Errorf("ошибка создания причины возврата: %w", err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return 0, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return reasonID, nil
}

// GetReturnReasonByID возвращает причину возврата по её ID с уровнем изоляции Repeatable Read
func GetReturnReasonByID(ctx context.Context, reasonID int, pool *pgxpool.Pool) (*model.ReturnReason, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.RepeatableRead)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var reason model.ReturnReason
	err = tx.QueryRow(ctx,
		`SELECT reason_id, reason FROM return_reasons WHERE reason_id = $1`, reasonID).
		Scan(&reason.ReasonID, &reason.Reason)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return nil, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка получения причины возврата с ID %d: %w", reasonID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return &reason, nil
}

// GetAllReturnReasons возвращает все причины возвратов с уровнем изоляции Serializable
func GetAllReturnReasons(ctx context.Context, pool *pgxpool.Pool) ([]model.ReturnReason, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := tx.Query(ctx, `SELECT reason_id, reason FROM return_reasons`)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return nil, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка получения причин возвратов: %w", err)
	}
	defer rows.Close()

	var reasons []model.ReturnReason
	for rows.Next() {
		var reason model.ReturnReason
		if err := rows.Scan(&reason.ReasonID, &reason.Reason); err != nil {
			if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
				return nil, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
			}
			return nil, fmt.Errorf("ошибка сканирования причины возврата: %w", err)
		}
		reasons = append(reasons, reason)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return reasons, nil
}

// UpdateReturnReason обновляет существующую причину возврата с уровнем изоляции Read Committed
func UpdateReturnReason(ctx context.Context, reason model.ReturnReason, pool *pgxpool.Pool) error {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = tx.Exec(ctx, `UPDATE return_reasons SET reason = $2 WHERE reason_id = $1`, reason.ReasonID, reason.Reason)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return fmt.Errorf("ошибка обновления причины возврата с ID %d: %w", reason.ReasonID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}

// DeleteReturnReason удаляет причину возврата по её ID с уровнем изоляции Serializable
func DeleteReturnReason(ctx context.Context, reasonID int, pool *pgxpool.Pool) error {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = tx.Exec(ctx, `DELETE FROM return_reasons WHERE reason_id = $1`, reasonID)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return fmt.Errorf("ошибка удаления причины возврата с ID %d: %w", reasonID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}

// CheckReturnReasonExists проверяет, существует ли причина возврата по её ID с уровнем изоляции Read Committed
func CheckReturnReasonExists(ctx context.Context, reasonID int, pool *pgxpool.Pool) (bool, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	var exists bool
	err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM return_reasons WHERE reason_id = $1)`, reasonID).Scan(&exists)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return false, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return false, fmt.Errorf("ошибка проверки существования причины возврата: %w", err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return false, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return exists, nil
}

// GetReturnReasonByName возвращает причину возврата по её имени
func GetReturnReasonByName(ctx context.Context, reasonName string, pool *pgxpool.Pool) (*model.ReturnReason, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var reason model.ReturnReason
	err = tx.QueryRow(ctx, `SELECT reason_id, reason FROM return_reasons WHERE reason = $1`, reasonName).
		Scan(&reason.ReasonID, &reason.Reason)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return nil, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка получения причины возврата с именем %s: %w", reasonName, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return &reason, nil
}
