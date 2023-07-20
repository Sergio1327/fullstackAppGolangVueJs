package repository_test

import (
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

	tx, err := db.Beginx()
	r.NoError(err)
	id, err := repo.AddProduct(tx, domain.Product{
		Name:     "dsfdsddcxcfz",
		Descr:    "sdsdsds",
		Addet_at: time.Now(),
		Tags:     "123,12i",
	})
	r.NoError(err)
	err = tx.Commit()
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

	tx, err := db.Beginx()
	r.NoError(err)
	repo := repository.NewPostgresProductRepository(db)
	err = repo.AddProductVariants(tx, id, domain.Variant{
		Weight: 440,
		Unit:   "г",
	})

	r.NoError(err)
	err = tx.Commit()
	r.NoError(err)
}

func TestCheckExists(t *testing.T) {
	r := require.New(t)
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()
	repo := repository.NewPostgresProductRepository(db)

	_, err = repo.CheckExists(domain.ProductPrice{
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
		Price:     decimal.New(15, 99),
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
		Price:     decimal.New(20, 15),
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
	isExists, err := repo.CheckProductsInStock(domain.AddProductInStock{
		VariantId: 4,
		StorageId: 2,
	})
	r.NoError(err)
	r.NotEmpty(isExists)
}

func TestFindSalesByFilters(t *testing.T) {
	r := require.New(t)

	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	//инициализация слоев
	repo := repository.NewPostgresProductRepository(db)

	startDate, err := time.Parse("02.01.2006", "01.07.2023")
	r.NoError(err)

	data, err := repo.FindSalesByFilters(domain.SaleQuery{
		StartDate: startDate,
		EndDate:   startDate.AddDate(0, 1, 0),
	})
	r.NoError(err)
	r.NotEmpty(data)
}
