package stock

import (
	"errors"
	"time"
)

// AddProductInStock структура для вставки продукта на склад
type AddProductInStock struct {
	VariantID int       `json:"variant_id" db:"variant_id"` // id варианта продукта
	StorageID int       `json:"storage_id" db:"storage_id"` // id склада куда будет помещен этот продукт
	AddedAt   time.Time `json:"added_at" db:"added_at" `    // дата добавления продукта на склад
	Quantity  int       `json:"quantity" db:"quantity"`     // кол-во продукта добавленного на склад
}

func (a AddProductInStock) IsNullFields() error {
	if a.StorageID == 0 || a.VariantID == 0 || a.AddedAt == (time.Time{}) || a.Quantity == 0 {
		return errors.New("поля: variant_id, storage_id, added_at, quantity не должны быть пустыми")
	}
	return nil
}

// Stock структура склада
type Stock struct {
	StorageID          int                 `db:"storage_id"`          // id склада
	StorageName        string              `db:"name"`                // название склада
	ProductVariantList []AddProductInStock `db:"products_in_storage"` // список продуктов на данном складе
}
