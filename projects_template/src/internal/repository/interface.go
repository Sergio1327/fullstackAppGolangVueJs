package repository

import (
	"product_storage/internal/entity/log"
	"product_storage/internal/entity/product"
	"product_storage/internal/entity/stock"
	"product_storage/internal/transaction"
)

type Logger interface {
	SaveLog(ts transaction.Session, row log.Row) error
	SaveLogWithReturnID(ts transaction.Session, row log.Row) (logID int, err error)
	SaveLogDetails(ts transaction.Session, logID int, details map[string]string) error
}

//go:generate mockgen -source=interface.go -destination=mocks/mock.go

type Product interface {
	AddProduct(ts transaction.Session, product product.ProductParams) (productID int, err error)
	AddProductVariantList(ts transaction.Session, productID int, variant product.Variant) error

	CheckExists(ts transaction.Session, p product.ProductPriceParams) (int, error)
	UpdateProductPrice(ts transaction.Session, p product.ProductPriceParams, id int) error
	AddProductPrice(ts transaction.Session, p product.ProductPriceParams) (int, error)

	CheckProductInStock(ts transaction.Session, p stock.ProductInStockParams) (bool, error)
	UpdateProductInstock(ts transaction.Session, p stock.ProductInStockParams) (int, error)
	AddProductInStock(ts transaction.Session, p stock.ProductInStockParams) (int, error)

	LoadProductInfo(ts transaction.Session, productID int) (product.ProductInfo, error)
	FindProductVariantList(ts transaction.Session, productID int) ([]product.Variant, error)
	FindCurrentPrice(ts transaction.Session, variantID int) (float64, error)
	InStorages(ts transaction.Session, variantID int) ([]int, error)

	FindProductListByTag(ts transaction.Session, tag string, limit int) ([]product.ProductInfo, error)
	LoadProductList(ts transaction.Session, limit int) ([]product.ProductInfo, error)

	LoadStockList(ts transaction.Session) ([]stock.Stock, error)
	FindStockListByProductId(ts transaction.Session, productID int) ([]stock.Stock, error)
	FindStocksVariantList(ts transaction.Session, storageID int) ([]stock.ProductInStockParams, error)

	SaveSale(ts transaction.Session, s product.SaleParams) (int, error)
	FindPrice(ts transaction.Session, variantID int) (float64, error)

	FindSaleListOnlyBySoldDate(ts transaction.Session, sq product.SaleQueryOnlyBySoldDateParam) ([]product.Sale, error)
	FindSaleListByFilters(ts transaction.Session, sq product.SaleQueryParam) ([]product.Sale, error)

}
