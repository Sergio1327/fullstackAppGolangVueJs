package test

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"product_storage/internal/entity/global"
	"product_storage/internal/entity/product"
	"product_storage/internal/transaction"
	"product_storage/rimport"
	"product_storage/tools/logger"
	"product_storage/uimport"
	"testing"
	"time"
)

var (
	testLogger = logger.NewNoFileLogger("test")
)

func TestBuy(t *testing.T) {
	r := require.New(t)

	type fields struct {
		ts *transaction.MockSession
		ri rimport.TestRepositoryImports
	}

	type args struct {
		sale product.Sale
	}

	fixedTime := time.Date(2023, 8, 9, 13, 52, 40, 0, time.UTC)
	argSale := product.Sale{
		VariantID: 1,
		StorageID: 1,
		SoldAt:    fixedTime,
		Quantity:  10,
	}

	tests := []struct {
		name         string
		prepare      func(f *fields)
		expectedData int
		args         args
		err          error
	}{
		{
			name: "успешный результат",
			prepare: func(f *fields) {
				var price float64
				saleID := 123

				gomock.InOrder(
					f.ri.MockRepository.Product.EXPECT().FindPrice(f.ts, gomock.Eq(argSale.VariantID)).Return(price, nil),
					f.ri.MockRepository.Product.EXPECT().Buy(f.ts, gomock.Eq(argSale)).Return(saleID, nil),
				)
			},
			args: args{
				sale: argSale,
			},
			expectedData: 123,
			err:          nil,
		},
		{
			name: "безуспешный результат",
			prepare: func(f *fields) {
				var price float64

				gomock.InOrder(
					f.ri.MockRepository.Product.EXPECT().FindPrice(f.ts, argSale.VariantID).Return(price, global.ErrNoData),
				)
			},
			args: args{
				sale: argSale,
			},
			expectedData: 0,
			err:          global.ErrInternalError,
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

			data, err := ui.Usecase.Product.Buy(f.ts, tt.args.sale)

			r.Equal(tt.err, err)
			r.Equal(tt.expectedData, data)
		})
	}
}
