package repository

import (
	"product_storage/internal/entity/log"
	"product_storage/internal/entity/product"
	"product_storage/internal/entity/stock"
	"product_storage/internal/transaction"
	"product_storage/tools/sqlnull"
)

type Logger interface {
	SaveLog(ts transaction.Session, row log.Row, operLogin sqlnull.NullString) error
	SaveLogWithReturnID(ts transaction.Session, row log.Row, contractID, seID sqlnull.NullInt64, operLogin sqlnull.NullString) (logID int, err error)
	SaveLogDetails(ts transaction.Session, logID int, details map[string]string) error
}

//go:generate mockgen -source=interface.go -destination=mocks/mock.go

type Product interface {
	AddProduct(ts transaction.Session, product product.Product) (productID int, err error)
	AddProductVariantList(ts transaction.Session, productID int, variant product.Variant) error

	CheckExists(ts transaction.Session, p product.ProductPrice) (int, error)
	UpdateProductPrice(ts transaction.Session, p product.ProductPrice, id int) error
	AddProductPrice(ts transaction.Session, p product.ProductPrice) (int, error)

	CheckProductInStock(ts transaction.Session, p stock.AddProductInStock) (bool, error)
	UpdateProductInstock(ts transaction.Session, p stock.AddProductInStock) (int, error)
	AddProductInStock(ts transaction.Session, p stock.AddProductInStock) (int, error)

	LoadProductInfo(ts transaction.Session, productID int) (product.ProductInfo, error)
	FindProductVariantList(ts transaction.Session, productID int) ([]product.Variant, error)
	FindCurrentPrice(ts transaction.Session, variantID int) (float64, error)
	InStorages(ts transaction.Session, variantID int) ([]int, error)

	FindProductListByTag(ts transaction.Session, tag string, limit int) ([]product.ProductInfo, error)
	LoadProductList(ts transaction.Session, limit int) ([]product.ProductInfo, error)

	LoadStockList(ts transaction.Session) ([]stock.Stock, error)
	FindStockListByProductId(ts transaction.Session, productID int) ([]stock.Stock, error)
	FindStocksVariantList(ts transaction.Session, storageID int) ([]stock.AddProductInStock, error)

	Buy(ts transaction.Session, s product.Sale) (int, error)
	FindPrice(ts transaction.Session, variantID int) (float64, error)

	FindSaleListOnlyBySoldDate(ts transaction.Session, sq product.SaleQueryOnlyBySoldDate) ([]product.Sale, error)
	FindSaleListByFilters(ts transaction.Session, sq product.SaleQuery) ([]product.Sale, error)
}
