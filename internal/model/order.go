package model

import "time"

// Order представляет заказ в системе
type Order struct {
	OrderID        int       `json:"order_id"`        // Уникальный идентификатор заказа
	UserID         int       `json:"user_id"`         // Внешний ключ на пользователя
	AcceptanceDate time.Time `json:"acceptance_date"` // Дата принятия заказа
	ExpirationDate time.Time `json:"expiration_date"` // Дата истечения срока хранения
	Weight         float64   `json:"weight"`          // Вес заказа
	BaseCost       float64   `json:"base_cost"`       // Базовая стоимость заказа
	PackagingCost  float64   `json:"packaging_cost"`  // Стоимость упаковки
	TotalCost      float64   `json:"total_cost"`      // Общая стоимость заказа
	PackagingID    int       `json:"packaging_id"`    // Внешний ключ на таблицу упаковки
	StatusID       int       `json:"status_id"`       // Внешний ключ на таблицу статусов
	IssueDate      time.Time `json:"issue_date"`      // Дата выдачи заказа (может быть NULL)
	WithFilm       bool      `json:"with_film"`       // Указание на упаковку с пленкой или без
}
