package repository_test

import (
	"errors"
	"go-back/internal/app/domain"
	"go-back/internal/app/repository"
	"go-back/internal/pkg/database"
	"go-back/internal/tools/sqlnull"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestAddProduct(t *testing.T) {
	r := require.New(t)

	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()
	repo := repository.NewPostgresProductRepository(db)

	tx, err := repo.TxBegin()
	r.NoError(err)
	id, err := repo.AddProduct(tx, domain.Product{
		Name:     "dsfdsddcxcfzsd",
		Descr:    "sdsdsds",
		Addet_at: time.Now(),
		Tags:     "123,12i",
	})
	r.NoError(err)
	r.NotEmpty(id)
}

func TestAddProductVariants(t *testing.T) {
	r := require.New(t)
	id := 1
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)

	tx, err := repo.TxBegin()
	defer tx.Rollback()
	r.NoError(err)

	err = repo.AddProductVariants(tx, id, domain.Variant{
		Weight: 440,
		Unit:   "г",
	})

	r.NoError(err)
}

func TestCheckExists(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()
	repo := repository.NewPostgresProductRepository(db)

	tx, err := repo.TxBegin()
	r.NoError(err)

	_, err = repo.CheckExists(tx, domain.ProductPrice{
		VariantId: 4,
		StartDate: time.Now(),
		Price:     decimal.New(15, 99),
	})
	r.NoError(err)
}

func TestUpdateProductPrice(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	tx, err := db.Beginx()
	defer tx.Rollback()
	r.NoError(err)
	repo := repository.NewPostgresProductRepository(db)
	err = repo.UpdateProductPrice(tx, domain.ProductPrice{
		EndDate: sqlnull.NewNullTime(time.Now()),
	}, 2)
	r.NoError(err)
}

func TestAddProductPriceWithEndDate(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	tx, err := db.Beginx()
	defer tx.Rollback()
	r.NoError(err)
	repo := repository.NewPostgresProductRepository(db)
	startDate, err := time.Parse("02.01.2006", "01.07.2023")
	r.NoError(err)
	endDate, err := time.Parse("02.01.2006", "20.07.2023")
	r.NoError(err)
	err = repo.AddProductPriceWithEndDate(tx, domain.ProductPrice{
		VariantId: 4,
		StartDate: startDate,
		EndDate:   sqlnull.NewNullTime(endDate),
		Price:     decimal.New(15, 2),
	})
	r.NoError(err)
}
func TestAddProductPrice(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	tx, err := db.Beginx()
	r.NoError(err)
	repo := repository.NewPostgresProductRepository(db)
	startDate, err := time.Parse("02.01.2006", "01.07.2023")
	r.NoError(err)

	err = repo.AddProductPrice(tx, domain.ProductPrice{
		VariantId: 5,
		StartDate: startDate,
		Price:     decimal.New(10, 2),
	})
	r.NoError(err)
}

func TestCheckProductsInStock(t *testing.T) {
	r := require.New(t)

	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)
	tx, err := repo.TxBegin()
	r.NoError(err)

	_, err = tx.Exec("insert into products_in_storage(variant_id,storage_id,quantity) values($1,$2,$3)", 4, 2, 10)
	r.NoError(err)

	_, err = repo.CheckProductsInStock(tx, domain.AddProductInStock{
		VariantId: 4,
		StorageId: 2,
	})
	r.NoError(err)
}

func TestUpdateProductsInStock(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()
	tx, err := db.Beginx()
	r.NoError(err)

	_, err = tx.Exec("insert into products_in_storage(variant_id,storage_id,quantity) values ($1,$2,$3)", 4, 2, 2)
	r.NoError(err)

	repo := repository.NewPostgresProductRepository(db)
	err = repo.UpdateProductsInstock(tx, domain.AddProductInStock{
		VariantId: 4,
		StorageId: 2,
		Quantity:  3,
	})
	r.NoError(err)
}
func TestAddProductInStock(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()
	tx, err := db.Beginx()
	defer tx.Rollback()
	r.NoError(err)

	repo := repository.NewPostgresProductRepository(db)
	err = repo.AddProductInStock(tx, domain.AddProductInStock{
		VariantId: 3,
		StorageId: 1,
		Quantity:  5,
	})
	r.NoError(err)
}

func TestLoadProductInfo(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()
	repo := repository.NewPostgresProductRepository(db)

	tx, err := repo.TxBegin()
	r.NoError(err)

	id, err := repo.AddProduct(tx, domain.Product{
		Name:     "Имя продукта",
		Descr:    "Описание продукта",
		Addet_at: time.Now(),
		Tags:     "tag",
	})
	r.NoError(err)

	productInfo, err := repo.LoadProductInfo(tx, id)
	r.NoError(err)
	r.NotEmpty(productInfo)
}

func TestFindProductVariants(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)
	tx, err := repo.TxBegin()
	r.NoError(err)

	id := 2
	err = repo.AddProductVariants(tx, id, domain.Variant{
		Weight: 500,
		Unit:   "г",
	})
	r.NoError(err)

	variants, err := repo.FindProductVariants(tx, id)
	r.NoError(err)
	r.NotEmpty(variants)

}

