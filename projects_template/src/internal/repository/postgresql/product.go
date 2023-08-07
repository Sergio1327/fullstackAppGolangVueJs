package postgresql

import (
	"database/sql"
	"product_storage/internal/entity/product"
	"product_storage/internal/entity/stock"
	"product_storage/internal/transaction"
)

type PostgresProductRepository struct {
}

func NewProduct() PostgresProductRepository {
	return PostgresProductRepository{}

}

// AddProduct вставка названия,описания,времени добавления и тегов в базу
func (r PostgresProductRepository) AddProduct(ts transaction.Session, product product.Product) (productID int, err error) {
	query := `insert into products
	(name, description, added_at, tags)
	values ($1, $2, $3, $4) 
	returning product_id`

	err = SqlxTx(ts).QueryRow(query, product.Name, product.Descr, product.AddetAt, product.Tags).Scan(&productID)

	return productID, err
}

// AddProductVariantList добавление вариантов продукта в продукт по его id
func (r PostgresProductRepository) AddProductVariantList(ts transaction.Session, productID int, variant product.Variant) error {
	query := `
	insert into product_variants 
	(product_id, weight, unit) 
	values ($1, $2, $3)`

	_, err := SqlxTx(ts).Exec(query, productID, variant.Weight, variant.Unit)
	return err
}

// CheckExists проверка наличия цен варианта продукта в указаный диапазон времени
func (r PostgresProductRepository) CheckExists(ts transaction.Session, p product.ProductPrice) (isExistsID int, err error) {
	query := `
	select price_id 
	from product_prices
	where variant_id = $1 
	and start_date = $2 
	and( end_date = $3 or end_date is null )`

	err = SqlxTx(ts).Get(&isExistsID, query, p.VariantID, p.StartDate, p.EndDate)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return 0, nil
		default:
			return 0, err
		}
	}

	return isExistsID, err
}

// UpdateProductPrice обновление цены варианта продукта
func (r PostgresProductRepository) UpdateProductPrice(ts transaction.Session, price product.ProductPrice, priceID int) error {
	_, err := SqlxTx(ts).Exec(`
	update product_prices
	set end_date = $1 
	where price_id = $2`,
		price.EndDate, priceID)

	return err
}

// AddProductPrice вставка цены варианта продукта в базу
func (r PostgresProductRepository) AddProductPrice(ts transaction.Session, price product.ProductPrice) (priceID int, err error) {
	err = SqlxTx(ts).QueryRow(`
	insert into product_prices
	( variant_id, price, start_date, end_date )
	values( $1, $2, $3, $4 )
	returning price_id`,
		price.VariantID, price.Price, price.StartDate, price.EndDate).Scan(&priceID)

	return priceID, err
}

// CheckProductInStock проверка есть ли на скалде продукт
func (r PostgresProductRepository) CheckProductInStock(ts transaction.Session, productInStock stock.AddProductInStock) (isExists bool, err error) {
	err = SqlxTx(ts).Get(&isExists,
		`select exists
		 (select 1 
		 from products_in_storage 
		 where variant_id = $1 
		 and storage_id = $2)`,
		productInStock.VariantID, productInStock.StorageID)

	return isExists, err
}

// UpdateProductInstock обновление колличества продукта
func (r PostgresProductRepository) UpdateProductInstock(ts transaction.Session, productInStock stock.AddProductInStock) (productStockID int, err error) {
	err = SqlxTx(ts).QueryRow(`
	update products_in_storage 
	set quantity = $1
	where variant_id = $2 
	and storage_id = $3
	returning pis_id`,
		productInStock.Quantity, productInStock.VariantID, productInStock.StorageID).Scan(&productStockID)

	return productStockID, err
}

// AddProductInStock добавление продукта на склад
func (r PostgresProductRepository) AddProductInStock(ts transaction.Session, productInStock stock.AddProductInStock) (productStockID int, err error) {
	err = SqlxTx(ts).QueryRow(`
	 insert into products_in_storage
	 ( variant_id, storage_id, added_at, quantity )
	 values ($1, $2, $3, $4)
	 returning pis_id`,
		productInStock.VariantID, productInStock.StorageID, productInStock.AddedAt, productInStock.Quantity).Scan(&productStockID)

	return productStockID, err
}

// LoadProductInfo получение информации о продукте
func (r PostgresProductRepository) LoadProductInfo(ts transaction.Session, productId int) (productInfo product.ProductInfo, err error) {
	err = SqlxTx(ts).Get(&productInfo,
		`select product_id, name, description  
	 	 from products 
	     where product_id = $1`, productId)

	return productInfo, err
}

// FindProductVariantList получение вариантов продукта по его id
func (r PostgresProductRepository) FindProductVariantList(ts transaction.Session, productID int) (variantList []product.Variant, err error) {
	err = SqlxTx(ts).Select(&variantList,
		`select product_id, variant_id, weight, unit, added_at
		 from product_variants	
		 where product_id = $1`, productID)

	return variantList, err
}

// FindCurrentPrice получение актуальной цены
func (r PostgresProductRepository) FindCurrentPrice(ts transaction.Session, variantID int) (price float64, err error) {
	err = SqlxTx(ts).Get(&price,
		`select price 
		 from product_prices 
		 where variant_id = $1 
		 and start_date < now() 
		 and ( end_date is null or end_date > now() )`,
		variantID)

	return price, err
}

