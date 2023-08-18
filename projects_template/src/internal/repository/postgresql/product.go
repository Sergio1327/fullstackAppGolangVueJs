package postgresql

import (
	"product_storage/internal/entity/product"
	"product_storage/internal/entity/stock"
	"product_storage/internal/repository"
	"product_storage/internal/transaction"
	"product_storage/tools/gensql"
)

type productRepository struct {
}

func NewProduct() repository.Product {
	return &productRepository{}

}

// AddProduct вставка названия,описания,времени добавления и тегов в базу
func (r *productRepository) AddProduct(ts transaction.Session, product product.ProductParams) (productID int, err error) {
	query := `insert into products
	(name, description, added_at, tags)
	values ($1, $2, $3, $4) 
	returning product_id`

	err = SqlxTx(ts).QueryRow(query, product.Name, product.Descr, product.AddetAt, product.Tags).Scan(&productID)

	return productID, err
}

// AddProductVariantList добавление вариантов продукта в продукт по его id
func (r *productRepository) AddProductVariantList(ts transaction.Session, productID int, variant product.Variant) error {
	query := `
	insert into product_variants 
	(product_id, weight, unit) 
	values ($1, $2, $3)`

	_, err := SqlxTx(ts).Exec(query, productID, variant.Weight, variant.Unit)
	return err
}

// CheckExists проверка наличия цен варианта продукта в указаный диапазон времени
func (r *productRepository) CheckExists(ts transaction.Session, p product.ProductPriceParams) (isExistsID int, err error) {
	query := `
	select price_id 
	from product_prices
	where variant_id = $1 
	and start_date = $2 
	and( end_date = $3 or end_date is null )`

	return gensql.Get[int](SqlxTx(ts), query, p.VariantID, p.StartDate, p.EndDate)
}

// UpdateProductPrice обновление цены варианта продукта
func (r *productRepository) UpdateProductPrice(ts transaction.Session, price product.ProductPriceParams, priceID int) error {
	_, err := SqlxTx(ts).Exec(`
	update product_prices
	set end_date = $1 
	where price_id = $2`,
		price.EndDate, priceID)

	return err
}

// AddProductPrice вставка цены варианта продукта в базу
func (r *productRepository) AddProductPrice(ts transaction.Session, price product.ProductPriceParams) (priceID int, err error) {
	err = SqlxTx(ts).QueryRow(`
	insert into product_prices
	( variant_id, price, start_date, end_date )
	values( $1, $2, $3, $4 )
	returning price_id`,
		price.VariantID, price.Price, price.StartDate, price.EndDate).Scan(&priceID)

	return priceID, err
}

// CheckProductInStock проверка есть ли на скалде продукт
func (r *productRepository) CheckProductInStock(ts transaction.Session, productInStock stock.ProductInStockParams) (isExists bool, err error) {
	query := `select exists
	(select 1 
	from products_in_storage 
	where variant_id = $1 
	and storage_id = $2)`

	return gensql.Get[bool](SqlxTx(ts), query, productInStock.VariantID, productInStock.StorageID)
}

