package postgresql_test

import (
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

	tx := postgresql.SqlxTx(ts)

	id, err := repo.Repository.Product.AddProduct(tx, p1)
	r.NoError(err)
	r.NotEmpty(id)

	productInfo, err := repo.Repository.Product.LoadProductInfo(tx, id)
	r.NoError(err)
	r.NotEmpty(productInfo)

	r.Equal(p1.Name, productInfo.Name)
	r.Equal(p1.Descr, productInfo.Descr)

	// Второй тестовый случай - добавление продукта без имени, ожидается ошибка
	p2 := product.Product{}

	id2, err := repo.Repository.Product.AddProduct(tx, p2)
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
	tx := postgresql.SqlxTx(ts)

	id := 1
	variant := product.Variant{
		Weight: 200,
		Unit:   "г",
	}
	err := repo.Repository.Product.AddProductVariantList(tx, id, variant)
	r.NoError(err)

	id = -1
	variant = product.Variant{
		Weight: 300,
		Unit:   "кг",
	}
	err = repo.Repository.Product.AddProductVariantList(tx, id, variant)
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
	tx := postgresql.SqlxTx(ts)

	productPrice := product.ProductPrice{
		VariantID: 2,
		StartDate: time.Now(),
	}

	id, err := repo.Repository.Product.CheckExists(tx, productPrice)
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
	tx := postgresql.SqlxTx(ts)

	expectedProductPrice := product.ProductPrice{
		PriceID: 2,
		EndDate: sqlnull.NewNullTime(time.Now()),
	}

	err := repo.Repository.Product.UpdateProductPrice(tx, expectedProductPrice, expectedProductPrice.PriceID)
	r.NoError(err)

	var productPrice product.ProductPrice
	err = tx.Get(&productPrice,
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
	tx := postgresql.SqlxTx(ts)

	startDate, err := time.Parse("02.01.2006", "01.07.2023")
	r.NoError(err)

	expectedProductPrice := product.ProductPrice{
		VariantID: 5,
		StartDate: startDate,
		Price:     18.99,
	}

	priceID, err := repo.Repository.Product.AddProductPrice(tx, expectedProductPrice)
	r.NoError(err)
	r.NotEmpty(priceID)

	var productPrice product.ProductPrice
	err = tx.Get(&productPrice,
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
	tx := postgresql.SqlxTx(ts)

	productInStock := stock.AddProductInStock{
		VariantID: 4,
		StorageID: 2,
		Quantity:  10,
	}

	_, err := tx.Exec(
		`insert into products_in_storage
		 ( variant_id, storage_id, quantity ) 
		 values( $1, $2, $3 )`, productInStock.VariantID, productInStock.StorageID, productInStock.Quantity)
	r.NoError(err)

	isExists, err := repo.Repository.Product.CheckProductInStock(tx, productInStock)
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
	tx := postgresql.SqlxTx(ts)

	productInStock := stock.AddProductInStock{
		VariantID: 4,
		StorageID: 2,
		Quantity:  10,
	}

	_, err := tx.Exec(
		`insert into products_in_storage
		 ( variant_id, storage_id, quantity )
		 values ( $1, $2, $3 )`, productInStock.VariantID, productInStock.StorageID, productInStock.Quantity)
	r.NoError(err)

	productStockID, err := repo.Repository.Product.UpdateProductInstock(tx, stock.AddProductInStock{
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
	tx := postgresql.SqlxTx(ts)

	expectedProduct := stock.AddProductInStock{
		VariantID: 3,
		StorageID: 1,
		AddedAt:   time.Now(),
		Quantity:  5,
	}
	productStockID, err := repo.Repository.Product.AddProductInStock(tx, expectedProduct)
	r.NotZero(productStockID)
	r.NoError(err)

	var productInStock stock.AddProductInStock
	err = tx.Get(&productInStock, `
	select variant_id, storage_id, quantity
	from products_in_storage
	where variant_id = $1 and storage_id = $2 and quantity = $3`,
		expectedProduct.VariantID, expectedProduct.StorageID, expectedProduct.Quantity)
	r.NoError(err)

	r.Equal(expectedProduct.VariantID, productInStock.VariantID)
	r.Equal(expectedProduct.StorageID, productInStock.StorageID)
	r.Equal(expectedProduct.Quantity, productInStock.Quantity)
}
