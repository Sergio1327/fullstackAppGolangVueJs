package uimport

import (
	"product_storage/internal/usecase"
	"product_storage/internal/usecase/product"
)

type Usecase struct {
	Logger         *usecase.LoggerUsecase
	ProdcutUsecase product.ProductUseCase
}
