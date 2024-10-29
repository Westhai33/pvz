package dao

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"homework1/internal/model"
	"log"
	"time"
)

// CreateOrder создает новый заказ с уровнем изоляции Read Committed
func CreateOrder(ctx context.Context, order model.Order, pool *pgxpool.Pool) (int, error) {
	tm := NewTransactionManager(pool)

	// Начало транзакции
	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
				log.Printf("Ошибка отката транзакции: %v", rollbackErr)
			}
		}
	}()

	// Выполнение SQL-запроса на вставку нового заказа
	query := `INSERT INTO orders (user_id, acceptance_date, expiration_date, weight, base_cost, packaging_cost, total_cost, packaging_id, status_id, issue_date, with_film) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING order_id;`

	var orderID int
	err = tx.QueryRow(ctx, query,
		order.UserID, order.AcceptanceDate, order.ExpirationDate, order.Weight,
		order.BaseCost, order.PackagingCost, order.TotalCost, order.PackagingID,
		order.StatusID, order.IssueDate, order.WithFilm).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания заказа: %w", err)
	}

	// Подтверждение транзакции
	if err := tm.CommitTransaction(ctx, tx, conn); err != nil {
		return 0, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return orderID, nil
}

// GetOrderByID возвращает заказ по его ID с уровнем изоляции Repeatable Read
func GetOrderByID(ctx context.Context, orderID int, pool *pgxpool.Pool) (*model.Order, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.RepeatableRead)
	if err != nil {
		return nil, err
	}

	var order model.Order
	err = tx.QueryRow(ctx,
		`SELECT order_id, user_id, acceptance_date, expiration_date, weight, base_cost, packaging_cost, total_cost, packaging_id, status_id, issue_date, with_film 
		 FROM orders WHERE order_id = $1`, orderID).
		Scan(&order.OrderID, &order.UserID, &order.AcceptanceDate, &order.ExpirationDate, &order.Weight, &order.BaseCost, &order.PackagingCost, &order.TotalCost, &order.PackagingID, &order.StatusID, &order.IssueDate, &order.WithFilm)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка получения заказа с ID %d: %w", orderID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return &order, nil
}

// UpdateOrder обновляет заказ с уровнем изоляции Repeatable Read
func UpdateOrder(ctx context.Context, order model.Order, pool *pgxpool.Pool) error {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.RepeatableRead)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`UPDATE orders SET user_id = $2, acceptance_date = $3, expiration_date = $4, weight = $5, base_cost = $6, packaging_cost = $7, total_cost = $8, packaging_id = $9, status_id = $10, issue_date = $11, with_film = $12
		 WHERE order_id = $1`,
		order.OrderID, order.UserID, order.AcceptanceDate, order.ExpirationDate, order.Weight, order.BaseCost, order.PackagingCost, order.TotalCost, order.PackagingID, order.StatusID, order.IssueDate, order.WithFilm)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return fmt.Errorf("ошибка обновления заказа с ID %d: %w", order.OrderID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}

// DeleteOrder удаляет заказ по его ID с уровнем изоляции Serializable
func DeleteOrder(ctx context.Context, orderID int, pool *pgxpool.Pool) error {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM orders WHERE order_id = $1`, orderID)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return fmt.Errorf("ошибка удаления заказа с ID %d: %w", orderID, err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}

// GetAllOrders возвращает список всех заказов с уровнем изоляции Serializable
func GetAllOrders(ctx context.Context, pool *pgxpool.Pool) ([]model.Order, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT order_id, user_id, acceptance_date, expiration_date, weight, base_cost, packaging_cost, total_cost, packaging_id, status_id, issue_date, with_film 
		FROM orders
	`

	rows, err := tx.Query(ctx, query)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка получения всех заказов: %w", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(&order.OrderID, &order.UserID, &order.AcceptanceDate, &order.ExpirationDate, &order.Weight, &order.BaseCost, &order.PackagingCost, &order.TotalCost, &order.PackagingID, &order.StatusID, &order.IssueDate, &order.WithFilm)
		if err != nil {
			if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
				log.Printf("Ошибка отката транзакции: %v", rollbackErr)
			}
			return nil, fmt.Errorf("ошибка сканирования заказа: %w", err)
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка итерации по строкам заказов: %w", err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return orders, nil
}

// GetOrdersByUserID возвращает заказы для конкретного пользователя с уровнем изоляции Repeatable Read
func GetOrdersByUserID(ctx context.Context, userID int, pool *pgxpool.Pool) ([]model.Order, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.RepeatableRead)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT order_id, user_id, acceptance_date, expiration_date, weight, base_cost, packaging_cost, total_cost, packaging_id, status_id, issue_date, with_film 
		FROM orders WHERE user_id = $1
	`

	rows, err := tx.Query(ctx, query, userID)
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка получения заказов для пользователя с ID %d: %w", userID, err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(&order.OrderID, &order.UserID, &order.AcceptanceDate, &order.ExpirationDate, &order.Weight, &order.BaseCost, &order.PackagingCost, &order.TotalCost, &order.PackagingID, &order.StatusID, &order.IssueDate, &order.WithFilm)
		if err != nil {
			if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
				log.Printf("Ошибка отката транзакции: %v", rollbackErr)
			}
			return nil, fmt.Errorf("ошибка сканирования заказа: %w", err)
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка итерации по строкам заказов: %w", err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return orders, nil
}

// GetExpiredOrders возвращает заказы с истекшим сроком хранения с уровнем изоляции Serializable
func GetExpiredOrders(ctx context.Context, pool *pgxpool.Pool) ([]model.Order, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx,
		`SELECT order_id, user_id, acceptance_date, expiration_date, weight, base_cost, packaging_cost, total_cost, packaging_id, status_id, issue_date, with_film 
		 FROM orders WHERE expiration_date < $1`, time.Now())
	if err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка получения истекших заказов: %w", err)
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(&order.OrderID, &order.UserID, &order.AcceptanceDate, &order.ExpirationDate, &order.Weight, &order.BaseCost, &order.PackagingCost, &order.TotalCost, &order.PackagingID, &order.StatusID, &order.IssueDate, &order.WithFilm)
		if err != nil {
			if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
				log.Printf("Ошибка отката транзакции: %v", rollbackErr)
			}
			return nil, fmt.Errorf("ошибка сканирования заказа: %w", err)
		}
		orders = append(orders, order)
	}
	if err = rows.Err(); err != nil {
		if rollbackErr := tm.RollbackTransaction(ctx, tx, conn); rollbackErr != nil {
			log.Printf("Ошибка отката транзакции: %v", rollbackErr)
		}
		return nil, fmt.Errorf("ошибка итерации по строкам заказов: %w", err)
	}

	if err = tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return orders, nil
}
