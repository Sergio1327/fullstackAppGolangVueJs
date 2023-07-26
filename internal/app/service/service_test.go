package service_test

import (
	"go-back/internal/app/repository"
	"go-back/internal/app/service"
	"go-back/internal/pkg/database"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindProductListByTag(t *testing.T) {
	r := require.New(t)
	tag := "молоко"
	limit := 3

	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	r.NoError(err)
	defer db.Close()

	tx, err := db.Beginx()
	defer tx.Rollback()
	r.NoError(err)

	repo := repository.NewPostgresProductRepository(db)
	serv := service.NewProductUseCase(repo)

	products, err := serv.FindProductList(tag, limit)
	r.NoError(err)
	r.NoError(err)
	r.NoError(err)
	r.NotEmpty(products)
}
