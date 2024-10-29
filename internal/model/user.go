package model

import "time"

type User struct {
	UserID    int       `json:"user_id" validate:"required"`    // Идентификатор пользователя (целое число)
	Username  string    `json:"username" validate:"required"`   // Имя пользователя
	CreatedAt time.Time `json:"created_at" validate:"required"` // Время создания пользователя
}
