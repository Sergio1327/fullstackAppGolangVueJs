package repository

import (
	"database/sql"
	"go-back/internal/app/domain"
	"log"

	"github.com/jmoiron/sqlx"
)

type ProductRepository interface {
	TxBegin() (*sqlx.Tx, error)

	AddProduct(tx *sqlx.Tx, product domain.Product) (int, error)
	AddProductVariantList(tx *sqlx.Tx, productID int, v domain.Variant) error

	CheckExists(tx *sqlx.Tx, p domain.ProductPrice) (int, error)
	UpdateProductPrice(tx *sqlx.Tx, p domain.ProductPrice, id int) error
	AddProductPrice(tx *sqlx.Tx, p domain.ProductPrice) (int, error)

	CheckProductInStock(tx *sqlx.Tx, p domain.AddProductInStock) (bool, error)
	UpdateProductInstock(tx *sqlx.Tx, p domain.AddProductInStock) (int, error)
	AddProductInStock(tx *sqlx.Tx, p domain.AddProductInStock) (int, error)

	LoadProductInfo(tx *sqlx.Tx, productID int) (domain.ProductInfo, error)
	FindProductVariantList(tx *sqlx.Tx, productID int) ([]domain.Variant, error)
	FindCurrentPrice(tx *sqlx.Tx, variantID int) (float64, error)
	InStorages(tx *sqlx.Tx, variantID int) ([]int, error)

	FindProductListByTag(tx *sqlx.Tx, tag string, limit int) ([]domain.ProductInfo, error)
	LoadProductList(tx *sqlx.Tx, limit int) ([]domain.ProductInfo, error)

	LoadStockList(tx *sqlx.Tx) ([]domain.Stock, error)
	FindStockListByProductId(tx *sqlx.Tx, productID int) ([]domain.Stock, error)
	FindStocksVariantList(tx *sqlx.Tx, storageID int) ([]domain.AddProductInStock, error)

	Buy(tx *sqlx.Tx, s domain.Sale) (int, error)
	FindPrice(tx *sqlx.Tx, variantID int) (float64, error)

	FindSaleListOnlyBySoldDate(tx *sqlx.Tx, sq domain.SaleQueryOnlyBySoldDate) ([]domain.Sale, error)
	FindSaleListByFilters(tx *sqlx.Tx, sq domain.SaleQuery) ([]domain.Sale, error)
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
func (r *PostgresProductRepository) AddProduct(tx *sqlx.Tx, product domain.Product) (productID int, err error) {
	err = tx.QueryRow(`
	insert into products
	(name, description, added_at, tags)
	values ($1, $2, $3, $4) 
	returning product_id`,
		product.Name, product.Descr, product.AddetAt, product.Tags).Scan(&productID)

	return productID, err
}

// AddProductVariantList добавление вариантов продукта в продукт по его id
func (r *PostgresProductRepository) AddProductVariantList(tx *sqlx.Tx, productID int, variant domain.Variant) error {
	_, err := tx.Exec(`
	insert into product_variants 
	(product_id, weight, unit) 
	values ($1, $2, $3)`, productID, variant.Weight, variant.Unit)

	return err
}

// CheckExists проверка наличия цен варианта продукта в указаный диапазон времени
func (r *PostgresProductRepository) CheckExists(tx *sqlx.Tx, p domain.ProductPrice) (isExists int, err error) {
	err = tx.Get(&isExists,
		`select price_id 
		 from product_prices
		 where variant_id = $1 
		 and start_date = $2 
		 and( end_date = $3 or end_date is null )`,
		p.VariantID, p.StartDate, p.EndDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return isExists, err
}

// UpdateProductPrice обновление цены варианта продукта
func (r *PostgresProductRepository) UpdateProductPrice(tx *sqlx.Tx, price domain.ProductPrice, priceID int) error {
	_, err := tx.Exec(`
	update product_prices
	set end_date = $1 
	where price_id = $2`,
		price.EndDate, priceID)

	return err
}

// AddProductPrice вставка цены варианта продукта в базу
func (r *PostgresProductRepository) AddProductPrice(tx *sqlx.Tx, price domain.ProductPrice) (priceID int, err error) {
	err = tx.QueryRow(`
	insert into product_prices
	( variant_id, price, start_date, end_date )
	values( $1, $2, $3, $4 )
	returning price_id`,
		price.VariantID, price.Price, price.StartDate, price.EndDate).Scan(&priceID)

	return priceID, err
}

// CheckProductInStock проверка есть ли на скалде продукт
func (r *PostgresProductRepository) CheckProductInStock(tx *sqlx.Tx, productInStock domain.AddProductInStock) (isExists bool, err error) {
	err = tx.Get(&isExists,
		`select exists
		 (select 1 
		 from products_in_storage 
		 where variant_id = $1 
		 and storage_id = $2)`,
		productInStock.VariantID, productInStock.StorageID)

	return isExists, err
}

// UpdateProductInstock обновление колличества продукта
func (r *PostgresProductRepository) UpdateProductInstock(tx *sqlx.Tx, productInStock domain.AddProductInStock) (productStockID int, err error) {
	err = tx.QueryRow(`
	update products_in_storage 
	set quantity = $1
	where variant_id = $2 
	and storage_id = $3
	returning pis_id`,
		productInStock.Quantity, productInStock.VariantID, productInStock.StorageID).Scan(&productStockID)

	return productStockID, err
}

// AddProductInStock добавление продукта на склад
func (r *PostgresProductRepository) AddProductInStock(tx *sqlx.Tx, productInStock domain.AddProductInStock) (productStockID int, err error) {
	err = tx.QueryRow(`
	 insert into products_in_storage
	 ( variant_id, storage_id, added_at, quantity )
	 values ($1, $2, $3, $4)
	 returning pis_id`,
		productInStock.VariantID, productInStock.StorageID, productInStock.AddedAt, productInStock.Quantity).Scan(&productStockID)

	return productStockID, err
}

// LoadProductInfo получение информации о продукте
func (r *PostgresProductRepository) LoadProductInfo(tx *sqlx.Tx, productId int) (productInfo domain.ProductInfo, err error) {
	err = tx.Get(&productInfo,
		`select product_id, name, description  
	 	 from products 
	     where product_id = $1`, productId)

	return productInfo, err
}

// FindProductVariantList получение вариантов продукта по его id
func (r *PostgresProductRepository) FindProductVariantList(tx *sqlx.Tx, productID int) (variantList []domain.Variant, err error) {
	err = tx.Select(&variantList,
		`select product_id, variant_id, weight, unit, added_at
		 from product_variants	
		 where product_id = $1`, productID)

	return variantList, err
}

// FindCurrentPrice получение актуальной цены
func (r *PostgresProductRepository) FindCurrentPrice(tx *sqlx.Tx, variantID int) (price float64, err error) {
	err = tx.Get(&price,
		`select price 
		 from product_prices 
		 where variant_id = $1 
		 and start_date < now() 
		 and ( end_date is null or end_date > now() )`,
		variantID)

	return price, err
}

// InStorages нахождение id складов в которых находится продукт
func (r *PostgresProductRepository) InStorages(tx *sqlx.Tx, varantID int) (inStorages []int, err error) {
	err = tx.Select(&inStorages,
		`SELECT storage_id 
	 	 FROM products_in_storage 
		 WHERE variant_id = $1`, varantID)

	return inStorages, err
}

// FindProductListByTag  поиск информации о продукте по его тегу
func (r *PostgresProductRepository) FindProductListByTag(tx *sqlx.Tx, tag string, limit int) (productList []domain.ProductInfo, err error) {
	err = tx.Select(&productList,
		`select product_id, name, description
	 	 from products 
	 	 where $1 = any ( string_to_array( tags,',' )) 
	 	 limit $2`,
		tag, limit)

	return productList, err
}

// LoadProductList получение списка продуктов с лимитом
func (r *PostgresProductRepository) LoadProductList(tx *sqlx.Tx, limit int) (productList []domain.ProductInfo, err error) {
	err = tx.Select(&productList,
		`select product_id, name, description
	 	 from products
	     limit $1`, limit)

	return productList, err
}

// LoadStockList получение информации о складах
func (r *PostgresProductRepository) LoadStockList(tx *sqlx.Tx) (stockList []domain.Stock, err error) {
	err = tx.Select(&stockList,
		`select  storage_id, name
		 from storages`)

	return stockList, err
}

// FindStockListByProductId получение информации о складах где есть определенный продукт
func (r *PostgresProductRepository) FindStockListByProductId(tx *sqlx.Tx, productID int) (stockList []domain.Stock, err error) {
	err = tx.Select(&stockList, `
	select s.storage_id ,s.name 
	from storages s
	join products_in_storage pis ON (s.storage_id = pis.storage_id)
	join product_variants pv ON (pis.variant_id = pv.variant_id)
	join products p ON (pv.product_id = p.product_id)
	where p.product_id = $1`, productID)

	return stockList, err
}

// FindStocksVariantList получение вариантов продукта на складе
func (r *PostgresProductRepository) FindStocksVariantList(tx *sqlx.Tx, storageID int) (variantList []domain.AddProductInStock, err error) {
	err = tx.Select(&variantList,
		`select variant_id, storage_id, added_at, quantity
	     from products_in_storage 
	     where storage_id = $1 `, storageID)

	return variantList, err
}

// FindPrice получение цены
func (r *PostgresProductRepository) FindPrice(tx *sqlx.Tx, variantID int) (price float64, err error) {
	err = tx.Get(&price,
		`select price
	 	 from product_prices
	  	 where variant_id = $1`, variantID)

	return price, err
}

// Buy запись о покупке в базу
func (r *PostgresProductRepository) Buy(tx *sqlx.Tx, sale domain.Sale) (saleID int, err error) {
	err = tx.QueryRow(`
	insert into sales
	( variant_id, storage_id, sold_at, quantity, total_price )
	values( $1, $2, $3, $4, $5 )
	returning sales_id`,
		sale.VariantID, sale.StorageID, sale.SoldAt, sale.Quantity, sale.TotalPrice).Scan(&saleID)

	return saleID, err
}

// FindSaleListOnlyBySoldDate получение списка всех продаж
func (r *PostgresProductRepository) FindSaleListOnlyBySoldDate(tx *sqlx.Tx, saleFilters domain.SaleQueryOnlyBySoldDate) (saleList []domain.Sale, err error) {
	query := `
	SELECT s.sales_id, s.variant_id, s.storage_id, s.sold_at, s.quantity, s.total_price, p.name 
	FROM sales s
	JOIN product_variants  pv ON ( pv.variant_id = s.variant_id )
	JOIN products  p ON ( p.product_id = pv.product_id )
	WHERE s.sold_at >= $1 AND s.sold_at <= $2
	LIMIT $3`

	err = tx.Select(&saleList, query, saleFilters.StartDate, saleFilters.EndDate, saleFilters.Limit)

	return saleList, err
}

// FindSaleListByFilters получение списка продаж по фильтрам
func (r *PostgresProductRepository) FindSaleListByFilters(tx *sqlx.Tx, saleFilters domain.SaleQuery) (saleList []domain.Sale, err error) {
	query := `
	SELECT s.sales_id, s.variant_id, s.storage_id, s.sold_at, s.quantity, s.total_price, p.name 
	FROM sales s
	JOIN product_variants pv ON (pv.variant_id = s.variant_id)
	JOIN products p ON (p.product_id = pv.product_id)
	WHERE s.sold_at > :start_date AND s.sold_at < :end_date
	AND ( cast(:product_name as varchar) IS NULL OR p.name = :product_name )
	AND ( cast(:storage_id as integer) IS NULL OR s.storage_id = :storage_id ) 
	LIMIT :limit`

	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	defer stmt.Close()

	err = stmt.Select(&saleList, saleFilters)

	return saleList, err
}
