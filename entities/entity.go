package entities

import (
	"time"

	"github.com/shopspring/decimal"
)

type Product struct {
	ProductId   int               `json:"product_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Added_at    time.Time         `json:"added_at"`
	Tags        []string          `json:"tags"`
	Variants    []ProductVariants `json:"variants"`
}

type ProductVariants struct {
	VariantId    int             `json:"variant_id"`
	Weight       decimal.Decimal `json:"weight"`
	Unit         string          `json:"unit"`
	CurrentPrice float64         `json:""`
}

type ProductPrice struct {
	VariantId int             `json:"variant_id"`
	Price     decimal.Decimal `json:"price"`
	StartDate time.Time       `json:"start_date"`
	EndDate   time.Time       `json:"endDate"`
}

type Product_in_stock struct {
	VariantId int `json:"variant_id"`
	Quantity  int `json:"quantity"`
	StockId   int `json:"storage_id"`
}

type ProductInfo struct {
	Product     Product `json:"prod ct"`
	StoragesIds []int
}

type Stock struct {
	StorageID    int    `json:"storage_id"`
	StorageName  string `json:"storage_name"`
	VariantId    int    `json:"variant_id"`
	ProductCount int    `json:"quantity"`
}

type Purchase struct {
	VariantId int `json:"variant_id"`
	Quantity  int `json:"quantity"`
}
type Sale struct {
	ID            int             `json:"sales_id"`
	ProductName   string          `json:"product_name"`
	VariantWeight float64         `json:"variant_weight"`
	VariantUnit   string          `json:"variant_unit"`
	Price         decimal.Decimal `json:"price"`
	Quantity      int             `json:"quantity"`
	TotalPrice    decimal.Decimal `json:"total_price"`
	SoldAt        time.Time       `json:"sold_at"`
}
