package service

import (
	"errors"
	"go-back/internal/app/domain"
	"go-back/internal/app/repository"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

type ProductService interface {
	AddProduct(p domain.Product) error
	AddProductPrice(pr domain.ProductPrice) error
	AddProductInStock(p domain.AddProductInStock) error
	FindProductInfoById(id int) (domain.ProductInfo, error)
	FindProductList(tag string, limit int) ([]domain.ProductInfo, error)
	FindProductsInStock(productId int) ([]domain.Stock, error)
	Buy(p domain.Sale) error
	LoadSales(sq domain.SaleQuery) ([]domain.Sale, error)
}

type ProductServiceImpl struct {
	repo repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) *ProductServiceImpl {
	return &ProductServiceImpl{
		repo: repo,
	}
}

// AddProduct логика добавление продукта в базу
func (u *ProductServiceImpl) AddProduct(p domain.Product) error {
	tx, err := u.repo.TxBegin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if p.Name == "" {
		return errors.New("имя продукта не может быть пустым")
	}
	productId, err := u.repo.AddProduct(tx, p)
	if err != nil {
		return err
	}
	for _, v := range p.Variants {
		err := u.repo.AddProductVariants(tx, productId, v)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// AddProductPrice  логика проверки цены и вставки в базу
func (u *ProductServiceImpl) AddProductPrice(p domain.ProductPrice) error {

	tx, err := u.repo.TxBegin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	variantID := strconv.Itoa(p.VariantId)
	if variantID == "" {
		return errors.New("нет варианта продукта с таким id")
	}

	if p.Price.IsZero() {
		return errors.New("цена не может быть пустой или равна 0")
	}

	if p.StartDate == (time.Time{}) {
		return errors.New("дата не может быть пустой")
	}
	//Проверка имеется ли запись уже в базе с заданным id продукта и дата начала цены
	isExistsId, err := u.repo.CheckExists(p)
	if err != nil {
		return err
	}
	//Если пользователь ввел дату окончания цены то
	//прооисходит проверка есть ли записи уже в базе
	if p.EndDate.Valid {
		//Если записи есть то вставляется дата окончания цены
		if isExistsId > 0 {
			p.EndDate.Time = time.Now()
			err := u.repo.UpdateProductPrice(tx, p, isExistsId)
			if err != nil {
				return err
			}

		} else {
			//Если же нет то просто добавляется запись в базу
			err := u.repo.AddProductPriceWithEndDate(tx, p)
			if err != nil {
				return err
			}

		}
		//Если пользователь не ввел дату окончания то просто вставляется новая запись в базу
	} else {
		err := u.repo.AddProductPrice(tx, p)
		if err != nil {
			return err
		}

	}
	err = tx.Commit()
	return err
}

// AddProductInStock Логика проверка продукта на складе и обновления или добавления на базу
func (u *ProductServiceImpl) AddProductInStock(p domain.AddProductInStock) error {

	tx, err := u.repo.TxBegin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = p.IsNullFields()
	if err != nil {
		return err
	}
	isExist, err := u.repo.CheckProductsInStock(p)
	if err != nil {
		return err
	}
	if isExist {
		err := u.repo.UpdateProductsInstock(tx, p)
		if err != nil {
			return err
		}

	} else {
		err := u.repo.AddProductInStock(tx, p)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	return err
}

// FindProductInfoById  Логика получения всей информации о продукте и его вариантах по id
func (u *ProductServiceImpl) FindProductInfoById(id int) (domain.ProductInfo, error) {
	if id == 0 || id < 0 {
		return domain.ProductInfo{}, errors.New("id не может быть меньше или равен 0")
	}

	product, err := u.repo.LoadProductInfo(id)
	if err != nil {
		return domain.ProductInfo{}, err
	}
	product.ProductId = id

	product.Variants, err = u.repo.FindProductVariants(product.ProductId)
	if err != nil {
		return domain.ProductInfo{}, err
	}
	for i, v := range product.Variants {
		price, err := u.repo.FindCurrentPrice(v.VariantId)
		if err != nil {
			return domain.ProductInfo{}, err
		}
		product.Variants[i].ProductId = id
		product.Variants[i].CurrentPrice = price

		inStorages, err := u.repo.InStorages(v.VariantId)
		if err != nil {
			return domain.ProductInfo{}, err
		}
		product.Variants[i].InStorages = inStorages
	}
	return product, nil
}

// LoadProductList  Логика получения списка продуктов по тегу и лимиту
func (u *ProductServiceImpl) FindProductList(tag string, limit int) ([]domain.ProductInfo, error) {
	if limit == 0 || limit < 0 {
		limit = 3
	}
	if tag != "" {
		products, err := u.repo.FindProductsByTag(tag, limit)
		if err != nil {
			return nil, err
		}
		for i := range products {
			vars, err := u.repo.FindProductVariants(products[i].ProductId)
			if err != nil {
				return nil, err
			}
			products[i].Variants = vars
			variants := products[i].Variants
			for j := range variants {
				price, err := u.repo.FindCurrentPrice(variants[j].VariantId)
				if err != nil {
					return nil, err
				}
				variants[j].ProductId = products[i].ProductId
				variants[j].CurrentPrice = price
				inStorages, err := u.repo.InStorages(variants[j].VariantId)
				if err != nil {
					return nil, err
				}
				variants[j].InStorages = inStorages
			}
		}
		return products, nil
	} else {
		products, err := u.repo.LoadProducts(limit)
		if err != nil {
			return nil, err
		}
		for i := range products {
			vars, err := u.repo.FindProductVariants(products[i].ProductId)
			if err != nil {
				return nil, err
			}
			products[i].Variants = vars
			variants := products[i].Variants
			for j := range variants {
				price, err := u.repo.FindCurrentPrice(variants[j].VariantId)
				if err != nil {
					return nil, err
				}
				variants[j].ProductId = products[i].ProductId
				variants[j].CurrentPrice = price
				inStorages, err := u.repo.InStorages(variants[j].VariantId)
				if err != nil {
					return nil, err
				}
				variants[j].InStorages = inStorages
			}
		}
		return products, nil
	}
}

// FindProductsInStock  Логика получения всех складов и продуктов в ней или фильтрация по продукту
func (u *ProductServiceImpl) FindProductsInStock(productId int) ([]domain.Stock, error) {
	if productId < 0 {
		return nil, errors.New("id продукта не может быть меньше нуля")
	}
	if productId == 0 {
		stocks, err := u.repo.LoadStocks()
		if err != nil {
			return nil, err
		}
		for i, v := range stocks {
			variants, err := u.repo.FindStocksVariants(v.StorageID)
			if err != nil {
				return nil, err
			}
			stocks[i].ProductVariants = variants
		}

		return stocks, nil
	} else {
		stocks, err := u.repo.FindStocksByProductId(productId)
		if err != nil {
			return nil, err
		}
		for i, v := range stocks {
			variants, err := u.repo.FindStocksVariants(v.StorageID)
			if err != nil {
				return nil, err
			}
			stocks[i].ProductVariants = variants
		}

		return stocks, nil

	}
}

// Buy  Логuка записи о покупке в базу
func (u *ProductServiceImpl) Buy(p domain.Sale) error {
	tx, err := u.repo.TxBegin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = p.IsNullFields()
	if err != nil {
		return err
	}
	p.SoldAt = time.Now()
	price, err := u.repo.FindPrice(p.VariantId)
	if err != nil {
		return err
	}
	p.TotalPrice = price.Mul(decimal.NewFromInt(int64(p.Quantity)))
	err = u.repo.Buy(tx, p)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

// LoadSales  Получение списка всех продаж или списка продаж по фильтрам
func (u *ProductServiceImpl) LoadSales(sq domain.SaleQuery) ([]domain.Sale, error) {
	if !sq.Limit.Valid {
		sq.Limit.Int64 = 3
	}
	if !sq.ProductName.Valid && !sq.StorageId.Valid {
		sales, err := u.repo.LoadSales(sq)
		if err != nil {
			return nil, err
		}
		return sales, nil
	} else {
		sales, err := u.repo.FindSalesByFilters(sq)
		if err != nil {
			return nil, err
		}
		return sales, nil
	}

}
