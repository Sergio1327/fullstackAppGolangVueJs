package rimport

import "product_storage/internal/repository"

type Repository struct {
	Logger   repository.Logger
	Template repository.Template
	Product  repository.ProductRepository
}

type MockRepository struct {
	Logger   *repository.MockLogger
	Template *repository.MockTemplate
}
