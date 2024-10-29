package model

import "time"

type Return struct {
	ReturnID      int       `json:"return_id"`      // Уникальный идентификатор возврата
	OrderID       int       `json:"order_id"`       // Идентификатор заказа
	UserID        int       `json:"user_id"`        // Идентификатор пользователя
	ReturnDate    time.Time `json:"return_date"`    // Дата возврата
	ReasonID      int       `json:"reason_id"`      // Идентификатор причины возврата
	BaseCost      float64   `json:"base_cost"`      // Базовая стоимость
	PackagingCost float64   `json:"packaging_cost"` // Стоимость упаковки
	PackagingID   int       `json:"packaging_id"`   // Внешний ключ на упаковку
	TotalCost     float64   `json:"total_cost"`     // Общая стоимость
	StatusID      int       `json:"status_id"`      // Статус возврата
}
