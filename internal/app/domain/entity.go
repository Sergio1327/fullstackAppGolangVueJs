package domain

import (
	"github.com/shopspring/decimal"
	"time"
)

type Product struct {
	ProductID  int       `json:"product_id"`
	Name       string    `json:"name"`
	Descr      string    `json:"description"`
	Addet_at   time.Time `json:"added_at"`
	Removed_at time.Time `json:"removed_at"`
	Tags       string  `json:"tags"`
	Variants   []Variant `json:"variants"`
}

type Variant struct {
	ProductId    int       `json:"product_id" db:"product_id"`
	VariantId    int       `json:"variant_id" db:"variant_id"`
	Weight       int       `json:"weight" db:"weight"`
	Unit         string    `json:"unit" db:"unit"`
	Added_at     time.Time `json:"added_at" db:"added_at"`
	CurrentPrice decimal.Decimal	`db:"price"`
	InStorages   []int
}

type ProductPrice struct {
	PriceId int 
	VariantId int             `json:"variant_id"`
	StartDate time.Time       `json:"start_date"`
	EndDate   time.Time       `json:"end_date"`
	Price     decimal.Decimal `json:"price"`
}

type AddProductInStock struct {
	VariantId  int       `json:"variant_id"`
	StorageId  int       `json:"storage_id"`
	Added_at   time.Time `json:"added_at"`
	Removed_at time.Time `json:"removed_at"`
	Quantity   int       `json:"quantity"`
}

type ProductInfo struct {
	ProductId int    `db:"product_id"`
	Name      string `db:"name"`
	Descr     string `db:"description"`
	Variants  []Variant
}

type Stock struct {
	StorageName string
	ProductID   int
	ProductName string
	VariantID   int
	VariantUnit string
	Weight      int
	Quantity    int
}

type Sale struct {
	SaleId      int
	VariantId   int `json:"variant_id"`
	StorageId   int `json:"storage_id"`
	SoldAt      time.Time
	Quantity    int `json:"quantity"`
	TotalPrice  decimal.Decimal
	ProductName string
}
type SaleQuery struct {
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Limit       int       `json:"limit"`
	StorageId   int       `json:"storage_id"`
	ProductName string    `json:"product_name"`
}
