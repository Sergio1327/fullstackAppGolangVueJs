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

	AddProduct(tx sqlx.Tx, p domain.Product) (int, error)
	AddProductVariants(tx sqlx.Tx, id int, v domain.Variant) error

	CheckExists(p domain.ProductPrice) (int, error)
	UpdateProductPrice(tx sqlx.Tx, p domain.ProductPrice, id int) error
	AddProductPriceWithEndDate(tx sqlx.Tx, p domain.ProductPrice) error
	AddProductPrice(tx sqlx.Tx, p domain.ProductPrice) error

	CheckProductsInStock(p domain.AddProductInStock) (bool, error)
	UpdateProductsInstock(tx sqlx.Tx, p domain.AddProductInStock) error
	AddProductInStock(tx sqlx.Tx, p domain.AddProductInStock) error

	LoadProductInfo(id int) (domain.ProductInfo, error)
	LoadProductVariants(id int) ([]domain.Variant, error)
	LoadCurrentPrice(variantId int) (decimal.Decimal, error)
	InStorages(id int) ([]int, error)

	FindProductsByTag(tag string, limit int) ([]domain.ProductInfo, error)
	LoadProducts(limit int) ([]domain.ProductInfo, error)

	LoadStocks() ([]domain.Stock, error)
	FindStocksByProductId(id int) ([]domain.Stock, error)
	LoadStocksVariants(storageId int) ([]domain.AddProductInStock, error)

	Buy(tx sqlx.Tx, s domain.Sale) error
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

// вставка названия,описания,времени добавления и тегов в базу
func (r *PostgresProductRepository) AddProduct(tx sqlx.Tx, p domain.Product) (int, error) {
	var productId int
	err := tx.QueryRow("insert into products(name,description,added_at,tags) values ($1,$2,$3,$4) returning product_id",
		p.Name, p.Descr, p.Addet_at, p.Tags).Scan(&productId)
	if err != nil {
		return 0, err
	}
	return productId, nil
}

// Добавление вариантов продукта в продукт по его id
func (r *PostgresProductRepository) AddProductVariants(tx sqlx.Tx, id int, v domain.Variant) error {
	_, err := tx.Exec("insert into product_variants(product_id,weight,unit)values($1,$2,$3)", id, v.Weight, v.Unit)
	if err != nil {
		return err
	}
	return nil
}

// Проверка наличия цен варианта продукта в указаный диапазон времени
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

// Обновление цены варианта продукта
func (r *PostgresProductRepository) UpdateProductPrice(tx sqlx.Tx, p domain.ProductPrice, id int) error {
	_, err := tx.Exec("update product_prices set end_date=$1 where price_id=$2",
		p.EndDate, id)
	if err != nil {
		return err
	}
	return nil
}

// Добавление цены варианта продукта в опеределенный диапазон времени
func (r *PostgresProductRepository) AddProductPriceWithEndDate(tx sqlx.Tx, p domain.ProductPrice) error {
	_, err := tx.Exec("insert into product_prices(variant_id,price,start_date,end_date) values($1,$2,$3,$4)",
		p.VariantId, p.Price, p.StartDate, p.EndDate)
	if err != nil {
		return err
	}
	return nil
}

// Вставка цены варианта продукта в базу
func (r *PostgresProductRepository) AddProductPrice(tx sqlx.Tx, p domain.ProductPrice) error {
	_, err := tx.Exec("insert into product_prices(variant_id,price,start_date) values($1,$2,$3)",
		p.VariantId, p.Price, p.StartDate)
	if err != nil {
		return err
	}
	return nil
}

// Проверка есть ли на скалде продукт
func (r *PostgresProductRepository) CheckProductsInStock(p domain.AddProductInStock) (bool, error) {
	var isExists bool
	err := r.db.Get(&isExists, "select exists(select 1 from products_in_storage where variant_id=$1 and storage_id=$2)",
		p.VariantId, p.StorageId)
	if err != nil {
		return false, err
	}
	return isExists, nil
}

// Обновление колличества продукта
func (r *PostgresProductRepository) UpdateProductsInstock(tx sqlx.Tx, p domain.AddProductInStock) error {
	_, err := tx.Exec("update products_in_storage set quantity=$1 where variant_id=$2 and storage_id=$3",
		p.Quantity, p.VariantId, p.StorageId)
	if err != nil {
		return err
	}
	return nil
}

