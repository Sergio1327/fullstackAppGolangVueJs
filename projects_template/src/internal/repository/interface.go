package repository

import (
	"product_storage/internal/entity/log"
	"product_storage/internal/entity/product"
	"product_storage/internal/entity/stock"
	"product_storage/internal/transaction"
	"product_storage/tools/sqlnull"

	"github.com/jmoiron/sqlx"	
)

type Logger interface {
	SaveLog(ts transaction.Session, row log.Row, contractID, seID sqlnull.NullInt64, operLogin sqlnull.NullString) error
	SaveLogWithReturnID(ts transaction.Session, row log.Row, contractID, seID sqlnull.NullInt64, operLogin sqlnull.NullString) (logID int, err error)
	SaveLogDetails(ts transaction.Session, logID int, details map[string]string) error
}

//go:generate mockgen -source=interface.go -destination=mocks/mock.go

type ProductRepository interface {
	AddProduct(tx *sqlx.Tx, product product.Product) (productID int, err error)
	AddProductVariantList(tx *sqlx.Tx, productID int, variant product.Variant) error

	CheckExists(tx *sqlx.Tx, p product.ProductPrice) (int, error)
	UpdateProductPrice(tx *sqlx.Tx, p product.ProductPrice, id int) error
	AddProductPrice(tx *sqlx.Tx, p product.ProductPrice) (int, error)

	CheckProductInStock(tx *sqlx.Tx, p stock.AddProductInStock) (bool, error)
	UpdateProductInstock(tx *sqlx.Tx, p stock.AddProductInStock) (int, error)
	AddProductInStock(tx *sqlx.Tx, p stock.AddProductInStock) (int, error)

	LoadProductInfo(tx *sqlx.Tx, productID int) (product.ProductInfo, error)
	FindProductVariantList(tx *sqlx.Tx, productID int) ([]product.Variant, error)
	FindCurrentPrice(tx *sqlx.Tx, variantID int) (float64, error)
	InStorages(tx *sqlx.Tx, variantID int) ([]int, error)

	FindProductListByTag(tx *sqlx.Tx, tag string, limit int) ([]product.ProductInfo, error)
	LoadProductList(tx *sqlx.Tx, limit int) ([]product.ProductInfo, error)

	LoadStockList(tx *sqlx.Tx) ([]stock.Stock, error)
	FindStockListByProductId(tx *sqlx.Tx, productID int) ([]stock.Stock, error)
	FindStocksVariantList(tx *sqlx.Tx, storageID int) ([]stock.AddProductInStock, error)

	Buy(tx *sqlx.Tx, s product.Sale) (int, error)
	FindPrice(tx *sqlx.Tx, variantID int) (float64, error)

	FindSaleListOnlyBySoldDate(tx *sqlx.Tx, sq product.SaleQueryOnlyBySoldDate) ([]product.Sale, error)
	FindSaleListByFilters(tx *sqlx.Tx, sq product.SaleQuery) ([]product.Sale, error)
}