// UpdateProductInstock обновление колличества продукта
func (r *productRepository) UpdateProductInstock(ts transaction.Session, productInStock stock.ProductInStockParams) (productStockID int, err error) {
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
func (r *productRepository) AddProductInStock(ts transaction.Session, productInStock stock.ProductInStockParams) (productStockID int, err error) {
	err = SqlxTx(ts).QueryRow(`
	 insert into products_in_storage
	 ( variant_id, storage_id, added_at, quantity )
	 values ($1, $2, $3, $4)
	 returning pis_id`,
		productInStock.VariantID, productInStock.StorageID, productInStock.AddedAt, productInStock.Quantity).Scan(&productStockID)

	return productStockID, err
}

// LoadProductInfo получение информации о продукте
func (r *productRepository) LoadProductInfo(ts transaction.Session, productId int) (productInfo product.ProductInfo, err error) {
	query := `
	select product_id, name, description  
	from products 
    where product_id = $1`

	return gensql.Get[product.ProductInfo](SqlxTx(ts), query, productId)
}

// FindProductVariantList получение вариантов продукта по его id
func (r *productRepository) FindProductVariantList(ts transaction.Session, productID int) (variantList []product.Variant, err error) {
	query := `
	select product_id, variant_id, weight, unit, added_at
	from product_variants	
	where product_id = $1`

	return gensql.Select[product.Variant](SqlxTx(ts), query, productID)
}

// FindCurrentPrice получение актуальной цены
func (r *productRepository) FindCurrentPrice(ts transaction.Session, variantID int) (price float64, err error) {
	query := `
	select price 
	from product_prices 
	where variant_id = $1 
	and start_date < now() 
	and ( end_date is null or end_date > now() )`

	return gensql.Get[float64](SqlxTx(ts), query, variantID)
}

// InStorages нахождение id складов в которых находится продукт
func (r *productRepository) InStorages(ts transaction.Session, varantID int) (inStorages []int, err error) {
	query := `
	SELECT storage_id 
	FROM products_in_storage 
    WHERE variant_id = $1`

	return gensql.Select[int](SqlxTx(ts), query, varantID)
}

// FindProductListByTag  поиск информации о продукте по его тегу
func (r *productRepository) FindProductListByTag(ts transaction.Session, tag string, limit int) (productList []product.ProductInfo, err error) {
	query := `
	select product_id, name, description
	from products 
	where $1 = any ( string_to_array( tags,',' )) 
	limit $2`

	return gensql.Select[product.ProductInfo](SqlxTx(ts), query, tag, limit)
}

// LoadProductList получение списка продуктов с лимитом
func (r *productRepository) LoadProductList(ts transaction.Session, limit int) (productList []product.ProductInfo, err error) {
	query := `
	select product_id, name, description
	from products
    limit $1`

	return gensql.Select[product.ProductInfo](SqlxTx(ts), query, limit)
}

// LoadStockList получение информации о складах
func (r *productRepository) LoadStockList(ts transaction.Session) (stockList []stock.Stock, err error) {
	query := `
	select  storage_id, name
	from storages`

	return gensql.Select[stock.Stock](SqlxTx(ts), query)
}

// FindStockListByProductId получение информации о складах где есть определенный продукт
func (r *productRepository) FindStockListByProductId(ts transaction.Session, productID int) (stockList []stock.Stock, err error) {
	query := `
	select s.storage_id ,s.name 
	from storages s
	join products_in_storage pis ON (s.storage_id = pis.storage_id)
	join product_variants pv ON (pis.variant_id = pv.variant_id)
	join products p ON (pv.product_id = p.product_id)
	where p.product_id = $1`

	return gensql.Select[stock.Stock](SqlxTx(ts), query, productID)
}

// FindStocksVariantList получение вариантов продукта на складе
func (r *productRepository) FindStocksVariantList(ts transaction.Session, storageID int) (variantList []stock.ProductInStockParams, err error) {
	query := `
	select variant_id, storage_id, added_at, quantity
	from products_in_storage 
	where storage_id = $1 `

	return gensql.Select[stock.ProductInStockParams](SqlxTx(ts), query, storageID)
}

// FindPrice получение цены
func (r *productRepository) FindPrice(ts transaction.Session, variantID int) (price float64, err error) {
	query :=
		`select price
	 	 from product_prices
	 	 where variant_id = $1`

	return gensql.Get[float64](SqlxTx(ts), query, variantID)
}

// SaveSale запись о покупке в базу
func (r *productRepository) SaveSale(ts transaction.Session, sale product.SaleParams) (saleID int, err error) {
	err = SqlxTx(ts).QueryRow(`
	insert into sales
	( variant_id, storage_id, sold_at, quantity, total_price )
	values( $1, $2, $3, $4, $5 )
	returning sales_id`,
		sale.VariantID, sale.StorageID, sale.SoldAt, sale.Quantity, sale.TotalPrice).Scan(&saleID)

	return saleID, err
}

// FindSaleListOnlyBySoldDate получение списка всех продаж
func (r *productRepository) FindSaleListOnlyBySoldDate(ts transaction.Session, saleFilters product.SaleQueryOnlyBySoldDateParam) (saleList []product.Sale, err error) {
	query := `
	SELECT s.sales_id, s.variant_id, s.storage_id, s.sold_at, s.quantity, s.total_price, p.name 
	FROM sales s
	JOIN product_variants  pv ON ( pv.variant_id = s.variant_id )
	JOIN products  p ON ( p.product_id = pv.product_id )
	WHERE s.sold_at >= $1 AND s.sold_at <= $2
	LIMIT $3`

	return gensql.Select[product.Sale](SqlxTx(ts), query, saleFilters.StartDate, saleFilters.EndDate, saleFilters.Limit)
}

// FindSaleListByFilters получение списка продаж по фильтрам
func (r *productRepository) FindSaleListByFilters(ts transaction.Session, saleFilters product.SaleQueryParam) (saleList []product.Sale, err error) {
	query := `
	SELECT s.sales_id, s.variant_id, s.storage_id, s.sold_at, s.quantity, s.total_price, p.name 
	FROM sales s
	JOIN product_variants pv ON (pv.variant_id = s.variant_id)
	JOIN products p ON (p.product_id = pv.product_id)
	WHERE s.sold_at > :start_date AND s.sold_at < :end_date
	AND ( cast(:product_name as varchar) IS NULL OR p.name = :product_name )
	AND ( cast(:storage_id as integer) IS NULL OR s.storage_id = :storage_id ) 
	LIMIT :limit`

	params := map[string]interface{}{
		"start_date":   saleFilters.StartDate,
		"end_date":     saleFilters.EndDate,
		"limit":        saleFilters.Limit,
		"product_name": saleFilters.ProductName,
		"storage_id":   saleFilters.StorageID,
	}

	return gensql.SelectNamed[product.Sale](SqlxTx(ts), query, params)
}

func (r productRepository) AddStock(ts transaction.Session, storage stock.StockParams) (stockID int, err error) {
	query := `
	insert into storages
	(name, added_at)
	values ($1, $2)
	returning storage_id
	`
	err = SqlxTx(ts).QueryRow(query, storage.StorageName, storage.Added_at).Scan(&stockID)
	return stockID, err
}
