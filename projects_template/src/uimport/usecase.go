package uimport

import (
	"product_storage/internal/usecase"
)

type Usecase struct {
	Logger         *usecase.LoggerUsecase
	ProdcutUsecase usecase.Product
}
