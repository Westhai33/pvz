package model

type Status struct {
	StatusID   int    `json:"status_id" validate:"required"`
	StatusName string `json:"status_name" validate:"required"`
}
