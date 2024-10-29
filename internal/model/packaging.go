package model

// PackagingFactory - интерфейс для фабрики упаковок
type PackagingFactory interface {
	CreatePackaging() PackagingOption
	CreateAdditionalOption() PackagingOption
}

// PackagingOption представляет параметры упаковки
type PackagingOption struct {
	PackagingID int     `json:"packaging_id"` // Уникальный идентификатор упаковки
	Type        string  `json:"type"`         // Тип упаковки
	Cost        float64 `json:"cost"`         // Стоимость
	MaxWeight   float64 `json:"max_weight"`   // Максимальный вес
}

// DynamicPackagingFactory - динамическая фабрика для создания упаковки
type DynamicPackagingFactory struct {
	Type      string
	Cost      float64
	MaxWeight float64
}

// CreatePackaging создает упаковку с параметрами, переданными при регистрации фабрики
func (f *DynamicPackagingFactory) CreatePackaging() PackagingOption {
	return PackagingOption{
		Type:      f.Type,
		Cost:      f.Cost,
		MaxWeight: f.MaxWeight,
	}
}

// CreateAdditionalOption создает дополнительную опцию (например, пленку)
func (f *DynamicPackagingFactory) CreateAdditionalOption() PackagingOption {
	return PackagingOption{
		Type:      "film",
		Cost:      1.0,
		MaxWeight: 0.0,
	}
}
