package dao

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"homework1/internal/model"
	"time"
)

// ReadReturns читает все возвраты из базы данных с уровнем изоляции Read Committed
func ReadReturns(ctx context.Context, pool *pgxpool.Pool) ([]model.Return, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `SELECT return_id, order_id, user_id, return_date, reason_id, base_cost, packaging_cost, packaging_id, total_cost, status_id FROM returns`

	rows, err := tx.Query(ctx, query)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return nil, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка чтения возвратов: %w", err)
	}
	defer rows.Close()

	var returns []model.Return
	for rows.Next() {
		var ret model.Return
		err := rows.Scan(&ret.ReturnID, &ret.OrderID, &ret.UserID, &ret.ReturnDate, &ret.ReasonID, &ret.BaseCost, &ret.PackagingCost, &ret.PackagingID, &ret.TotalCost, &ret.StatusID)
		if err != nil {
			if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
				return nil, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
			}
			return nil, fmt.Errorf("ошибка обработки строки: %w", err)
		}
		returns = append(returns, ret)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return returns, nil
}

// WriteReturns записывает возвраты в базу данных с уровнем изоляции Read Committed
func WriteReturns(ctx context.Context, returns []model.Return, pool *pgxpool.Pool) error {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return err
	}
	defer conn.Release()

	for _, ret := range returns {
		query := `INSERT INTO returns (order_id, user_id, return_date, reason_id, base_cost, packaging_cost, packaging_id, total_cost, status_id)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
            ON CONFLICT (return_id) DO UPDATE
            SET order_id = $1, user_id = $2, return_date = $3, reason_id = $4, base_cost = $5, packaging_cost = $6, packaging_id = $7, total_cost = $8, status_id = $9`

		_, err := tx.Exec(ctx, query, ret.OrderID, ret.UserID, ret.ReturnDate, ret.ReasonID, ret.BaseCost, ret.PackagingCost, ret.PackagingID, ret.TotalCost, ret.StatusID)
		if err != nil {
			if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
				return fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
			}
			return fmt.Errorf("ошибка записи возврата: %w", err)
		}
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}

// FindReturnByOrderID ищет возврат по идентификатору заказа с уровнем изоляции Repeatable Read
func FindReturnByOrderID(ctx context.Context, orderID int, pool *pgxpool.Pool) (*model.Return, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.RepeatableRead)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `SELECT return_id, order_id, user_id, return_date, reason_id, base_cost, packaging_cost, packaging_id, total_cost, status_id FROM returns WHERE order_id = $1`

	row := tx.QueryRow(ctx, query, orderID)

	var ret model.Return
	err = row.Scan(&ret.ReturnID, &ret.OrderID, &ret.UserID, &ret.ReturnDate, &ret.ReasonID, &ret.BaseCost, &ret.PackagingCost, &ret.PackagingID, &ret.TotalCost, &ret.StatusID)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return nil, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return nil, fmt.Errorf("возврат с order_id %d не найден: %w", orderID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return &ret, nil
}

// DeleteReturn удаляет возврат из базы данных по идентификатору возврата с уровнем изоляции Serializable
func DeleteReturn(ctx context.Context, returnID int, pool *pgxpool.Pool) error {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := `DELETE FROM returns WHERE return_id = $1`

	_, err = tx.Exec(ctx, query, returnID)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return fmt.Errorf("ошибка удаления возврата с ID %d: %w", returnID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}

// FindReturnsByUserID ищет возвраты по идентификатору пользователя с уровнем изоляции Read Committed
func FindReturnsByUserID(ctx context.Context, userID int, pool *pgxpool.Pool) ([]model.Return, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `SELECT return_id, order_id, user_id, return_date, reason_id, base_cost, packaging_cost, packaging_id, total_cost, status_id FROM returns WHERE user_id = $1`

	rows, err := tx.Query(ctx, query, userID)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return nil, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка поиска возвратов для пользователя с ID %d: %w", userID, err)
	}
	defer rows.Close()

	var returns []model.Return
	for rows.Next() {
		var ret model.Return
		err := rows.Scan(&ret.ReturnID, &ret.OrderID, &ret.UserID, &ret.ReturnDate, &ret.ReasonID, &ret.BaseCost, &ret.PackagingCost, &ret.PackagingID, &ret.TotalCost, &ret.StatusID)
		if err != nil {
			if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
				return nil, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
			}
			return nil, fmt.Errorf("ошибка обработки строки: %w", err)
		}
		returns = append(returns, ret)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return returns, nil
}

