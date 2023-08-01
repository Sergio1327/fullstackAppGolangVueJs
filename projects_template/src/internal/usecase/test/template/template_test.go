// package template_test
package template_test

// import (
// 	"product_storage/bimport"
// 	"product_storage/internal/entity/global"
// 	"product_storage/internal/entity/template"
// 	"product_storage/internal/transaction"
// 	"product_storage/rimport"
// 	"product_storage/tools/logger"
// 	"product_storage/uimport"

// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// )

// var (
// 	testLogger = logger.NewNoFileLogger("test")
// )

// func TestAwesomePublicMethod(t *testing.T) {
// 	r := assert.New(t)

// 	type fields struct {
// 		ri rimport.TestRepositoryImports
// 		bi *bimport.TestBridgeImports
// 		ts *transaction.MockSession
// 	}
// 	type args struct {
// 		id int
// 	}

// 	const (
// 		id = 1
// 	)

// 	tests := []struct {
// 		name         string
// 		prepare      func(f *fields)
// 		args         args
// 		expectedData template.TemplateObject
// 		err          error
// 	}{
// 		{
// 			name: "успешный результат",
// 			prepare: func(f *fields) {
// 				templateData := template.TemplateObject{}

// 				gomock.InOrder(
// 					f.ri.MockRepository.Template.EXPECT().FindTemplateObj(f.ts, id).Return(templateData, nil),
// 				)
// 			},
// 			args: args{id: id},
// 			err:  nil,
// 		},
// 		{
// 			name: "неуспешный результат",
// 			prepare: func(f *fields) {
// 				templateData := template.TemplateObject{}

// 				gomock.InOrder(
// 					f.ri.MockRepository.Template.EXPECT().FindTemplateObj(f.ts, id).Return(templateData, global.ErrNoData),
// 				)
// 			},
// 			args: args{id: id},
// 			err:  global.ErrInternalError,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctrl := gomock.NewController(t)
// 			defer ctrl.Finish()
// 			f := fields{
// 				ri: rimport.NewTestRepositoryImports(ctrl),
// 				ts: transaction.NewMockSession(ctrl),
// 			}
// 			if tt.prepare != nil {
// 				tt.prepare(&f)
// 			}

// 			sm := transaction.NewMockSessionManager(ctrl)
// 			ui := uimport.NewUsecaseImports(testLogger, testLogger, f.ri.RepositoryImports(), f.bi.BridgeImports(), sm)

// 			// data, err := ui.Usecase.ProdcutUsecase..AwesomePublicMethod(f.ts, tt.args.id)
// 			// r.Equal(tt.err, err)
// 			r.Equal(tt.expectedData, data)

// 		})
// 	}
// }
