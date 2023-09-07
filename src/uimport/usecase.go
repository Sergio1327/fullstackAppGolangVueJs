package uimport

import (
	"product_storage/internal/usecase"
)

type Usecase struct {
	Logger  *usecase.Logger
	Product *usecase.ProductUseCase
}
