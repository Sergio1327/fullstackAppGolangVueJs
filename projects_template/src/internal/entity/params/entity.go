package params

import (
	"product_storage/tools/sqlnull"
	"time"
)

// SaleQuery фильтры продаж по которым нужно вывести информацию
type SaleQuery struct {
	StartDate   time.Time          `json:"start_date" db:"start_date"`     // дата начала продаж(обязательные поля)
	EndDate     time.Time          `json:"end_date"  db:"end_date"`        // дата конца прдаж (обязательные поля)
	Limit       sqlnull.NullInt64  `json:"limit" db:"limit"`               // лимит вывода продаж
	StorageID   sqlnull.NullInt64  `json:"storage_id" db:"storage_id"`     // id склада
	ProductName sqlnull.NullString `json:"product_name" db:"product_name"` // название продукта
}

// SaleQuery фильтер продаж только дата продажи по которым нужно вывести информацию
type SaleQueryOnlyBySoldDate struct {
	StartDate time.Time         `json:"start_date" db:"sold_at"` // дата начала продаж
	EndDate   time.Time         `json:"end_date" db:"sold_at"`   // дата конца продаж
	Limit     sqlnull.NullInt64 `json:"limit" db:"limit"`        // лимит вывода
}
