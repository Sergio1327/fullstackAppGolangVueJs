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

	AddProduct(tx *sqlx.Tx, p domain.Product) (int, error)
	AddProductVariants(tx *sqlx.Tx, id int, v domain.Variant) error

	CheckExists(p domain.ProductPrice) (int, error)
	UpdateProductPrice(tx *sqlx.Tx, p domain.ProductPrice, id int) error
	AddProductPriceWithEndDate(tx *sqlx.Tx, p domain.ProductPrice) error
	AddProductPrice(tx *sqlx.Tx, p domain.ProductPrice) error

	CheckProductsInStock(p domain.AddProductInStock) (bool, error)
	UpdateProductsInstock(tx *sqlx.Tx, p domain.AddProductInStock) error
	AddProductInStock(tx *sqlx.Tx, p domain.AddProductInStock) error

	LoadProductInfo(id int) (domain.ProductInfo, error)
	FindProductVariants(id int) ([]domain.Variant, error)
	FindCurrentPrice(variantId int) (decimal.Decimal, error)
	InStorages(id int) ([]int, error)

	FindProductsByTag(tag string, limit int) ([]domain.ProductInfo, error)
	LoadProducts(limit int) ([]domain.ProductInfo, error)

	LoadStocks() ([]domain.Stock, error)
	FindStocksByProductId(id int) ([]domain.Stock, error)
	FindStocksVariants(storageId int) ([]domain.AddProductInStock, error)

	Buy(tx *sqlx.Tx, s domain.Sale) error
	FindPrice(id int) (decimal.Decimal, error)

	LoadSales(sq domain.SaleQuery) ([]domain.Sale, error)
	FindSalesByFilters(sq domain.SaleQuery) ([]domain.Sale, error)
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

// AddProduct вставка названия,описания,времени добавления и тегов в базу
func (r *PostgresProductRepository) AddProduct(tx *sqlx.Tx, p domain.Product) (int, error) {
	var productId int
	err := tx.QueryRow("insert into products(name,description,added_at,tags) values ($1,$2,$3,$4) returning product_id",
		p.Name, p.Descr, p.Addet_at, p.Tags).Scan(&productId)
	if err != nil {
		return 0, err
	}
	return productId, nil
}

// AddProductVariants  добавление вариантов продукта в продукт по его id
func (r *PostgresProductRepository) AddProductVariants(tx *sqlx.Tx, id int, v domain.Variant) error {
	_, err := tx.Exec("insert into product_variants(product_id,weight,unit)values($1,$2,$3)", id, v.Weight, v.Unit)
	return err
}

