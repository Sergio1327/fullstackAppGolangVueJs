package product

import (
	"product_storage/tools/sqlnull"
	"time"
)

// Variant структура варианта, продукта представляем с собой информацию о продукте который нужно внести в базу
type Variant struct {
	ProductID    int          `json:"product_id" db:"product_id"` // id продука
	VariantID    int          `json:"variant_id" db:"variant_id"` // id конкретного варианта продукта
	Weight       int          `json:"weight" db:"weight"`         // масса или вес продукта
	Unit         string       `json:"unit" db:"unit"`             // единица измерения
	AddedAt      time.Time    `json:"added_at" db:"added_at"`     // дата добавления определенного варианта
	CurrentPrice float64      `json:"price" db:"price"`           // актуальная цена
	InStorages   []VarStorage `json:"in_storages"`                // список названий складов в которых есть этот вариант
}
type VarStorage struct {
	StorageID   int    `db:"storage_id"`
	StorageName string `db:"name"`
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
