package product

import (
	"errors"
	"product_storage/tools/sqlnull"
	"time"

	"github.com/sirupsen/logrus"
)

// Product cтруктура продукта для записи в базу
type ProductParams struct {
	ProductID   int              `json:"product_id"`  // id продукта
	Name        string           `json:"name"`        // название продукта
	Descr       string           `json:"description"` // описание продукта
	AddetAt     time.Time        `json:"added_at"`    // дата добавления продукта
	RemovedAt   sqlnull.NullTime `json:"removed_at"`  // дата удаления продукта
	Tags        string           `json:"tags"`        // теги продукта
	VariantList []Variant        `json:"variants"`    // cписок вариантов продукта
}

func (p ProductParams) Log() logrus.Fields {
	return logrus.Fields{"product_ID": p.ProductID}
}

// ProductPrice структура для вставки цены продукта
type ProductPriceParams struct {
	PriceID   int              `json:"price_id" db:"price_id"`     // id цены продукта
	VariantID int              `json:"variant_id" db:"variant_id"` // id варианта продука
	StartDate time.Time        `json:"start_date" db:"start_date"` // дата начала цены
	EndDate   sqlnull.NullTime `json:"end_date" db:"end_date"`     // дата конца цены
	Price     float64          `json:"price" db:"price"`           // цена продукта
}

func (p ProductPriceParams) Log() logrus.Fields {
	return logrus.Fields{"priceID": p.PriceID}
}

func (p ProductPriceParams) IsNullFields() error {
	if p.VariantID == 0 || p.Price == 0 || p.StartDate == (time.Time{}) {
		return errors.New("id варианта или цена или дата начала не могут быть пустыми")
	}
	return nil
}

// Sale структура продажи
type SaleParams struct {
	SaleID      int                `db:"sales_id"`                     // id продажи
	ProductName sqlnull.NullString `db:"name"`                         // id продукта
	VariantID   int                `json:"variant_id" db:"variant_id"` // id варианта продукта
	StorageID   int                `json:"storage_id" db:"storage_id"` // id склада из которого произошла продажа продукта
	SoldAt      time.Time          `db:"sold_at"`                      // дата продажи
	Quantity    int                `json:"quantity" db:"quantity"`     // кол-во проданного продукта
	TotalPrice  float64            `db:"total_price"`                  // общая стоимость с учетом кол-ва продукта
}

// IsNullFields проверка полей нва нулевые значения
func (s SaleParams) IsNullFields() error {
	if s.VariantID == 0 || s.StorageID == 0 || s.Quantity == 0 {
		return errors.New("variant_id, storage_id или quantity являются пустыми полями")
	}
	return nil
}
func (s SaleParams) Log() logrus.Fields {
	return logrus.Fields{"sale_ID": s.SaleID}
}

// SaleQuery фильтры продаж по которым нужно вывести информацию
type SaleQueryParam struct {
	StartDate   time.Time          `json:"start_date" db:"start_date"`     // дата начала продаж(обязательные поля)
	EndDate     time.Time          `json:"end_date"  db:"end_date"`        // дата конца прдаж (обязательные поля)
	Limit       sqlnull.NullInt64  `json:"limit" db:"limit"`               // лимит вывода продаж
	StorageID   sqlnull.NullInt64  `json:"storage_id" db:"storage_id"`     // id склада
	ProductName sqlnull.NullString `json:"product_name" db:"product_name"` // название продукта
}

func (s SaleQueryParam) Log() logrus.Fields {
	return logrus.Fields{"sale_query_params": s}
}

// SaleQuery фильтер продаж только дата продажи по которым нужно вывести информацию
type SaleQueryOnlyBySoldDateParam struct {
	StartDate time.Time         `json:"start_date" db:"sold_at"` // дата начала продаж
	EndDate   time.Time         `json:"end_date" db:"sold_at"`   // дата конца продаж
	Limit     sqlnull.NullInt64 `json:"limit" db:"limit"`        // лимит вывода
}

func (s SaleQueryOnlyBySoldDateParam) Log() logrus.Fields {
	return logrus.Fields{"sale_query_params": s}
}
