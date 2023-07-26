package repository_test

import (
	"errors"
	"go-back/internal/app/domain"
	"go-back/internal/app/repository"
	"go-back/internal/pkg/database"
	"go-back/internal/tools/sqlnull"
	"testing"
	"time"

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
	defer tx.Rollback()

	p := domain.Product{
		Name:    "dsfdsddcxcfzsd",
		Descr:   "sdsdsds",
		AddetAt: time.Now(),
		Tags:    "123,12i",
	}

	id, err := repo.AddProduct(tx, p)
	r.NoError(err)
	r.NotEmpty(id)

	product, err := repo.LoadProductInfo(tx, id)
	r.NoError(err)
	r.Equal(id, product.ProductID)
	r.NotEmpty(product)
}

func TestAddProductVariantList(t *testing.T) {
	r := require.New(t)

	id := 1

	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)

	tx, err := repo.TxBegin()
	r.NoError(err)
	defer tx.Rollback()

	varQuery := domain.Variant{
		Weight: 440,
		Unit:   "г",
	}

	err = repo.AddProductVariantList(tx, id, varQuery)
	r.NoError(err)

	var variant domain.Variant
	err = tx.Get(&variant,
		`select variant_id 
		 from product_variants 
		 where product_id=$1 
		 and weight=$2 
		 and unit=$3`,
		id, varQuery.Weight, varQuery.Unit)

	r.NoError(err)
	r.NotEmpty(variant)
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
		VariantID: 4,
		StartDate: time.Now(),
		Price:     15.2,
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

func TestAddProductPrice(t *testing.T) {
	r := require.New(t)

	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	tx, err := db.Beginx()
	r.NoError(err)
	defer tx.Rollback()

	repo := repository.NewPostgresProductRepository(db)

	startDate, err := time.Parse("02.01.2006", "01.07.2023")
	r.NoError(err)

	priceID, err := repo.AddProductPrice(tx, domain.ProductPrice{
		VariantID: 5,
		StartDate: startDate,
		Price:     18.99,
	})

	r.NotZero(priceID)
	r.NoError(err)
}

func TestCheckProductInStock(t *testing.T) {
	r := require.New(t)

	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)
	tx, err := repo.TxBegin()
	r.NoError(err)

	_, err = tx.Exec(
		`insert into products_in_storage
		 (variant_id,storage_id,quantity) 
		 values($1,$2,$3)`, 4, 2, 10)
	r.NoError(err)

	_, err = repo.CheckProductInStock(tx, domain.AddProductInStock{
		VariantID: 4,
		StorageID: 2,
	})
	r.NoError(err)
}

func TestUpdateProductInStock(t *testing.T) {
	r := require.New(t)

	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	tx, err := db.Beginx()
	r.NoError(err)
	defer tx.Rollback()

	_, err = tx.Exec(
		`insert into products_in_storage
		 (variant_id,storage_id,quantity)
		 values ($1,$2,$3)`, 4, 2, 2)
	r.NoError(err)

	repo := repository.NewPostgresProductRepository(db)
	productStockID, err := repo.UpdateProductInstock(tx, domain.AddProductInStock{
		VariantID: 4,
		StorageID: 2,
		Quantity:  3,
	})
	r.NotZero(productStockID)
	r.NoError(err)
}
func TestAddProductInStock(t *testing.T) {
	r := require.New(t)

	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	tx, err := db.Beginx()
	r.NoError(err)
	defer tx.Rollback()

	repo := repository.NewPostgresProductRepository(db)
	productStockID, err := repo.AddProductInStock(tx, domain.AddProductInStock{
		VariantID: 3,
		StorageID: 1,
		AddedAt:   time.Now(),
		Quantity:  5,
	})

	r.NotZero(productStockID)
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
	defer tx.Rollback()

	id, err := repo.AddProduct(tx, domain.Product{
		Name:    "Имя продукта",
		Descr:   "Описание продукта",
		AddetAt: time.Now(),
		Tags:    "tag",
	})
	r.NoError(err)

	productInfo, err := repo.LoadProductInfo(tx, id)
	r.NoError(err)
	r.NotEmpty(productInfo)
}