// InStorages нахождение id складов в которых находится продукт
func (r PostgresProductRepository) InStorages(ts transaction.Session, varantID int) (inStorages []int, err error) {
	err = SqlxTx(ts).Select(&inStorages,
		`SELECT storage_id 
	 	 FROM products_in_storage 
		 WHERE variant_id = $1`, varantID)

	return inStorages, err
}

// FindProductListByTag  поиск информации о продукте по его тегу
func (r PostgresProductRepository) FindProductListByTag(ts transaction.Session, tag string, limit int) (productList []product.ProductInfo, err error) {
	err = SqlxTx(ts).Select(&productList,
		`select product_id, name, description
	 	 from products 
	 	 where $1 = any ( string_to_array( tags,',' )) 
	 	 limit $2`,
		tag, limit)

	return productList, err
}

// LoadProductList получение списка продуктов с лимитом
func (r PostgresProductRepository) LoadProductList(ts transaction.Session, limit int) (productList []product.ProductInfo, err error) {
	err = SqlxTx(ts).Select(&productList,
		`select product_id, name, description
	 	 from products
	     limit $1`, limit)

	return productList, err
}

// LoadStockList получение информации о складах
func (r PostgresProductRepository) LoadStockList(ts transaction.Session) (stockList []stock.Stock, err error) {
	err = SqlxTx(ts).Select(&stockList,
		`select  storage_id, name
		 from storages`)

	return stockList, err
}

// FindStockListByProductId получение информации о складах где есть определенный продукт
func (r PostgresProductRepository) FindStockListByProductId(ts transaction.Session, productID int) (stockList []stock.Stock, err error) {
	err = SqlxTx(ts).Select(&stockList, `
	select s.storage_id ,s.name 
	from storages s
	join products_in_storage pis ON (s.storage_id = pis.storage_id)
	join product_variants pv ON (pis.variant_id = pv.variant_id)
	join products p ON (pv.product_id = p.product_id)
	where p.product_id = $1`, productID)

	return stockList, err
}

// FindStocksVariantList получение вариантов продукта на складе
func (r PostgresProductRepository) FindStocksVariantList(ts transaction.Session, storageID int) (variantList []stock.AddProductInStock, err error) {
	err = SqlxTx(ts).Select(&variantList,
		`select variant_id, storage_id, added_at, quantity
	     from products_in_storage 
	     where storage_id = $1 `, storageID)

	return variantList, err
}

// FindPrice получение цены
func (r PostgresProductRepository) FindPrice(ts transaction.Session, variantID int) (price float64, err error) {
	err = SqlxTx(ts).Get(&price,
		`select price
	 	 from product_prices
	  	 where variant_id = $1`, variantID)

	return price, err
}

// Buy запись о покупке в базу
func (r PostgresProductRepository) Buy(ts transaction.Session, sale product.Sale) (saleID int, err error) {
	err = SqlxTx(ts).QueryRow(`
	insert into sales
	( variant_id, storage_id, sold_at, quantity, total_price )
	values( $1, $2, $3, $4, $5 )
	returning sales_id`,
		sale.VariantID, sale.StorageID, sale.SoldAt, sale.Quantity, sale.TotalPrice).Scan(&saleID)

	return saleID, err
}

// FindSaleListOnlyBySoldDate получение списка всех продаж
func (r PostgresProductRepository) FindSaleListOnlyBySoldDate(ts transaction.Session, saleFilters product.SaleQueryOnlyBySoldDate) (saleList []product.Sale, err error) {
	query := `
	SELECT s.sales_id, s.variant_id, s.storage_id, s.sold_at, s.quantity, s.total_price, p.name 
	FROM sales s
	JOIN product_variants  pv ON ( pv.variant_id = s.variant_id )
	JOIN products  p ON ( p.product_id = pv.product_id )
	WHERE s.sold_at >= $1 AND s.sold_at <= $2
	LIMIT $3`

	err = SqlxTx(ts).Select(&saleList, query, saleFilters.StartDate, saleFilters.EndDate, saleFilters.Limit)

	return saleList, err
}

// FindSaleListByFilters получение списка продаж по фильтрам
func (r PostgresProductRepository) FindSaleListByFilters(ts transaction.Session, saleFilters product.SaleQuery) (saleList []product.Sale, err error) {
	query := `
	SELECT s.sales_id, s.variant_id, s.storage_id, s.sold_at, s.quantity, s.total_price, p.name 
	FROM sales s
	JOIN product_variants pv ON (pv.variant_id = s.variant_id)
	JOIN products p ON (p.product_id = pv.product_id)
	WHERE s.sold_at > :start_date AND s.sold_at < :end_date
	AND ( cast(:product_name as varchar) IS NULL OR p.name = :product_name )
	AND ( cast(:storage_id as integer) IS NULL OR s.storage_id = :storage_id ) 
	LIMIT :limit`

	stmt, err := SqlxTx(ts).PrepareNamed(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.Select(&saleList, saleFilters)

	return saleList, err
}
