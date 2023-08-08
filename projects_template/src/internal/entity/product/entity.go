package product

import (
	"errors"
	"product_storage/tools/sqlnull"
	"time"
)

// Product cтруктура продукта для записи в базу
type Product struct {
	ProductID   int              `json:"product_id"`  // id продукта
	Name        string           `json:"name"`        // название продукта
	Descr       string           `json:"description"` // описание продукта
	AddetAt     time.Time        `json:"added_at"`    // дата добавления продукта
	RemovedAt   sqlnull.NullTime `json:"removed_at"`  // дата удаления продукта
	Tags        string           `json:"tags"`        // теги продукта
	VariantList []Variant        `json:"variants"`    // cписок вариантов продукта
}

// Variant структура варианта, продукта представляем с собой информацию о продукте который нужно внести в базу
type Variant struct {
	ProductID    int       `json:"product_id" db:"product_id"` // id продука
	VariantID    int       `json:"variant_id" db:"variant_id"` // id конкретного варианта продукта
	Weight       int       `json:"weight" db:"weight"`         // масса или вес продукта
	Unit         string    `json:"unit" db:"unit"`             // единица измерения
	AddedAt      time.Time `json:"added_at" db:"added_at"`     // дата добавления определенного варианта
	CurrentPrice float64   `json:"price" db:"price"`           // актуальная цена
	InStorages   []int     `json:"in_storages"`                // список id складов в которых есть этот вариант
}

// ProductPrice структура для вставки цены продукта
type ProductPrice struct {
	PriceID   int              `json:"price_id" db:"price_id"`     // id цены продукта
	VariantID int              `json:"variant_id" db:"variant_id"` // id варианта продука
	StartDate time.Time        `json:"start_date" db:"start_date"` // дата начала цены
	EndDate   sqlnull.NullTime `json:"end_date" db:"end_date"`     // дата конца цены
	Price     float64          `json:"price" db:"price"`           // цена продукта
}

// ProductInfo структура информации о продукте о котором нужно получить информацию
type ProductInfo struct {
	ProductID   int       `db:"product_id"`       // id продукта
	Name        string    `db:"name"`             // название продукта
	Descr       string    `db:"description"`      // описание продукта
	VariantList []Variant `db:"product_variants"` // список вариантов продукта
}

// Sale структура продажи
type Sale struct {
	SaleID      int                `db:"sales_id"`                     // id продажи
	ProductName sqlnull.NullString `db:"name"`                         // id продукта
	VariantID   int                `json:"variant_id" db:"variant_id"` // id варианта продукта
	StorageID   int                `json:"storage_id" db:"storage_id"` // id склада из которого произошла продажа продукта
	SoldAt      time.Time          `db:"sold_at"`                      // дата продажи
	Quantity    int                `json:"quantity" db:"quantity"`     // кол-во проданного продукта
	TotalPrice  float64            `db:"total_price"`                  // общая стоимость с учетом кол-ва продукта
}

// IsNullFields проверка полей нва нулевые значения
func (s Sale) IsNullFields() error {
	if s.VariantID == 0 || s.StorageID == 0 || s.Quantity == 0 {
		return errors.New("variant_id, storage_id или quantity являются пустыми полями")
	}
	return nil
}

