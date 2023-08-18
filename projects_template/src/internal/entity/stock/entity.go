package stock

import (
	
)

// Stock структура склада
type Stock struct {
	StorageID          int                    `db:"storage_id"`          // id склада
	StorageName        string                 `db:"name"`                // название склада
	ProductVariantList []ProductInStockParams `db:"products_in_storage"` // список продуктов на данном складе
}


