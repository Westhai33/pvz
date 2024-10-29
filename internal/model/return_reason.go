package model

type ReturnReason struct {
	ReasonID int    `json:"reason_id"` // Уникальный идентификатор причины возврата
	Reason   string `json:"reason"`    // Описание причины возврата
}
