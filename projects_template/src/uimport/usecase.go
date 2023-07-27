package uimport

import "projects_template/internal/usecase"

type Usecase struct {
	Logger   *usecase.LoggerUsecase
	Template *usecase.TemplateUsecase
}
