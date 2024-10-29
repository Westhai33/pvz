package dao

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"homework1/internal/model"
)

// CreatePackaging создает новую упаковку с уровнем изоляции Read Committed
func CreatePackaging(ctx context.Context, packaging model.PackagingOption, pool *pgxpool.Pool) (int, error) {
	tm := NewTransactionManager(pool)

	// Начинаем транзакцию с уровнем изоляции Read Committed
	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return 0, err
	}
	defer conn.Release() // Освобождаем соединение обратно в пул

	var packagingID int
	query := `INSERT INTO packaging (type, cost, max_weight) VALUES ($1, $2, $3) RETURNING packaging_id`

	err = tx.QueryRow(ctx, query, packaging.Type, packaging.Cost, packaging.MaxWeight).Scan(&packagingID)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return 0, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return 0, fmt.Errorf("ошибка создания упаковки: %w", err)
	}

	// Подтверждаем транзакцию
	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return 0, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return packagingID, nil
}

// GetAllPackaging возвращает все упаковки с уровнем изоляции Repeatable Read
func GetAllPackaging(ctx context.Context, pool *pgxpool.Pool) ([]model.PackagingOption, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.RepeatableRead)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := tx.Query(ctx, "SELECT packaging_id, type, cost, max_weight FROM packaging")
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return nil, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка получения упаковок: %w", err)
	}
	defer rows.Close()

	var packagingOptions []model.PackagingOption
	for rows.Next() {
		var packaging model.PackagingOption
		if err := rows.Scan(&packaging.PackagingID, &packaging.Type, &packaging.Cost, &packaging.MaxWeight); err != nil {
			if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
				return nil, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
			}
			return nil, fmt.Errorf("ошибка сканирования упаковки: %w", err)
		}
		packagingOptions = append(packagingOptions, packaging)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return packagingOptions, nil
}

// GetPackagingByID возвращает упаковку по ID с уровнем изоляции Serializable
func GetPackagingByID(ctx context.Context, packagingID int, pool *pgxpool.Pool) (*model.PackagingOption, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var packaging model.PackagingOption
	err = tx.QueryRow(ctx, "SELECT packaging_id, type, cost, max_weight FROM packaging WHERE packaging_id = $1", packagingID).
		Scan(&packaging.PackagingID, &packaging.Type, &packaging.Cost, &packaging.MaxWeight)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return nil, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка получения упаковки с ID %d: %w", packagingID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return &packaging, nil
}

// UpdatePackaging обновляет упаковку с уровнем изоляции Read Committed
func UpdatePackaging(ctx context.Context, packaging model.PackagingOption, pool *pgxpool.Pool) error {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = tx.Exec(ctx, `UPDATE packaging SET type = $2, cost = $3, max_weight = $4 WHERE packaging_id = $1`,
		packaging.PackagingID, packaging.Type, packaging.Cost, packaging.MaxWeight)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return fmt.Errorf("ошибка обновления упаковки с ID %d: %w", packaging.PackagingID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}

// DeletePackaging удаляет упаковку по ID с уровнем изоляции Serializable
func DeletePackaging(ctx context.Context, packagingID int, pool *pgxpool.Pool) error {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = tx.Exec(ctx, "DELETE FROM packaging WHERE packaging_id = $1", packagingID)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return fmt.Errorf("ошибка удаления упаковки с ID %d: %w", packagingID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}

// CheckPackagingExists проверяет, существует ли упаковка по ее ID с уровнем изоляции Read Committed
func CheckPackagingExists(ctx context.Context, packagingID int, pool *pgxpool.Pool) (bool, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	var exists bool
	err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM packaging WHERE packaging_id = $1)`, packagingID).Scan(&exists)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return false, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return false, fmt.Errorf("ошибка проверки существования упаковки: %w", err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return false, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return exists, nil
}

// GetPackagingByName возвращает упаковку по названию с уровнем изоляции Serializable
func GetPackagingByName(ctx context.Context, packagingType string, pool *pgxpool.Pool) (*model.PackagingOption, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var packaging model.PackagingOption
	err = tx.QueryRow(ctx, "SELECT packaging_id, type, cost, max_weight FROM packaging WHERE type = $1", packagingType).
		Scan(&packaging.PackagingID, &packaging.Type, &packaging.Cost, &packaging.MaxWeight)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			return nil, fmt.Errorf("ошибка отката транзакции: %w", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка получения упаковки по названию %s: %w", packagingType, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return &packaging, nil
}
