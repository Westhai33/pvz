package dao

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"homework1/internal/model"
)

// CreateUser создает нового пользователя с уровнем изоляции Read Committed.
func CreateUser(ctx context.Context, user model.User, pool *pgxpool.Pool) (int, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return 0, err
	}
	defer conn.Release() // Освобождаем соединение обратно в пул.

	var userID int
	err = tx.QueryRow(ctx,
		`INSERT INTO users (username, created_at) 
         VALUES ($1, $2) RETURNING user_id`,
		user.Username, user.CreatedAt).
		Scan(&userID)
	if err != nil {
		_ = tm.RollbackTransaction(ctx, tx, conn)
		return 0, fmt.Errorf("ошибка создания пользователя: %w", err)
	}

	if err := tm.CommitTransaction(ctx, tx, conn); err != nil {
		return 0, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return userID, nil
}

// GetUserByID возвращает пользователя по его ID с уровнем изоляции Repeatable Read.
func GetUserByID(ctx context.Context, userID int, pool *pgxpool.Pool) (*model.User, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.RepeatableRead)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var user model.User
	err = tx.QueryRow(ctx,
		`SELECT user_id, username, created_at 
         FROM users WHERE user_id = $1`, userID).
		Scan(&user.UserID, &user.Username, &user.CreatedAt)
	if err != nil {
		_ = tm.RollbackTransaction(ctx, tx, conn)
		return nil, fmt.Errorf("ошибка получения пользователя с ID %d: %w", userID, err)
	}

	if err := tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return &user, nil
}

// GetAllUsers возвращает всех пользователей с уровнем изоляции Serializable.
func GetAllUsers(ctx context.Context, pool *pgxpool.Pool) ([]model.User, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	query := `SELECT user_id, username, created_at FROM users`

	rows, err := tx.Query(ctx, query)
	if err != nil {
		_ = tm.RollbackTransaction(ctx, tx, conn)
		return nil, fmt.Errorf("ошибка получения пользователей: %w", err)
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.CreatedAt); err != nil {
			_ = tm.RollbackTransaction(ctx, tx, conn)
			return nil, fmt.Errorf("ошибка сканирования пользователя: %w", err)
		}
		users = append(users, user)
	}

	if err := tm.CommitTransaction(ctx, tx, conn); err != nil {
		return nil, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return users, nil
}

// UpdateUser обновляет данные пользователя с уровнем изоляции Read Committed.
func UpdateUser(ctx context.Context, user model.User, pool *pgxpool.Pool) error {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = tx.Exec(ctx,
		`UPDATE users SET username = $2, created_at = $3 
         WHERE user_id = $1`,
		user.UserID, user.Username, user.CreatedAt)
	if err != nil {
		_ = tm.RollbackTransaction(ctx, tx, conn)
		return fmt.Errorf("ошибка обновления пользователя с ID %d: %w", user.UserID, err)
	}

	if err := tm.CommitTransaction(ctx, tx, conn); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}

// DeleteUser удаляет пользователя по его ID с уровнем изоляции Serializable.
func DeleteUser(ctx context.Context, userID int, pool *pgxpool.Pool) error {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.Serializable)
	if err != nil {
		return err
	}
	defer conn.Release()

	query := `DELETE FROM users WHERE user_id = $1`

	_, err = tx.Exec(ctx, query, userID)
	if err != nil {
		_ = tm.RollbackTransaction(ctx, tx, conn)
		return fmt.Errorf("ошибка удаления пользователя с ID %d: %w", userID, err)
	}

	if err := tm.CommitTransaction(ctx, tx, conn); err != nil {
		return fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return nil
}

// CheckUserExists проверяет, существует ли пользователь по его ID с уровнем изоляции Read Committed.
func CheckUserExists(ctx context.Context, userID int, pool *pgxpool.Pool) (bool, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.ReadCommitted)
	if err != nil {
		return false, err
	}
	defer conn.Release()

	var exists bool
	err = tx.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)`, userID).
		Scan(&exists)
	if err != nil {
		_ = tm.RollbackTransaction(ctx, tx, conn)
		return false, fmt.Errorf("ошибка проверки существования пользователя: %w", err)
	}

	if err := tm.CommitTransaction(ctx, tx, conn); err != nil {
		return false, fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return exists, nil
}

// GetUserNameByID возвращает имя пользователя по его ID с уровнем изоляции Repeatable Read.
func GetUserNameByID(ctx context.Context, userID int, pool *pgxpool.Pool) (string, error) {
	tm := NewTransactionManager(pool)

	tx, conn, err := tm.BeginTransaction(ctx, pgx.RepeatableRead)
	if err != nil {
		return "", err
	}
	defer conn.Release()

	var username string
	err = tx.QueryRow(ctx,
		`SELECT username FROM users WHERE user_id = $1`, userID).
		Scan(&username)
	if err != nil {
		_ = tm.RollbackTransaction(ctx, tx, conn)
		return "", fmt.Errorf("ошибка получения имени пользователя с ID %d: %w", userID, err)
	}

	if err := tm.CommitTransaction(ctx, tx, conn); err != nil {
		return "", fmt.Errorf("ошибка подтверждения транзакции: %w", err)
	}

	return username, nil
}
