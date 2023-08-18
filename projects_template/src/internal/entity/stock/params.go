package stock

import (
	"errors"
	"product_storage/tools/sqlnull"
	"time"

	"github.com/sirupsen/logrus"
)

// AddProductInStock структура для вставки продукта на склад
type ProductInStockParams struct {
	ProductInStorageID int
	VariantID          int       `json:"variant_id" db:"variant_id"` // id варианта продукта
	StorageID          int       `json:"storage_id" db:"storage_id"` // id склада куда будет помещен этот продукт
	AddedAt            time.Time `json:"added_at" db:"added_at" `    // дата добавления продукта на склад
	Quantity           int       `json:"quantity" db:"quantity"`     // кол-во продукта добавленного на склад
}

func (p ProductInStockParams) Log() logrus.Fields {
	return logrus.Fields{
		"product_in_stock_ID": p.ProductInStorageID,
		"variant_ID":          p.VariantID,
		"storage_ID":          p.StorageID,
		"added_at":            p.AddedAt,
		"quantity":            p.Quantity,
	}
}

// IsNullFields проверка полей на нулевые значения
func (p ProductInStockParams) IsNullFields() error {
	if p.StorageID == 0 || p.VariantID == 0 || p.AddedAt == (time.Time{}) || p.Quantity == 0 {
		return errors.New("поля: variant_id, storage_id, added_at, quantity не должны быть пустыми")
	}
	return nil
}

type StockParams struct {
	StorageID   int              `db:"storage_id" json:"storage_id"`
	StorageName string           `db:"name" json:"storage_name"`
	Added_at    sqlnull.NullTime `db:"added_at" json:"added_at"`
}

func (s StockParams) Log() logrus.Fields {
	return logrus.Fields{
		"Storage_ID":  s.StorageID,
		"StorageName": s.StorageName,
		"Added_at":    s.Added_at,
	}
}