// Добавление продукта на склад
func (r *PostgresProductRepository) AddProductInStock(tx sqlx.Tx, p domain.AddProductInStock) error {
	_, err := r.db.Exec("insert into products_in_storage(variant_id,storage_id,added_at,quantity) values ($1,$2,$3,$4)",
		p.VariantId, p.StorageId, p.Added_at, p.Quantity)
	if err != nil {
		return err
	}
	return nil
}

// Получени информации о продукте
func (r *PostgresProductRepository) LoadProductInfo(id int) (domain.ProductInfo, error) {
	var p domain.ProductInfo
	err := r.db.Get(&p, "select name,description from products where product_id=$1", id)
	if err != nil {
		return domain.ProductInfo{}, err
	}
	return p, nil
}

// Получение вариантов продукта по его id
func (r *PostgresProductRepository) LoadProductVariants(id int) ([]domain.Variant, error) {
	var v []domain.Variant
	err := r.db.Select(&v, "select variant_id,weight,unit,added_at from product_variants where product_id=$1", id)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// Получение актуальной цены
func (r *PostgresProductRepository) LoadCurrentPrice(variantId int) (decimal.Decimal, error) {
	var price decimal.Decimal
	err := r.db.Get(&price, "select price from product_prices where variant_id=$1 and start_date<$2 and (end_date is null or end_date>$2)",
		variantId, time.Now())
	if err != nil {
		return decimal.Decimal{}, err
	}
	return price, nil
}

// Нахождение id складов в которых находится продукт
func (r *PostgresProductRepository) InStorages(id int) ([]int, error) {
	var inStorages []int
	err := r.db.Select(&inStorages, "SELECT storage_id FROM products_in_storage WHERE variant_id = $1", id)
	if err != nil {
		return nil, err
	}
	return inStorages, nil
}

// Поиск информации о продукте по его тегу
func (r *PostgresProductRepository) FindProductsByTag(tag string, limit int) ([]domain.ProductInfo, error) {
	var products []domain.ProductInfo
	err := r.db.Select(&products, "select product_id,name,description from products where $1 = any (string_to_array(tags,',')) limit $2", tag, limit)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// Получение списка продуктов с лимитом
func (r *PostgresProductRepository) LoadProducts(limit int) ([]domain.ProductInfo, error) {
	var products []domain.ProductInfo
	err := r.db.Select(&products, "select product_id,name,description from products limit $1", limit)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// Получение информации о складах
func (r *PostgresProductRepository) LoadStocks() ([]domain.Stock, error) {
	var stocks []domain.Stock
	err := r.db.Select(&stocks, "select  storage_id,name from storages")
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

// Получение информации о складах где есть определенный продукт
func (r *PostgresProductRepository) FindStocksByProductId(id int) ([]domain.Stock, error) {
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

// Получение вариантов продукта на складе
func (r *PostgresProductRepository) LoadStocksVariants(storageId int) ([]domain.AddProductInStock, error) {
	var v []domain.AddProductInStock
	err := r.db.Select(&v, "select variant_id,storage_id,added_at,quantity from products_in_storage where storage_id=$1 ", storageId)
	if err != nil {
		return nil, err
	}
	return v, err
}

// Получение цены
func (r *PostgresProductRepository) FindPrice(id int) (decimal.Decimal, error) {
	var price decimal.Decimal
	err := r.db.Get(&price, "select price from product_prices where variant_id = $1", id)
	if err != nil {
		return decimal.Decimal{}, nil
	}
	return price, nil
}

// Запись о покупке в базу
func (r *PostgresProductRepository) Buy(tx sqlx.Tx, s domain.Sale) error {
	_, err := tx.Exec("insert into sales(variant_id,storage_id,sold_at,quantity,total_price) values($1,$2,$3,$4,$5)",
		s.VariantId, s.StorageId, s.SoldAt, s.Quantity, s.TotalPrice)
	if err != nil {
		return err
	}
	return nil
}

// Получение списка всех продаж
func (r *PostgresProductRepository) LoadSales(sq domain.SaleQuery) ([]domain.Sale, error) {
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

// Получение списка продаж по фильтрам
func (r *PostgresProductRepository) FindSalesByFilters(sq domain.SaleQuery) ([]domain.Sale, error) {
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
	if err != nil {
		return nil, err
	}
	return sales, nil
}