func TestFindProductVariantList(t *testing.T) {
	r := require.New(t)

	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)

	tx, err := repo.TxBegin()
	r.NoError(err)
	defer tx.Rollback()

	varquery := domain.Variant{
		Weight: 500,
		Unit:   "г",
	}

	id := 2
	err = repo.AddProductVariantList(tx, id, varquery)
	r.NoError(err)

	variants, err := repo.FindProductVariantList(tx, id)
	r.NoError(err)
	r.NotEmpty(variants)

	var variant domain.Variant

	err = tx.Get(&variant,
		`select variant_id 
		 from product_variants
		 where product_id=$1
		 and weight=$2 
		 and unit = $3`,
		id, varquery.Weight, varquery.Unit)
	r.NoError(err)
	r.NotEmpty(variant)
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
	defer tx.Rollback()

	var id int
	pp := domain.ProductPrice{
		VariantID: 1,
		Price:     14.99,
		StartDate: time.Now(),
	}

	err = tx.QueryRow(
		`insert into product_prices
		 (variant_id,price,start_date)
	 	 values($1,$2,$3)
		 returning variant_id`,
		pp.VariantID, pp.Price, pp.StartDate).Scan(&id)
	r.NoError(err)

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
	defer tx.Rollback()

	var id int
	product := domain.AddProductInStock{
		VariantID: 5,
		StorageID: 1,
		AddedAt:   time.Now(),
		Quantity:  3,
	}

	err = tx.QueryRow(
		`insert into products_in_storage
		 (variant_id,storage_id,added_at,quantity)
		 values($1,$2,$3,$4)
		 returning variant_id`,
		product.VariantID, product.StorageID, product.AddedAt, product.Quantity).Scan(&id)
	r.NoError(err)

	inStorages, err := repo.InStorages(tx, id)
	r.NoError(err)
	r.NotEmpty(inStorages)
}

func TestFindProductListByTag(t *testing.T) {
	r := require.New(t)

	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)

	tx, err := repo.TxBegin()
	r.NoError(err)
	defer tx.Rollback()

	product := domain.Product{
		Name:    "sldlsd",
		Descr:   "sdsdsdsds",
		AddetAt: time.Now(),
		Tags:    "напиток",
	}

	product2 := domain.Product{
		Name:    "авывывы",
		Descr:   "sdsdsdsds",
		AddetAt: time.Now(),
		Tags:    "стирка",
	}

	_, err = tx.Exec(`
	insert into products
	(name,description,added_at,tags)
	values($1,$2,$3,$4)`,
		product.Name, product.Descr, product.AddetAt, product.Tags)
	r.NoError(err)

	_, err = tx.Exec(`
	insert into products
	(name,description,added_at,tags)
	values($1,$2,$3,$4)`,
		product2.Name, product2.Descr, product2.AddetAt, product2.Tags)
	r.NoError(err)

	tag := "напиток"
	limit := 3

	products, err := repo.FindProductListByTag(tx, tag, limit)
	r.NoError(err)
	r.NotEmpty(products)

	tag = "стирка"
	limit = 1

	products, err = repo.FindProductListByTag(tx, tag, limit)
	r.NoError(err)
	r.NotEmpty(products)
}

func TestLoadProductList(t *testing.T) {
	r := require.New(t)

	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)

	tx, err := repo.TxBegin()
	r.NoError(err)
	defer tx.Rollback()

	limit := 3

	products, err := repo.LoadProductList(tx, limit)
	r.NoError(err)
	r.NotEmpty(products)

	limit = 1

	products, err = repo.LoadProductList(tx, limit)
	r.NoError(err)
	r.NotEmpty(products)

	if len(products) > 1 {
		r.Error(errors.New("кол-во продуктов больше чем указано в лимите"))
	}
}

func TestLoadStockList(t *testing.T) {
	r := require.New(t)

	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)

	tx, err := repo.TxBegin()
	r.NoError(err)
	defer tx.Rollback()

	stocks, err := repo.LoadStockList(tx)
	r.NoError(err)
	r.NotEmpty(stocks)
}