// FindExpiredReturns ищет возвраты, которые не были завершены курьером в установленные сроки с уровнем изоляции Serializable
func FindExpiredReturns(ctx context.Context, pool *pgxpool.Pool) ([]model.Return, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `SELECT return_id, order_id, user_id, return_date, reason_id, base_cost, packaging_cost, packaging_id, total_cost, status_id 
              FROM returns WHERE return_date < $1 AND status_id != $2`

	rows, err := tx.Query(ctx, query, time.Now(), 4) // Assuming status ID 4 means "Completed"
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return nil, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка поиска истекших возвратов: %w", err)
	}
	defer rows.Close()

	var expiredReturns []model.Return
	for rows.Next() {
		var ret model.Return
		err := rows.Scan(&ret.ReturnID, &ret.OrderID, &ret.UserID, &ret.ReturnDate, &ret.ReasonID, &ret.BaseCost, &ret.PackagingCost, &ret.PackagingID, &ret.TotalCost, &ret.StatusID)
		if err != nil {
			if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
				return nil, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
			}
			return nil, fmt.Errorf("ошибка обработки строки: %w", err)
		}
		expiredReturns = append(expiredReturns, ret)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return expiredReturns, nil
}

// UpdateReturn обновляет данные о возврате в базе данных с уровнем изоляции Read Committed
func UpdateReturn(ctx context.Context, ret model.Return, pool *pgxpool.Pool) error {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = tx.Exec(ctx,
		`UPDATE returns 
		SET order_id = $1, user_id = $2, return_date = $3, reason_id = $4, base_cost = $5, packaging_cost = $6, packaging_id = $7, total_cost = $8, status_id = $9
		WHERE return_id = $10`,
		ret.OrderID, ret.UserID, ret.ReturnDate, ret.ReasonID, ret.BaseCost, ret.PackagingCost, ret.PackagingID, ret.TotalCost, ret.StatusID, ret.ReturnID)
	if err != nil {
		_ = tm.RollbackTransaction(ctx, tx, conn)
		return fmt.Errorf("ошибка обновления возврата с ID %d: %w", ret.ReturnID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}

// CreateReturn создает новый возврат на основе данных заказа с уровнем изоляции Read Committed
func CreateReturn(ctx context.Context, newReturn model.Return, pool *pgxpool.Pool) error {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = tx.Exec(ctx,
		`INSERT INTO returns (order_id, user_id, return_date, reason_id, base_cost, packaging_cost, packaging_id, total_cost, status_id)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		newReturn.OrderID, newReturn.UserID, newReturn.ReturnDate, newReturn.ReasonID, newReturn.BaseCost, newReturn.PackagingCost, newReturn.PackagingID, newReturn.TotalCost, newReturn.StatusID)

	if err != nil {
		_ = tm.RollbackTransaction(ctx, tx, conn)
		return fmt.Errorf("ошибка вставки возврата в базу данных: %w", err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}

// CheckReturnExists проверяет, существует ли возврат с данным returnID
func CheckReturnExists(ctx context.Context, returnID int, pool *pgxpool.Pool) (bool, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM returns WHERE return_id = $1)`
	err = tx.QueryRow(ctx, query, returnID).Scan(&exists)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return false, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return false, fmt.Errorf("ошибка проверки существования возврата с ID %d: %w", returnID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return false, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return exists, nil
}