// CheckExists проверка наличия цен варианта продукта в указаный диапазон времени
func (r *PostgresProductRepository) CheckExists(p domain.ProductPrice) (int, error) {
	var isExists int
	err := r.db.Get(&isExists,
		"select price_id from product_prices where variant_id=$1 and start_date=$2 and(end_date=$3 or end_date is null)",
		p.VariantId, p.StartDate, p.EndDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return isExists, nil
}

// UpdateProductPrice  обновление цены варианта продукта
func (r *PostgresProductRepository) UpdateProductPrice(tx *sqlx.Tx, p domain.ProductPrice, id int) error {
	_, err := tx.Exec("update product_prices set end_date=$1 where price_id=$2",
		p.EndDate, id)
	return err
}

// AddProductPriceWithEndDate  добавление цены варианта продукта в опеределенный диапазон времени
func (r *PostgresProductRepository) AddProductPriceWithEndDate(tx *sqlx.Tx, p domain.ProductPrice) error {
	_, err := tx.Exec("insert into product_prices(variant_id,price,start_date,end_date) values($1,$2,$3,$4)",
		p.VariantId, p.Price, p.StartDate, p.EndDate)
	return err
}

// AddProductPrice вставка цены варианта продукта в базу
func (r *PostgresProductRepository) AddProductPrice(tx *sqlx.Tx, p domain.ProductPrice) error {
	_, err := tx.Exec("insert into product_prices(variant_id,price,start_date) values($1,$2,$3)",
		p.VariantId, p.Price, p.StartDate)
	return err
}

// CheckProductsInStock  проверка есть ли на скалде продукт
func (r *PostgresProductRepository) CheckProductsInStock(p domain.AddProductInStock) (bool, error) {
	var isExists bool
	err := r.db.Get(&isExists, "select exists(select 1 from products_in_storage where variant_id=$1 and storage_id=$2)",
		p.VariantId, p.StorageId)

	return isExists, err
}

// UpdateProductsInstock обновление колличества продукта
func (r *PostgresProductRepository) UpdateProductsInstock(tx *sqlx.Tx, p domain.AddProductInStock) error {
	_, err := tx.Exec("update products_in_storage set quantity=$1 where variant_id=$2 and storage_id=$3",
		p.Quantity, p.VariantId, p.StorageId)
	return err
}

// AddProductInStock добавление продукта на склад
func (r *PostgresProductRepository) AddProductInStock(tx *sqlx.Tx, p domain.AddProductInStock) error {
	_, err := tx.Exec("insert into products_in_storage(variant_id,storage_id,added_at,quantity) values ($1,$2,$3,$4)",
		p.VariantId, p.StorageId, p.Added_at, p.Quantity)
	return err
}

// LoadProductInfo получение информации о продукте
func (r *PostgresProductRepository) LoadProductInfo(id int) (domain.ProductInfo, error) {
	var p domain.ProductInfo
	err := r.db.Get(&p, "select name,description from products where product_id=$1", id)
	return p, err
}

// FindProductVariants  получение вариантов продукта по его id
func (r *PostgresProductRepository) FindProductVariants(id int) ([]domain.Variant, error) {
	var v []domain.Variant
	err := r.db.Select(&v, "select variant_id,weight,unit,added_at from product_variants where product_id=$1", id)
	return v, err
}

// FindCurrentPrice  получение актуальной цены
func (r *PostgresProductRepository) FindCurrentPrice(variantId int) (decimal.Decimal, error) {
	var price decimal.Decimal
	err := r.db.Get(&price, "select price from product_prices where variant_id=$1 and start_date<$2 and (end_date is null or end_date>$2)",
		variantId, time.Now())

	return price, err
}

// InStorages  нахождение id складов в которых находится продукт
func (r *PostgresProductRepository) InStorages(id int) ([]int, error) {
	var inStorages []int
	err := r.db.Select(&inStorages, "SELECT storage_id FROM products_in_storage WHERE variant_id = $1", id)
	return inStorages, err
}

// FindProductsByTag  поиск информации о продукте по его тегу
func (r *PostgresProductRepository) FindProductsByTag(tag string, limit int) ([]domain.ProductInfo, error) {
	var products []domain.ProductInfo
	err := r.db.Select(&products, "select product_id,name,description from products where $1 = any (string_to_array(tags,',')) limit $2",
		tag, limit)
	return products, err
}

// LoadProducts  получение списка продуктов с лимитом
func (r *PostgresProductRepository) LoadProducts(limit int) ([]domain.ProductInfo, error) {
	var products []domain.ProductInfo
	err := r.db.Select(&products, "select product_id,name,description from products limit $1", limit)

	return products, err
}

// LoadStocks  Пплучение информации о складах
func (r *PostgresProductRepository) LoadStocks() ([]domain.Stock, error) {
	var stocks []domain.Stock
	err := r.db.Select(&stocks, "select  storage_id,name from storages")
	return stocks, err
}

// FindStocksByProductId получение информации о складах где есть определенный продукт
func (r *PostgresProductRepository) FindStocksByProductId(id int) ([]domain.Stock, error) {
	var stocks []domain.Stock
	err := r.db.Select(&stocks, `
	select s.storage_id ,s.name 
	from storages s
	join products_in_storage pis ON s.storage_id = pis.storage_id
	join product_variants pv ON pis.variant_id = pv.variant_id
	join products p ON pv.product_id = p.product_id
	where p.product_id=$1`, id)

	return stocks, err
}

// LoadStocksVariants  получение вариантов продукта на складе
func (r *PostgresProductRepository) FindStocksVariants(storageId int) ([]domain.AddProductInStock, error) {
	var v []domain.AddProductInStock
	err := r.db.Select(&v, "select variant_id,storage_id,added_at,quantity from products_in_storage where storage_id=$1 ", storageId)
	return v, err
}

// FindPrice  получение цены
func (r *PostgresProductRepository) FindPrice(id int) (decimal.Decimal, error) {
	var price decimal.Decimal
	err := r.db.Get(&price, "select price from product_prices where variant_id = $1", id)
	return price, err
}

// Buy запись о покупке в базу
func (r *PostgresProductRepository) Buy(tx *sqlx.Tx, s domain.Sale) error {
	_, err := tx.Exec("insert into sales(variant_id,storage_id,sold_at,quantity,total_price) values($1,$2,$3,$4,$5)",
		s.VariantId, s.StorageId, s.SoldAt, s.Quantity, s.TotalPrice)
	return err
}

// LoadSales получение списка всех продаж
func (r *PostgresProductRepository) LoadSales(sq domain.SaleQuery) ([]domain.Sale, error) {
	var sales []domain.Sale
	query := `SELECT s.sales_id, s.variant_id, s.storage_id, s.sold_at, s.quantity, s.total_price, p.name 
	FROM sales  s
	JOIN product_variants  pv ON (pv.variant_id = s.variant_id)
	JOIN products  p ON (p.product_id = pv.product_id)
	WHERE s.sold_at >= $1 AND s.sold_at <= $2
	LIMIT $3`

	err := r.db.Select(&sales, query, sq.StartDate, sq.EndDate, sq.Limit)
	return sales, err
}

// FindSalesByFilters  получение списка продаж по фильтрам
func (r *PostgresProductRepository) FindSalesByFilters(sq domain.SaleQuery) ([]domain.Sale, error) {
	var sales []domain.Sale

	query := `SELECT s.sales_id, s.variant_id, s.storage_id, s.sold_at, s.quantity, s.total_price, p.name 
	FROM sales s
	JOIN product_variants pv ON (pv.variant_id = s.variant_id)
	JOIN products p ON (p.product_id = pv.product_id)
	WHERE s.sold_at > :start_date AND s.sold_at < :end_date
	AND ( cast(:product_name as varchar) IS NULL OR p.name = :product_name )
	AND ( cast(:storage_id as integer) IS NULL OR s.storage_id = :storage_id ) 
	LIMIT :limit`

	params := map[string]interface{}{
		"start_date":   sq.StartDate,
		"end_date":     sq.EndDate,
		"product_name": sq.ProductName,
		"storage_id":   sq.StorageId,
		"limit":        sq.Limit,
	}

	stmt, err := r.db.PrepareNamed(query)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	defer stmt.Close()

	err = stmt.Select(&sales, params)
	return sales, err
}