func TestFindCurrentPrice(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)
	tx, err := repo.TxBegin()
	r.NoError(err)

	id := 1
	price, err := repo.FindCurrentPrice(tx, id)
	r.NoError(err)
	r.NotEmpty(price)

}

func TestInStorages(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)
	tx, err := repo.TxBegin()
	r.NoError(err)

	id := 1
	inStorages, err := repo.InStorages(tx, id)
	r.NoError(err)
	r.NotEmpty(inStorages)
}

func TestFindProductsByTag(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)
	tx, err := repo.TxBegin()
	r.NoError(err)

	tag := "напиток"
	limit := 3

	products, err := repo.FindProductsByTag(tx, tag, limit)
	r.NoError(err)
	r.NotEmpty(products)

	tag = "стирка"
	limit = 1

	products, err = repo.FindProductsByTag(tx, tag, limit)
	r.NoError(err)
	r.NotEmpty(products)
}

func TestLoadProducts(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)
	tx, err := repo.TxBegin()
	r.NoError(err)
	limit := 3

	products, err := repo.LoadProducts(tx, limit)
	r.NoError(err)
	r.NotEmpty(products)

	limit = 1

	products, err = repo.LoadProducts(tx, limit)
	r.NoError(err)
	r.NotEmpty(products)
	if len(products) > 1 {
		r.Error(errors.New("кол-во продуктов больше чем указано в лимите"))
	}
}

func TestLoadStocks(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)
	tx, err := repo.TxBegin()
	r.NoError(err)

	stocks, err := repo.LoadStocks(tx)
	r.NoError(err)
	r.NotEmpty(stocks)
}

func TestFindStocksByProductId(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)
	tx, err := repo.TxBegin()
	r.NoError(err)

	id := 1
	stocks, err := repo.FindStocksByProductId(tx, id)
	r.NoError(err)
	r.NotEmpty(stocks)

	id = 2
	stocks, err = repo.FindStocksByProductId(tx, id)
	r.NoError(err)
	r.NotEmpty(stocks)
}

func TestFindStockVariants(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()
	repo := repository.NewPostgresProductRepository(db)
	tx, err := repo.TxBegin()
	r.NoError(err)
	storageId := 1

	variants, err := repo.FindStocksVariants(tx, storageId)
	r.NoError(err)
	r.NotEmpty(variants)

	storageId = 2
	variants, err = repo.FindStocksVariants(tx, storageId)
	r.NoError(err)
	r.NotEmpty(variants)
}

func TestFindPrice(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()
	repo := repository.NewPostgresProductRepository(db)

	tx, err := repo.TxBegin()
	r.NoError(err)

	variantId := 2
	price, err := repo.FindPrice(tx, variantId)
	r.NoError(err)
	r.NotEmpty(price)

	variantId = 3
	price, err = repo.FindPrice(tx, variantId)
	r.NoError(err)
	r.NotEmpty(price)
}

func TestBuy(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()
	tx, err := db.Beginx()
	r.NoError(err)
	repo := repository.NewPostgresProductRepository(db)
	err = repo.Buy(tx, domain.Sale{
		VariantId:  1,
		StorageId:  3,
		Quantity:   2,
		TotalPrice: decimal.New(10, 2),
	})
	r.NoError(err)
}

func TestFindSales(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	startDate, err := time.Parse("02.01.2006", "01.07.2023")
	r.NoError(err)
	endDate, err := time.Parse("02.01.2006", "20.07.2023")
	r.NoError(err)

	repo := repository.NewPostgresProductRepository(db)
	tx, err := repo.TxBegin()
	r.NoError(err)
	sales, err := repo.FindSales(tx, domain.SaleQueryWithoutFilters{
		Limit:     sqlnull.NewInt64(3),
		StartDate: startDate,
		EndDate:   endDate,
	})
	r.NoError(err)
	r.NotEmpty(sales)
}

func TestFindSalesByFilters(t *testing.T) {
	r := require.New(t)

	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	//инициализация слоев
	repo := repository.NewPostgresProductRepository(db)
	tx, err := repo.TxBegin()
	r.NoError(err)

	startDate, err := time.Parse("02.01.2006", "01.07.2023")
	r.NoError(err)

	data, err := repo.FindSalesByFilters(tx, domain.SaleQuery{
		StartDate: startDate,
		EndDate:   startDate.AddDate(0, 1, 0),
	})
	r.NoError(err)
	r.NotEmpty(data)

	data2, err := repo.FindSalesByFilters(tx, domain.SaleQuery{
		StartDate: startDate,
		EndDate:   startDate.AddDate(0, 1, 0),
		StorageId: sqlnull.NewNullInt64(1),
	})

	r.NoError(err)
	r.NotEmpty(data2)

	data3, err := repo.FindSalesByFilters(tx, domain.SaleQuery{
		StartDate:   startDate,
		EndDate:     startDate.AddDate(0, 1, 0),
		StorageId:   sqlnull.NewNullInt64(1),
		ProductName: sqlnull.NewNullString("Вода Hydrolife"),
	})

	r.NoError(err)
	r.NotEmpty(data3)
}
