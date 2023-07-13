package repository

import (
	"database/sql"
	"go-back/internal/app/domain"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type ProductRepository interface {
	TxBegin() (*sqlx.Tx, error)

	AddProduct(p *domain.Product) (int, error)
	AddProductVariants(id int, v *domain.Variant) error

	CheckExists(p *domain.ProductPrice) (int, error)
	UpdateProductPrice(p *domain.ProductPrice, id int) error
	AddProductPriceWithEndDate(p *domain.ProductPrice) error
	AddProductPrice(p *domain.ProductPrice) error
	AddProductInStock(p *domain.AddProductInStock) error
	GetProductInfoById(id int) (domain.ProductInfo, error)
	GetProductListByTag(tag string, limit int) ([]domain.ProductInfo, error)
	GetProductList(limit int) ([]domain.ProductInfo, error)
	GetProductsInStock() ([]domain.Stock, error)
	GetProductsInStockById(id int) ([]domain.Stock, error)
	Buy(s *domain.Sale) error
	GetSales(sq *domain.SaleQuery) ([]domain.Sale, error)
	GetSalesByFilters(sq *domain.SaleQuery) ([]domain.Sale, error)
}

type PostgresProductRepository struct {
	db *sqlx.DB
}

func (p *PostgresProductRepository) TxBegin() (*sqlx.Tx, error) {
	tx, err := p.db.Beginx()
	if err != nil {
		log.Fatal(err)
	}
	return tx, nil
}

func NewPostgresProductRepository(db *sqlx.DB) *PostgresProductRepository {
	return &PostgresProductRepository{
		db: db,
	}
}

func (r *PostgresProductRepository) AddProduct(p *domain.Product) (int, error) {
	var productId int
	err := r.db.QueryRow("insert into products(name,description,added_at,tags) values ($1,$2,$3,$4) returning product_id",
		p.Name, p.Descr, p.Addet_at, p.Tags).Scan(&productId)
	if err != nil {
		log.Println(11111)
		return 0, err
	}
	return productId, nil
}

func (r *PostgresProductRepository) AddProductVariants(id int, v *domain.Variant) error {
	_, err := r.db.Exec("insert into product_variants(product_id,weight,unit)values($1,$2,$3)", id, v.Weight, v.Unit)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresProductRepository) CheckExists(p *domain.ProductPrice) (int, error) {
	var isExists int
	err := r.db.Get(&isExists,
		"select price_id from product_prices where variant_id=$1 and start_date=$2 and(end_date=$3 or end_date is null)",
		p.VariantId, p.StartDate, p.EndDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		log.Println(err)
		return 0, err
	}
	log.Println(isExists)
	return isExists, nil

}

func (r *PostgresProductRepository) UpdateProductPrice(p *domain.ProductPrice, id int) error {
	_, err := r.db.Exec("update product_prices set end_date=$1 where price_id=$2",
		p.EndDate, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresProductRepository) AddProductPriceWithEndDate(p *domain.ProductPrice) error {
	_, err := r.db.Exec("insert into product_prices(variant_id,price,start_date,end_date) values($1,$2,$3,$4)",
		p.VariantId, p.Price, p.StartDate, p.EndDate)
	if err != nil {
		log.Println(err, 2)
		return err
	}
	return nil
}
func (r *PostgresProductRepository) AddProductPrice(p *domain.ProductPrice) error {
	_, err := r.db.Exec("insert into product_prices(variant_id,price,start_date) values($1,$2,$3)",
		p.VariantId, p.Price, p.StartDate)
	if err != nil {
		log.Println(err, 3)
		return err
	}
	return nil
}




func (r *PostgresProductRepository) AddProductInStock(p *domain.AddProductInStock) error {

	var isExistId bool

	err := r.db.QueryRow("select exists(select 1 from products_in_storage where variant_id=$1 and storage_id=$2 and added_at=$3)",
		p.VariantId, p.StorageId, p.Added_at).Scan(&isExistId)
	if err != nil {

		return err
	}

	if isExistId {
		_, err := r.db.Exec("update products_in_storage set quantity=$1 where variant_id=$2 and storage_id=$3 and added_at=$4",
			p.Quantity, p.VariantId, p.StorageId, p.Added_at)
		if err != nil {
			return err
		}
	} else {
		_, err := r.db.Exec("insert into products_in_storage(variant_id,storage_id,added_at,quantity) values($1,$2,$3,$4)",
			p.VariantId, p.StorageId, p.Added_at, p.Quantity)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresProductRepository) GetProductInfoById(id int) (domain.ProductInfo, error) {
	var productInfo domain.ProductInfo
	tx, err := r.db.Begin()
	if err != nil {
		return domain.ProductInfo{}, err
	}
	err = r.db.QueryRow("select name,description from products where product_id=$1", id).
		Scan(&productInfo.Name, &productInfo.Descr)
	if err != nil {
		tx.Rollback()
		return domain.ProductInfo{}, err
	}

	rows, err := r.db.Queryx("select variant_id,weight,unit from product_variants where product_id=$1", id)
	if err != nil {
		tx.Rollback()
		return domain.ProductInfo{}, err
	}
	defer rows.Close()
	var p_variants []domain.Variant
	for rows.Next() {
		var variant domain.Variant
		err := rows.Scan(&variant.VariantId, &variant.Weight, &variant.Unit)
		if err != nil {
			tx.Rollback()
			return domain.ProductInfo{}, err
		}

		err = r.db.QueryRow("select price from product_prices where variant_id=$1 and start_date<$2 and (end_date is null or end_date>$2)",
			variant.VariantId, time.Now()).Scan(&variant.CurrentPrice)
		if err != nil {
			return domain.ProductInfo{}, err
		}

		rows, err := r.db.Query("SELECT storage_id FROM products_in_storage WHERE variant_id = $1", variant.VariantId)
		if err != nil {
			return domain.ProductInfo{}, err
		}
		defer rows.Close()

		for rows.Next() {
			var storageID int
			err := rows.Scan(&storageID)
			if err != nil {
				return domain.ProductInfo{}, err
			}

			variant.InStorages = append(variant.InStorages, storageID)
		}

		p_variants = append(p_variants, variant)
	}
	productInfo.Variants = p_variants
	return productInfo, nil
}

func (r *PostgresProductRepository) GetProductListByTag(tag string, limit int) ([]domain.ProductInfo, error) {
	var products []domain.ProductInfo
	rows, err := r.db.Query(`select product_id,name,description from products where $1 = any (select jsonb_array_elements_text(tags)) limit $2`,
		tag, limit)
	if err != nil {
		log.Println(err, 1)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var product domain.ProductInfo
		err := rows.Scan(&product.ProductId, &product.Name, &product.Descr)
		if err != nil {
			log.Println(err, 2)
			return nil, err
		}
		vRows, err := r.db.Query("select variant_id,weight,unit from product_variants where product_id=$1", product.ProductId)
		if err != nil {
			log.Println(err, 3)
			return nil, err
		}
		defer vRows.Close()
		for vRows.Next() {
			var variant domain.Variant
			err := vRows.Scan(&variant.VariantId, &variant.Weight, &variant.Unit)
			if err != nil {
				log.Println(4)
				return nil, err
			}

			err = r.db.QueryRow("select price from product_prices where variant_id=$1 and start_date<$2 and(end_date>$2 or end_date is null)",
				variant.VariantId, time.Now()).Scan(&variant.CurrentPrice)
			if err != nil {
				log.Println(5)
				return nil, err
			}

			product.Variants = append(product.Variants, variant)

		}
		products = append(products, product)
	}
	return products, nil
}

func (r *PostgresProductRepository) GetProductList(limit int) ([]domain.ProductInfo, error) {
	var products []domain.ProductInfo
	rows, err := r.db.Queryx("select product_id,name,description from products limit $1", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var product domain.ProductInfo
		err := rows.Scan(&product.ProductId, &product.Name, &product.Descr)
		if err != nil {
			return nil, err
		}
		vRows, err := r.db.Query("select variant_id,weight,unit from product_variants where product_id=$1", product.ProductId)
		if err != nil {
			return nil, err
		}
		defer vRows.Close()
		for vRows.Next() {
			var variant domain.Variant
			err := vRows.Scan(&variant.VariantId, &variant.Weight, &variant.Unit)
			if err != nil {
				return nil, err
			}

			err = r.db.QueryRow("select price from product_prices where variant_id=$1 and start_date<$2 and(end_date>$2 or end_date is null)",
				variant.VariantId, time.Now()).Scan(&variant.CurrentPrice)
			if err != nil {
				return nil, err
			}

			product.Variants = append(product.Variants, variant)

		}
		products = append(products, product)
	}
	return products, nil
}

func (r *PostgresProductRepository) GetProductsInStock() ([]domain.Stock, error) {
	var stocks []domain.Stock
	rows, err := r.db.Query(`
	select s.name AS storage_name, p.product_id, p.name AS product_name, pv.variant_id, pv.unit AS variant_unit, pv.weight, pis.quantity
	from storages s
	inner join products_in_storage pis ON s.storage_id = pis.storage_id
	inner join product_variants pv ON pis.variant_id = pv.variant_id
	inner join products p ON pv.product_id = p.product_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stock domain.Stock
		err := rows.Scan(&stock.StorageName, &stock.ProductID, &stock.ProductName, &stock.VariantID, &stock.VariantUnit, &stock.Weight, &stock.Quantity)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}
	return stocks, nil
}

func (r *PostgresProductRepository) GetProductsInStockById(id int) ([]domain.Stock, error) {
	var stocks []domain.Stock
	rows, err := r.db.Query(`
	select s.name AS storage_name, p.product_id, p.name AS product_name, pv.variant_id, pv.unit AS variant_unit, pv.weight, pis.quantity
	from storages s
	inner join products_in_storage pis ON s.storage_id = pis.storage_id
	inner join product_variants pv ON pis.variant_id = pv.variant_id
	inner join products p ON pv.product_id = p.product_id
	where p.product_id=$1`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stock domain.Stock
		err := rows.Scan(&stock.StorageName, &stock.ProductID, &stock.ProductName, &stock.VariantID, &stock.VariantUnit, &stock.Weight, &stock.Quantity)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}
	return stocks, nil
}

func (r *PostgresProductRepository) Buy(s *domain.Sale) error {
	s.SoldAt = time.Now()
	var price decimal.Decimal
	err := r.db.QueryRow("select price from product_prices where variant_id=$1", s.VariantId).Scan(&price)
	if err != nil {
		return err
	}
	s.TotalPrice = price.Mul(decimal.NewFromInt(int64(s.Quantity)))
	_, err = r.db.Exec("insert into sales(variant_id,storage_id,sold_at,quantity,total_price) values($1,$2,$3,$4,$5)",
		s.VariantId, s.StorageId, s.SoldAt, s.Quantity, s.TotalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresProductRepository) GetSales(sq *domain.SaleQuery) ([]domain.Sale, error) {
	var sales []domain.Sale
	query := `select sales.sales_id, sales.variant_id, sales.storage_id, sales.sold_at, sales.quantity, sales.total_price,
	products.name 
	from sales
	inner join product_variants AS product_variants ON product_variants.variant_id = sales.variant_id
	inner join products AS products ON products.product_id = product_variants.product_id
	where sales.sold_at >= $1 AND sales.sold_at <= $2 limit $3`

	rows, err := r.db.Queryx(query, sq.StartDate, sq.EndDate, sq.Limit)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var sale domain.Sale
		err := rows.Scan(&sale.SaleId, &sale.VariantId, &sale.StorageId, &sale.SoldAt, &sale.Quantity, &sale.TotalPrice, &sale.ProductName)
		if err != nil {
			return nil, err
		}
		sales = append(sales, sale)
	}

	return sales, nil
}

func (r *PostgresProductRepository) GetSalesByFilters(sq *domain.SaleQuery) ([]domain.Sale, error) {
	var sales []domain.Sale
	query := `select sales.sales_id, sales.variant_id, sales.storage_id, sales.sold_at, sales.quantity, sales.total_price,
	products.name 
	from sales
	inner join product_variants AS product_variants ON product_variants.variant_id = sales.variant_id
	inner join products AS products ON products.product_id = product_variants.product_id
	where sales.sold_at >= $1 AND sales.sold_at <= $2`

	args := []interface{}{sq.StartDate, sq.EndDate}
	if sq.StorageId != 0 && sq.ProductName == "" {
		args = append(args, sq.StorageId)
		query += ` and sales.storage_id = $3 limit $4 `
	} else if sq.StorageId == 0 && sq.ProductName != "" {
		args = append(args, sq.ProductName)
		query += ` and products.name = $3 limit $4`
	} else {
		args = append(args, sq.StorageId, sq.ProductName)
		query += ` and storage_id = $3 and products.name = $4 limit $5`
	}

	args = append(args, sq.Limit)
	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var sale domain.Sale
		err := rows.Scan(&sale.SaleId, &sale.VariantId, &sale.StorageId, &sale.SoldAt, &sale.Quantity, &sale.TotalPrice, &sale.ProductName)
		if err != nil {
			return nil, err
		}
		sales = append(sales, sale)
	}

	return sales, nil
}
