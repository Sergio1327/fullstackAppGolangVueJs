package product

import (
	"errors"
	"projects_template/internal/entity/product"
	"projects_template/internal/transaction"
	"projects_template/rimport"

	"github.com/sirupsen/logrus"
)

type ProductUseCaseImpl struct {
	log   *logrus.Logger
	dbLog *logrus.Logger
	rimport.RepositoryImports
}

func NewProductUseCaseImpl(log, dblog *logrus.Logger, ri rimport.RepositoryImports) ProductUseCaseImpl {
	return ProductUseCaseImpl{
		log:               log,
		dbLog:             dblog,
		RepositoryImports: ri,
	}
}

func (u *ProductUseCaseImpl) AddProduct(ts transaction.Session, product product.Product) (productID int, err error) {
	// если имя продукта не введено то возвращается ошибка
	err = ts.Start()
	defer ts.Rollback()
	if product.Name == "" {
		return 0, errors.New("имя продукта не может быть пустым")
	}

	// добавляется продукт в базу
	productID, err = u.Repository.Product.AddProduct(ts, product)
	if err != nil {
		return 0, errors.New("не удалось добавить продукт в базу данных")
	}
	
	// если пользователь не ввел варианты продукта то данные о продукте просто запишутся в базу
	if product.VariantList == nil {
		err = ts.Commit()
		if err != nil {
			return 0, errors.New("ошибка в добавлении вариантов продукта")
		}
		return productID, nil
	}

	// добавляются варианты продукта
	for _, v := range product.VariantList {
		err := u.Repository.Product.AddProductVariantList(ts, productID, v)
		if err != nil {
			return 0, errors.New("не удалось добавить варианты продукта")
		}
	}

	err = ts.Commit()
	return productID, err
}
