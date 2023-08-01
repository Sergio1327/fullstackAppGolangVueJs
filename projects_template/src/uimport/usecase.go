package uimport

import (
	"projects_template/internal/usecase"
	"projects_template/internal/usecase/product"
)

type Usecase struct {
	Logger   *usecase.LoggerUsecase
	ProdcutUsecase product.ProductUseCase
}
