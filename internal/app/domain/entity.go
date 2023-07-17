package domain

import (
	"go-back/internal/tools/sqlnull"
	"time"

	"github.com/shopspring/decimal"
)

// Структрура продукта
type Product struct {
	ProductID  int              `json:"product_id"`
	Name       string           `json:"name"`
	Descr      string           `json:"description"`
	Addet_at   time.Time        `json:"added_at"`
	Removed_at sqlnull.NullTime `json:"removed_at"`
	Tags       string           `json:"tags"`
	Variants   []Variant        `json:"variants"`
}

// Вариант Продукта
type Variant struct {
	ProductId    int             `json:"product_id" db:"product_id"`
	VariantId    int             `json:"variant_id" db:"variant_id"`
	Weight       int             `json:"weight" db:"weight"`
	Unit         string          `json:"unit" db:"unit"`
	Added_at     time.Time       `json:"added_at" db:"added_at"`
	CurrentPrice decimal.Decimal `db:"price"`
	InStorages   []int
}

// Структура для вставки цены продукта
type ProductPrice struct {
	PriceId   int
	VariantId int              `json:"variant_id"`
	StartDate time.Time        `json:"start_date"`
	EndDate   sqlnull.NullTime `json:"end_date"`
	Price     decimal.Decimal  `json:"price"`
}

// Структура для вставки продукта на склад
type AddProductInStock struct {
	VariantId int       `json:"variant_id" db:"variant_id"`
	StorageId int       `json:"storage_id" db:"storage_id"`
	Added_at  time.Time `json:"added_at" db:"added_at" `
	Quantity  int       `json:"quantity" db:"quantity"`
}

type ProductInfo struct {
	ProductId int    `db:"product_id"`
	Name      string `db:"name"`
	Descr     string `db:"description"`
	Variants  []Variant
}

// струкьуоа склада
type Stock struct {
	StorageID       int    `db:"storage_id"`
	StorageName     string `db:"name"`
	ProductVariants []AddProductInStock
}

// Структура продажи
type Sale struct {
	SaleId      int             `db:"sales_id"`
	ProductName string          `db:"name"`
	VariantId   int             `json:"variant_id" db:"variant_id"`
	StorageId   int             `json:"storage_id" db:"storage_id"`
	SoldAt      time.Time       `db:"sold_at"`
	Quantity    int             `json:"quantity" db:"quantity"`
	TotalPrice  decimal.Decimal `db:"total_price"`
}
type SaleQuery struct {
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Limit       int       `json:"limit"`
	StorageId   int       `json:"storage_id"`
	ProductName string    `json:"product_name"`
}
