package repository_test

import (
	"go-back/internal/app/domain"
	"go-back/internal/app/repository"
	"go-back/internal/pkg/database"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

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
