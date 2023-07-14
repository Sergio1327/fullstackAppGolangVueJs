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

	CheckProductsInStock(p *domain.AddProductInStock) (bool, error)
	UpdateProductsInstock(p *domain.AddProductInStock) error
	AddProductInStock(p *domain.AddProductInStock) error

	GetProductInfo(id int, p *domain.ProductInfo) error
	GetProductVariants(id int, v []domain.Variant, p *domain.ProductInfo) error
	GetCurrentPrice(v *domain.Variant) error
	InStorages(v *domain.Variant) error

	GetProductsByTag(tag string, limit int) ([]domain.ProductInfo, error)
	GetProducts(limit int) ([]domain.ProductInfo, error)

	GetStocks() ([]domain.Stock, error)
	GetStocksByProductId(id int) ([]domain.Stock, error)
	GetStocksVariants(storageId int) ([]domain.AddProductInStock, error)

	Buy(s *domain.Sale) error
	GetPrice(id int) (decimal.Decimal, error)

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

func (r *PostgresProductRepository) CheckProductsInStock(p *domain.AddProductInStock) (bool, error) {
	var isExists bool
	err := r.db.Get(&isExists, "select exists(select 1 from products_in_storage where variant_id=$1 and storage_id=$2)",
		p.VariantId, p.StorageId)
	if err != nil {
		return false, err
	}
	return isExists, nil
}
func (r *PostgresProductRepository) UpdateProductsInstock(p *domain.AddProductInStock) error {
	_, err := r.db.Exec("update products_in_storage set quantity=$1 where variant_id=$2 and storage_id=$3",
		p.Quantity, p.VariantId, p.StorageId)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresProductRepository) AddProductInStock(p *domain.AddProductInStock) error {
	_, err := r.db.Exec("insert into products_in_storage(variant_id,storage_id,added_at,quantity) values ($1,$2,$3,$4)",
		p.VariantId, p.StorageId, p.Added_at, p.Quantity)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresProductRepository) GetProductInfo(id int, p *domain.ProductInfo) error {
	err := r.db.Get(p, "select name,description from products where product_id=$1", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresProductRepository) GetProductVariants(id int, v []domain.Variant, p *domain.ProductInfo) error {
	err := r.db.Select(&v, "select variant_id,weight,unit from product_variants where product_id=$1", id)
	if err != nil {
		return err
	}
	p.Variants = v
	return nil
}

func (r *PostgresProductRepository) GetCurrentPrice(v *domain.Variant) error {
	err := r.db.Get(&v.CurrentPrice, "select price from product_prices where variant_id=$1 and start_date<$2 and (end_date is null or end_date>$2)",
		v.VariantId, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresProductRepository) InStorages(v *domain.Variant) error {
	var inStorages []int
	err := r.db.Select(&inStorages, "SELECT storage_id FROM products_in_storage WHERE variant_id = $1", v.VariantId)
	if err != nil {
		return err
	}
	v.InStorages = inStorages
	return nil
}

func (r *PostgresProductRepository) GetProductsByTag(tag string, limit int) ([]domain.ProductInfo, error) {
	var products []domain.ProductInfo
	err := r.db.Select(&products, "select product_id,name,description from products where $1 = any (string_to_array(tags,',')) limit $2", tag, limit)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *PostgresProductRepository) GetProducts(limit int) ([]domain.ProductInfo, error) {
	var products []domain.ProductInfo
	err := r.db.Select(&products, "select product_id,name,description from products limit $1", limit)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *PostgresProductRepository) GetStocks() ([]domain.Stock, error) {
	var stocks []domain.Stock
	err := r.db.Select(&stocks, "select  storage_id,name from storages")
	if err != nil {
		return nil, err
	}
	return stocks, nil
}
func (r *PostgresProductRepository) GetStocksByProductId(id int) ([]domain.Stock, error) {
	var stocks []domain.Stock
	err := r.db.Select(&stocks, `
	select s.storage_id ,s.name 
	from storages s
	inner join products_in_storage pis ON s.storage_id = pis.storage_id
	inner join product_variants pv ON pis.variant_id = pv.variant_id
	inner join products p ON pv.product_id = p.product_id
	where p.product_id=$1`, id)
	if err != nil {
		return nil, err
	}
	return stocks, nil
}
func (r *PostgresProductRepository) GetStocksVariants(storageId int) ([]domain.AddProductInStock, error) {
	var v []domain.AddProductInStock
	err := r.db.Select(&v, "select variant_id,storage_id,added_at,quantity from products_in_storage where storage_id=$1 ", storageId)
	if err != nil {
		return nil, err
	}
	return v, err
}

func (r *PostgresProductRepository) GetPrice(id int) (decimal.Decimal, error) {
	var price decimal.Decimal
	err := r.db.Get(&price, "select price from product_prices where variant_id = $1", id)
	if err != nil {
		return decimal.Decimal{}, nil
	}
	return price, nil
}

func (r *PostgresProductRepository) Buy(s *domain.Sale) error {
	_, err := r.db.Exec("insert into sales(variant_id,storage_id,sold_at,quantity,total_price) values($1,$2,$3,$4,$5)",
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

	err := r.db.Select(&sales, query, sq.StartDate, sq.EndDate, sq.Limit)
	if err != nil {
		return nil, err
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
	err := r.db.Select(&sales, query, args...)
	if err!=nil{
		return nil,err
	}
	return sales, nil
}
