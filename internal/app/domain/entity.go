package domain

import (
	"errors"
	"go-back/internal/tools/sqlnull"
	"time"
)

// Product cтруктура продукта для записи в базу
type Product struct {
	ProductID  int              `json:"product_id"`  // Id продукта
	Name       string           `json:"name"`        // Название продукта
	Descr      string           `json:"description"` // Описание продукта
	AddetAt   time.Time        `json:"added_at"`    //Дата добавления продукта
	RemovedAt sqlnull.NullTime `json:"removed_at"`  // дата удаления продукта
	Tags       string           `json:"tags"`        //теги продукта
	VariantList   []Variant        `json:"variants"`    //	Список вариантов продукта
}

// Variant структура варианта, продукта представляем с собой информацию о продукте который нужно внести в базу
type Variant struct {
	ProductId    int       `json:"product_id" db:"product_id"` //id продука
	VariantId    int       `json:"variant_id" db:"variant_id"` //id конкретного варианта продукта
	Weight       int       `json:"weight" db:"weight"`         // масса или вес продукта
	Unit         string    `json:"unit" db:"unit"`             //единица измерения
	AddedAt     time.Time `json:"added_at" db:"added_at"`     // дата добавления определенного варианта
	CurrentPrice float64   `json:"price" db:"price"`           //актуальная цена
	InStorages   []int     `json:"in_storages"`                //список id складов в которых есть этот вариант
}

// ProductPrice структура для вставки цены продукта
type ProductPrice struct {
	PriceId   int              //id цены продукта
	VariantId int              `json:"variant_id"` //id варианта продука
	StartDate time.Time        `json:"start_date"` // дата начала цены
	EndDate   sqlnull.NullTime `json:"end_date"`   //дата конца цены
	Price     float64          `json:"price"`      //цена продукта
}

// AddProductInStock  структура для вставки продукта на склад
type AddProductInStock struct {
	VariantId int       `json:"variant_id" db:"variant_id"` //id варианта продукта
	StorageId int       `json:"storage_id" db:"storage_id"` //id склада куда будет помещен этот продукт
	AddedAt  time.Time `json:"added_at" db:"added_at" `    //дата добавления продукта на склад
	Quantity  int       `json:"quantity" db:"quantity"`     //кол-во продукта добавленного на склад
}

func (a *AddProductInStock) IsNullFields() error {
	if a.StorageId == 0 || a.VariantId == 0 || a.AddedAt == (time.Time{}) || a.Quantity == 0 {
		return errors.New("поля: variant_id, storage_id, added_at, quantity не должны быть пустыми")
	}
	return nil
}

// ProductInfo структура информации о продукте о котором нужно получить информацию
type ProductInfo struct {
	ProductId int       `db:"product_id"`  //id продукта
	Name      string    `db:"name"`        //название продукта
	Descr     string    `db:"description"` //описание продукта
	VariantList  []Variant //список вариантов продукта
}

// Stock структура склада
type Stock struct {
	StorageID       int                 `db:"storage_id"` //id склада
	StorageName     string              `db:"name"`       //название склада
	ProductVariantList []AddProductInStock //список продуктов на данном складе
}

// Sale структура продажи
type Sale struct {
	SaleId      int                `db:"sales_id"`                     //id продажи
	ProductName sqlnull.NullString `db:"name"`                         //id продукта
	VariantId   int                `json:"variant_id" db:"variant_id"` //id варианта продукта
	StorageId   int                `json:"storage_id" db:"storage_id"` //id склада из которого произошла продажа продукта
	SoldAt      time.Time          `db:"sold_at"`                      //дата продажи
	Quantity    int                `json:"quantity" db:"quantity"`     //кол-во проданного продукта
	TotalPrice  float64            `db:"total_price"`                  //общая стоимость с учетом кол-ва продукта
}

func (s *Sale) IsNullFields() error {
	if s.VariantId == 0 || s.StorageId == 0 || s.Quantity == 0 {
		return errors.New("variant_id, storage_id или quantity являются пустыми полями")
	}
	return nil
}

// SaleQuery  фильтры продаж по которым нужно вывести информацию
type SaleQuery struct {
	StartDate   time.Time          `json:"start_date" db:"sold_at"`        //дата начала продаж(обязательные поля)
	EndDate     time.Time          `json:"end_date"  db:"sold_at"`         //дата конца прдаж (обязательные поля)
	Limit       sqlnull.NullInt64  `json:"limit" db:"limit"`               //лимит вывода продаж
	StorageId   sqlnull.NullInt64  `json:"storage_id" db:"storage_id"`     //id склада
	ProductName sqlnull.NullString `json:"product_name" db:"product_name"` //название продукта
}

type SaleQueryWithoutFilters struct {
	StartDate time.Time         `json:"start_date" db:"sold_at"`
	EndDate   time.Time         `json:"end_date" db:"sold_at"`
	Limit     sqlnull.NullInt64 `json:"limit" db:"limit"`
}
