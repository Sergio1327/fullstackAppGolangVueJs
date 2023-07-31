package postgresql

import (
	"projects_template/internal/entity/product"
	"projects_template/internal/transaction"
)

type PostgresProductRepository struct {
}

func NewPostgresProductRepository() PostgresProductRepository {
	return PostgresProductRepository{}
}

func (r PostgresProductRepository) AddProduct(ts transaction.Session, product product.Product) (productID int, err error) {
	tx := SqlxTx(ts)

	query := `insert into products
	(name, description, added_at, tags)
	values ($1, $2, $3, $4) 
	returning product_id`

	err = tx.QueryRow(query, product.Name, product.Descr, product.AddetAt, product.Tags).Scan(&productID)

	return productID, err
}

func (r PostgresProductRepository) AddProductVariantList(ts transaction.Session, productID int, variant product.Variant) error {
	tx := SqlxTx(ts)
	query := `
	insert into product_variants 
	(product_id, weight, unit) 
	values ($1, $2, $3)`

	_, err 	:= tx.Exec(query, productID, variant.Weight, variant.Unit)
	return err
}
