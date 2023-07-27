package oracle

import (
	"projects_template/internal/entity/template"
	"projects_template/internal/repository"
	"projects_template/internal/transaction"
	"projects_template/tools/gensql"
)

type templateRepository struct {
}

func NewTemplate() repository.Template {
	return &templateRepository{}
}

func (r *templateRepository) FindTemplateObj(ts transaction.Session, id int) (template.TemplateObject, error) {
	sqlQuery := `SELECT * from template_table where id = :1`

	return gensql.Get[template.TemplateObject](SqlxTx(ts), sqlQuery, id)
}

func (r *templateRepository) LoadAllTemplateObj(ts transaction.Session) ([]template.TemplateObject, error) {
	sqlQuery := `SELECT * from template_table`

	return gensql.Select[template.TemplateObject](SqlxTx(ts), sqlQuery)
}
