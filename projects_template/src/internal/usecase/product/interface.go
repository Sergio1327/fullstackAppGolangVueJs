package product

import (
	"product_storage/internal/entity/product"
	"product_storage/internal/entity/stock"

	"github.com/jmoiron/sqlx"
)

type ProductUseCase interface {
	AddProduct(tx *sqlx.Tx, product product.Product) (productID int, err error)
	AddProductPrice(tx *sqlx.Tx, pr product.ProductPrice) (int, error)
	AddProductInStock(tx *sqlx.Tx, p stock.AddProductInStock) (int, error)
	FindProductInfoById(tx *sqlx.Tx, productID int) (product.ProductInfo, error)
	FindProductList(tx *sqlx.Tx, tag string, limit int) ([]product.ProductInfo, error)
	FindProductsInStock(tx *sqlx.Tx, productID int) ([]stock.Stock, error)
	Buy(tx *sqlx.Tx, p product.Sale) (int, error)
	FindSaleList(tx *sqlx.Tx, sq product.SaleQuery) ([]product.Sale, error)
}
