package usecase

import (
	"product_storage/internal/entity/product"
	"product_storage/internal/entity/stock"
	"product_storage/internal/transaction"
)

type Product interface {
	AddProduct(ts transaction.Session, product product.ProductParams) (productID int, err error)
	AddProductPrice(ts transaction.Session, pr product.ProductPriceParams) (int, error)
	AddProductInStock(ts transaction.Session, p stock.ProductInStockParams) (int, error)
	FindProductInfoById(ts transaction.Session, productID int) (product.ProductInfo, error)
	FindProductList(ts transaction.Session, tag string, limit int) ([]product.ProductInfo, error)
	FindProductsInStock(ts transaction.Session, productID int) ([]stock.Stock, error)
	SaveSale(ts transaction.Session, p product.SaleParams) (int, error)
	FindSaleList(ts transaction.Session, sq product.SaleQueryParam) ([]product.Sale, error)
	LoadStockList(ts transaction.Session) ([]stock.Stock, error)
}
