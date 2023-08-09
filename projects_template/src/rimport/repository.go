package rimport

import "product_storage/internal/repository"

type Repository struct {
	Logger  repository.Logger
	Product repository.Product
}

type MockRepository struct {
	Logger  *repository.MockLogger
	Product *repository.MockProduct
}
