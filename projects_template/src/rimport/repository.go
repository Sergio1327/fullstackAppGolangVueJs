package rimport

import "projects_template/internal/repository"

type Repository struct {
	Logger   repository.Logger
	Template repository.Template
}

type MockRepository struct {
	Logger   *repository.MockLogger
	Template *repository.MockTemplate
}