func TestFindStockListByProductId(t *testing.T) {
	r := require.New(t)

	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)

	tx, err := repo.TxBegin()
	r.NoError(err)
	defer tx.Rollback()

	id := 1
	stocks, err := repo.FindStockListByProductId(tx, id)
	r.NoError(err)
	r.NotEmpty(stocks)

	id = 2
	stocks, err = repo.FindStockListByProductId(tx, id)
	r.NoError(err)
	r.NotEmpty(stocks)
}

func TestFindStockVariantList(t *testing.T) {
	r := require.New(t)

	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)

	tx, err := repo.TxBegin()
	r.NoError(err)
	defer tx.Rollback()

	storageId := 1
	variants, err := repo.FindStocksVariantList(tx, storageId)
	r.NoError(err)
	r.NotEmpty(variants)

	storageId = 2
	variants, err = repo.FindStocksVariantList(tx, storageId)
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
	defer tx.Rollback()

	var variantID int
	priceQuery := domain.ProductPrice{
		VariantID: 1,
		Price:     9.99,
		StartDate: time.Now(),
	}

	err = tx.QueryRow(`
	insert into product_prices
	(variant_id,price,start_date)
	values ($1,$2,$3) 
	returning variant_id`,
		priceQuery.VariantID, priceQuery.Price, priceQuery.StartDate).Scan(&variantID)
	r.NoError(err)

	price, err := repo.FindPrice(tx, variantID)
	r.NoError(err)
	r.NotEmpty(price)
}

func TestBuy(t *testing.T) {
	r := require.New(t)

	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)

	tx, err := repo.TxBegin()
	r.NoError(err)
	defer tx.Rollback()

	saleQuery := domain.Sale{
		VariantID:  1,
		StorageID:  3,
		Quantity:   2,
		TotalPrice: 19.99,
	}

	saleID, err := repo.Buy(tx, saleQuery)
	r.NoError(err)

	var sale domain.Sale
	err = tx.Get(&sale, `
	select variant_id 
	from sales 
	where variant_id=$1 
	and storage_id=$2 
	and quantity=$3`,
		saleQuery.VariantID, saleQuery.StorageID, saleQuery.Quantity)
	r.NoError(err)
	r.NotZero(saleID)
	r.NotEmpty(sale)
}

func TestFindSaleList(t *testing.T) {
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
	defer tx.Rollback()

	sales, err := repo.FindSaleListOnlyBySoldDate(tx, domain.SaleQueryOnlyBySoldDate{
		Limit:     sqlnull.NewInt64(3),
		StartDate: startDate,
		EndDate:   endDate,
	})
	r.NoError(err)
	r.NotEmpty(sales)
}

func TestFindSaleListByFilters(t *testing.T) {
	r := require.New(t)

	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	//инициализация слоев
	repo := repository.NewPostgresProductRepository(db)

	tx, err := repo.TxBegin()
	r.NoError(err)
	defer tx.Rollback()

	startDate, err := time.Parse("02.01.2006", "01.07.2023")
	r.NoError(err)

	data, err := repo.FindSaleListByFilters(tx, domain.SaleQuery{
		StartDate: startDate,
		EndDate:   startDate.AddDate(0, 1, 0),
	})
	r.NoError(err)
	r.NotEmpty(data)

	data2, err := repo.FindSaleListByFilters(tx, domain.SaleQuery{
		StartDate: startDate,
		EndDate:   startDate.AddDate(0, 1, 0),
		StorageId: sqlnull.NewNullInt64(1),
	})

	r.NoError(err)
	r.NotEmpty(data2)

	data3, err := repo.FindSaleListByFilters(tx, domain.SaleQuery{
		StartDate:   startDate,
		EndDate:     startDate.AddDate(0, 1, 0),
		StorageId:   sqlnull.NewNullInt64(1),
		ProductName: sqlnull.NewNullString("Вода Hydrolife"),
	})

	r.NoError(err)
	r.NotEmpty(data3)
}
