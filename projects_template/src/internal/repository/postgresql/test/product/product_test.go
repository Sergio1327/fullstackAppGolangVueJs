package product_test

import (
	"product_storage/internal/entity/params"
	"product_storage/internal/entity/product"
	"product_storage/internal/entity/stock"
	"product_storage/internal/repository/postgresql"
	"product_storage/internal/transaction"
	"product_storage/rimport"
	"product_storage/tools/pgdb"
	"product_storage/tools/sqlnull"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAddProduct(t *testing.T) {
	r := require.New(t)

	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()
	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	ts := sm.CreateSession()
	ts.Start()
	defer ts.Rollback()

	// Первый тестовый случай - успешное добавление продукта
	p1 := product.Product{
		Name:    "sdlds",
		Descr:   "dlldld",
		AddetAt: time.Now(),
		Tags:    "ddd",
	}

	id, err := repo.Repository.Product.AddProduct(ts, p1)
	r.NoError(err)
	r.NotEmpty(id)

	productInfo, err := repo.Repository.Product.LoadProductInfo(ts, id)
	r.NoError(err)
	r.NotEmpty(productInfo)

	r.Equal(p1.Name, productInfo.Name)
	r.Equal(p1.Descr, productInfo.Descr)

	// Второй тестовый случай - добавление продукта без имени, ожидается ошибка
	p2 := product.Product{}

	id2, err := repo.Repository.Product.AddProduct(ts, p2)
	r.NoError(err)
	r.NotEmpty(id2)
}

func TestAddProductVariantList(t *testing.T) {
	r := require.New(t)
	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()
	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	ts := sm.CreateSession()
	ts.Start()
	defer ts.Rollback()

	id := 1
	variant := product.Variant{
		Weight: 200,
		Unit:   "г",
	}
	err := repo.Repository.Product.AddProductVariantList(ts, id, variant)
	r.NoError(err)

	id = -1
	variant = product.Variant{
		Weight: 300,
		Unit:   "кг",
	}
	err = repo.Repository.Product.AddProductVariantList(ts, id, variant)
	r.Error(err)
}

func TestCheckExists(t *testing.T) {
	r := require.New(t)
	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()
	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	ts := sm.CreateSession()
	ts.Start()
	defer ts.Rollback()

	productPrice := product.ProductPrice{
		VariantID: 2,
		StartDate: time.Now(),
	}

	id, err := repo.Repository.Product.CheckExists(ts, productPrice)
	r.NoError(err)
	r.Zero(id)
}

func TestUpdateProductPrice(t *testing.T) {
	r := require.New(t)
	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()
	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	ts := sm.CreateSession()
	ts.Start()
	defer ts.Rollback()

	expectedProductPrice := product.ProductPrice{
		PriceID: 2,
		EndDate: sqlnull.NewNullTime(time.Now()),
	}

	err := repo.Repository.Product.UpdateProductPrice(ts, expectedProductPrice, expectedProductPrice.PriceID)
	r.NoError(err)

	var productPrice product.ProductPrice
	err = postgresql.SqlxTx(ts).Get(&productPrice,
		`select price_id 
	     from product_prices
		 where price_id = $1 
		 and end_date = $2`, expectedProductPrice.PriceID, expectedProductPrice.EndDate)

	r.NoError(err)
	r.Equal(expectedProductPrice.PriceID, productPrice.PriceID)
}

func TestAddProductPrice(t *testing.T) {
	r := require.New(t)

	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()
	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	ts := sm.CreateSession()
	ts.Start()
	defer ts.Rollback()

	startDate, err := time.Parse("02.01.2006", "01.07.2023")
	r.NoError(err)

	expectedProductPrice := product.ProductPrice{
		VariantID: 5,
		StartDate: startDate,
		Price:     18.99,
	}

	priceID, err := repo.Repository.Product.AddProductPrice(ts, expectedProductPrice)
	r.NoError(err)
	r.NotEmpty(priceID)

	var productPrice product.ProductPrice
	err = postgresql.SqlxTx(ts).Get(&productPrice,
		`select variant_id, price
		 from product_prices
		 where variant_id = $1
		 and start_date = $2 
		 and price = $3`, expectedProductPrice.VariantID, expectedProductPrice.StartDate, expectedProductPrice.Price)
	r.NoError(err)

	r.Equal(expectedProductPrice.VariantID, productPrice.VariantID)
	r.Equal(expectedProductPrice.Price, productPrice.Price)
}

func TestCheckProductsInStock(t *testing.T) {
	r := require.New(t)

	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()
	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	ts := sm.CreateSession()
	ts.Start()
	defer ts.Rollback()

	productInStock := stock.AddProductInStock{
		VariantID: 4,
		StorageID: 2,
		Quantity:  10,
	}

	_, err := postgresql.SqlxTx(ts).Exec(
		`insert into products_in_storage
		 ( variant_id, storage_id, quantity ) 
		 values( $1, $2, $3 )`, productInStock.VariantID, productInStock.StorageID, productInStock.Quantity)
	r.NoError(err)

	isExists, err := repo.Repository.Product.CheckProductInStock(ts, productInStock)
	r.NoError(err)
	r.True(isExists)
}

func TestUpdateProductsnStock(t *testing.T) {
	r := require.New(t)

	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()
	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	ts := sm.CreateSession()
	ts.Start()
	defer ts.Rollback()

	productInStock := stock.AddProductInStock{
		VariantID: 4,
		StorageID: 2,
		Quantity:  10,
	}

	_, err := postgresql.SqlxTx(ts).Exec(
		`insert into products_in_storage
		 ( variant_id, storage_id, quantity )
		 values ( $1, $2, $3 )`, productInStock.VariantID, productInStock.StorageID, productInStock.Quantity)
	r.NoError(err)

	productStockID, err := repo.Repository.Product.UpdateProductInstock(ts, stock.AddProductInStock{
		VariantID: 4,
		StorageID: 2,
		Quantity:  5,
	})
	r.NoError(err)
	r.NotEmpty(productStockID)
}

func TestAddProductInStock(t *testing.T) {
	r := require.New(t)

	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()
	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	ts := sm.CreateSession()
	ts.Start()
	defer ts.Rollback()

	expectedProduct := stock.AddProductInStock{
		VariantID: 3,
		StorageID: 1,
		AddedAt:   time.Now(),
		Quantity:  5,
	}
	productStockID, err := repo.Repository.Product.AddProductInStock(ts, expectedProduct)
	r.NotZero(productStockID)
	r.NoError(err)

	var productInStock stock.AddProductInStock
	err = postgresql.SqlxTx(ts).Get(&productInStock, `
	select variant_id, storage_id, quantity
	from products_in_storage
	where variant_id = $1 and storage_id = $2 and quantity = $3`,
		expectedProduct.VariantID, expectedProduct.StorageID, expectedProduct.Quantity)
	r.NoError(err)

	r.Equal(expectedProduct.VariantID, productInStock.VariantID)
	r.Equal(expectedProduct.StorageID, productInStock.StorageID)
	r.Equal(expectedProduct.Quantity, productInStock.Quantity)
}

func TestLoadProductInfo(t *testing.T) {
	r := require.New(t)

	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()
	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	ts := sm.CreateSession()
	ts.Start()
	defer ts.Rollback()

	product := product.Product{
		Name:    "Имя продукта",
		Descr:   "Описание продукта",
		AddetAt: time.Now(),
		Tags:    "tag",
	}

	id, err := repo.Repository.Product.AddProduct(ts, product)
	r.NoError(err)
	r.NotEmpty(id)

	productInfo, err := repo.Repository.Product.LoadProductInfo(ts, id)
	r.NoError(err)
	r.NotEmpty(productInfo)
	r.Equal(product.Name, productInfo.Name)
	r.Equal(product.Descr, productInfo.Descr)
}

func TestFindProductVariantList(t *testing.T) {
	r := require.New(t)

	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()
	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	ts := sm.CreateSession()
	ts.Start()
	defer ts.Rollback()

	varquery := product.Variant{
		Weight: 500,
		Unit:   "г",
	}

	id := 2
	err := repo.Repository.Product.AddProductVariantList(ts, id, varquery)
	r.NoError(err)

	variants, err := repo.Repository.Product.FindProductVariantList(ts, id)
	r.NoError(err)
	r.NotEmpty(variants)

	var variant product.Variant

	err = postgresql.SqlxTx(ts).Get(&variant,
		`select  weight, unit 
		 from product_variants
		 where product_id = $1
		 and weight = $2 
		 and unit = $3`,
		id, varquery.Weight, varquery.Unit)

	r.NoError(err)
	r.NotEmpty(variant)
	r.Equal(varquery.Weight, variant.Weight)
	r.Equal(varquery.Unit, variant.Unit)
}

func TestFindCurrentPrice(t *testing.T) {
	r := require.New(t)

	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()
	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	ts := sm.CreateSession()
	ts.Start()
	defer ts.Rollback()

	pp := product.ProductPrice{
		VariantID: 1,
		Price:     14.99,
		StartDate: time.Now(),
	}

	var id int
	err := postgresql.SqlxTx(ts).QueryRow(
		`insert into product_prices
		 ( variant_id, price, start_date )
	 	 values( $1, $2, $3 ) 
		 returning variant_id`,
		pp.VariantID, pp.Price, pp.StartDate).Scan(&id)
	r.NoError(err)

	price, err := repo.Repository.Product.FindCurrentPrice(ts, id)
	r.NoError(err)
	r.NotEmpty(price)
}

func TestInStorages(t *testing.T) {
	r := require.New(t)

	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()
	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	ts := sm.CreateSession()
	ts.Start()
	defer ts.Rollback()

	var id int
	product := stock.AddProductInStock{
		VariantID: 5,
		StorageID: 1,
		AddedAt:   time.Now(),
		Quantity:  3,
	}

	err := postgresql.SqlxTx(ts).QueryRow(
		`insert into products_in_storage
		 ( variant_id, storage_id, added_at, quantity )
		 values( $1, $2, $3, $4 )
		 returning variant_id`,
		product.VariantID, product.StorageID, product.AddedAt, product.Quantity).Scan(&id)
	r.NoError(err)

	inStorages, err := repo.Repository.Product.InStorages(ts, id)
	r.NoError(err)
	r.NotEmpty(inStorages)
}

func TestFindProductListByTag(t *testing.T) {
	r := require.New(t)

	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()
	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	ts := sm.CreateSession()
	ts.Start()
	defer ts.Rollback()

	p1 := product.Product{
		Name:    "sldlsd",
		Descr:   "sdsdsdsds",
		AddetAt: time.Now(),
		Tags:    "напиток",
	}

	p2 := product.Product{
		Name:    "авывывы",
		Descr:   "sdsdsdsds",
		AddetAt: time.Now(),
		Tags:    "стирка",
	}

	_, err := postgresql.SqlxTx(ts).Exec(`
	insert into products
	( name, description, added_at, tags )
	values( $1, $2, $3, $4 )`,
		p1.Name, p1.Descr, p1.AddetAt, p1.Tags)
	r.NoError(err)

	_, err = postgresql.SqlxTx(ts).Exec(`
	insert into products
	(name, description, added_at, tags )
	values($1, $2, $3, $4 )`,
		p2.Name, p2.Descr, p2.AddetAt, p2.Tags)
	r.NoError(err)

	tag := "напиток"
	limit := 3

	products, err := repo.Repository.Product.FindProductListByTag(ts, tag, limit)
	r.NoError(err)
	r.NotEmpty(products)

	tag = "стирка"
	limit = 1

	products, err = repo.Repository.Product.FindProductListByTag(ts, tag, limit)
	r.NoError(err)
	r.NotEmpty(products)
}

func TestFindStockVariantList(t *testing.T) {
	r := require.New(t)

	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()
	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	ts := sm.CreateSession()
	ts.Start()
	defer ts.Rollback()

	storageId := 1
	variants, err := repo.Repository.Product.FindStocksVariantList(ts, storageId)
	r.NoError(err)
	r.NotEmpty(variants)

	storageId = 2
	variants, err = repo.Repository.Product.FindStocksVariantList(ts, storageId)
	r.NoError(err)
	r.NotEmpty(variants)
}

func TestFindPrice(t *testing.T) {
	r := require.New(t)

	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()
	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	ts := sm.CreateSession()
	ts.Start()
	defer ts.Rollback()

	var variantID int
	priceQuery := product.ProductPrice{
		VariantID: 1,
		Price:     9.99,
		StartDate: time.Now(),
	}

	err := postgresql.SqlxTx(ts).QueryRow(`
	insert into product_prices
	( variant_id, price, start_date )
	values ( $1, $2, $3 ) 
	returning variant_id`,
		priceQuery.VariantID, priceQuery.Price, priceQuery.StartDate).Scan(&variantID)
	r.NoError(err)

	price, err := repo.Repository.Product.FindPrice(ts, variantID)
	r.NoError(err)
	r.NotEmpty(price)
}

func TestBuy(t *testing.T) {
	r := require.New(t)

	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()
	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	ts := sm.CreateSession()
	ts.Start()
	defer ts.Rollback()

	saleQuery := product.Sale{
		VariantID:  1,
		StorageID:  3,
		Quantity:   2,
		TotalPrice: 19.99,
	}

	saleID, err := repo.Repository.Product.Buy(ts, saleQuery)
	r.NoError(err)

	var sale product.Sale
	err = postgresql.SqlxTx(ts).Get(&sale, `
	select variant_id 
	from sales 
	where variant_id = $1 
	and storage_id = $2 
	and quantity = $3`,
		saleQuery.VariantID, saleQuery.StorageID, saleQuery.Quantity)
	r.NoError(err)
	r.NotZero(saleID)
	r.NotEmpty(sale)
}
func TestFindSaleList(t *testing.T) {
	r := require.New(t)

	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()
	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	ts := sm.CreateSession()
	ts.Start()
	defer ts.Rollback()

	startDate, err := time.Parse("02.01.2006", "01.07.2023")
	r.NoError(err)
	endDate, err := time.Parse("02.01.2006", "20.07.2023")
	r.NoError(err)

	sales, err := repo.Repository.Product.FindSaleListOnlyBySoldDate(ts, params.SaleQueryOnlyBySoldDate{
		Limit:     sqlnull.NewInt64(3),
		StartDate: startDate,
		EndDate:   endDate,
	})
	r.NoError(err)
	r.NotEmpty(sales)
}

func TestFindSaleListByFilters(t *testing.T) {
	r := require.New(t)

	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()
	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	ts := sm.CreateSession()
	ts.Start()
	defer ts.Rollback()

	startDate, err := time.Parse("02.01.2006", "01.07.2023")
	r.NoError(err)

	data, err := repo.Repository.Product.FindSaleListByFilters(ts, params.SaleQuery{
		StartDate: startDate,
		EndDate:   startDate.AddDate(0, 1, 0),
	})
	r.NoError(err)
	r.NotEmpty(data)

	data2, err := repo.Repository.Product.FindSaleListByFilters(ts, params.SaleQuery{
		StartDate: startDate,
		EndDate:   startDate.AddDate(0, 1, 0),
		StorageID: sqlnull.NewInt64(1),
	})

	r.NoError(err)
	r.NotEmpty(data2)

	data3, err := repo.Repository.Product.FindSaleListByFilters(ts, params.SaleQuery{
		StartDate:   startDate,
		EndDate:     startDate.AddDate(0, 1, 0),
		StorageID:   sqlnull.NewInt64(1),
		ProductName: sqlnull.NewString("Вода Hydrolife"),
	})

	r.NoError(err)
	r.NotEmpty(data3)
}
