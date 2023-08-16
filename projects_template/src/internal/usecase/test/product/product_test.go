package test

import (
	"product_storage/internal/entity/global"
	"product_storage/internal/entity/product"
	"product_storage/internal/transaction"
	"product_storage/rimport"
	"product_storage/tools/logger"
	"product_storage/uimport"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	testLogger = logger.NewNoFileLogger("test")
)

func TestSaveSale(t *testing.T) {
	r := require.New(t)

	type fields struct {
		ts *transaction.MockSession
		ri rimport.TestRepositoryImports
	}

	type args struct {
		sale product.SaleParams
	}

	fixedTime := time.Date(2023, 8, 9, 13, 52, 40, 0, time.UTC)

	argSale := product.SaleParams{
		VariantID: 1,
		StorageID: 1,
		SoldAt:    fixedTime,
		Quantity:  2,
	}

	tests := []struct {
		name               string
		prepare            func(f *fields)
		expectedID         int
		expectedTotalPrice float64
		args               args
		err                error
	}{
		{
			name: "успешный результат",
			prepare: func(f *fields) {
				price := 5.99
				saleID := 123
				sale := product.SaleParams{
					VariantID:  1,
					StorageID:  1,
					Quantity:   2,
					SoldAt:     fixedTime,
					TotalPrice: price * float64(argSale.Quantity),
				}
				f.ri.MockRepository.Product.EXPECT().FindPrice(f.ts, sale.VariantID).Return(price, nil)
				f.ri.MockRepository.Product.EXPECT().SaveSale(f.ts, sale).Return(saleID, nil)
			},
			args: args{
				sale: argSale,
			},
			expectedID:         123,
			expectedTotalPrice: 5.99,
			err:                nil,
		},
		{
			name: "безуспешный результат",
			prepare: func(f *fields) {
				price := 0.0 // Assuming price is 0.0 in this case
				f.ri.MockRepository.Product.EXPECT().FindPrice(f.ts, argSale.VariantID).Return(price, global.ErrNoData)
			},
			args: args{
				argSale,
			},
			expectedTotalPrice: 0,
			expectedID:         0,
			err:                global.ErrInternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := fields{
				ri: rimport.NewTestRepositoryImports(ctrl),
				ts: transaction.NewMockSession(ctrl),
			}

			if tt.prepare != nil {
				tt.prepare(&f)
			}

			sm := transaction.NewMockSessionManager(ctrl)
			ui := uimport.NewUsecaseImports(testLogger, testLogger, f.ri.RepositoryImports(), sm)

			data, err := ui.Usecase.Product.SaveSale(f.ts, tt.args.sale)

			r.Equal(tt.err, err)
			r.Equal(tt.expectedID, data)
		})
	}
}
