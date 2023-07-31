package product

import (
	"projects_template/internal/entity/product"
	"projects_template/internal/transaction"
)

type ProductUseCase interface {
	AddProduct(ts transaction.Session, product product.Product) (productID int, err error)
}
