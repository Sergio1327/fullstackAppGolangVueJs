package repository

import (
	"projects_template/internal/entity/log"
	"projects_template/internal/entity/product"
	"projects_template/internal/entity/template"
	"projects_template/internal/transaction"
	"projects_template/tools/sqlnull"
)

type Logger interface {
	SaveLog(ts transaction.Session, row log.Row, contractID, seID sqlnull.NullInt64, operLogin sqlnull.NullString) error
	SaveLogWithReturnID(ts transaction.Session, row log.Row, contractID, seID sqlnull.NullInt64, operLogin sqlnull.NullString) (logID int, err error)
	SaveLogDetails(ts transaction.Session, logID int, details map[string]string) error
}

type Template interface {
	FindTemplateObj(ts transaction.Session, id int) (template.TemplateObject, error)
	LoadAllTemplateObj(ts transaction.Session) ([]template.TemplateObject, error)
}

type ProductRepository interface {
	AddProduct(ts transaction.Session, product product.Product) (productID int, err error)
	AddProductVariantList(ts transaction.Session, productID int, variant product.Variant) error
}
