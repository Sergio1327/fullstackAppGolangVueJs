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
	AddProduct(p *domain.Product) error
	AddProductPrice(pr *domain.ProductPrice) error
	AddProductInStock(p *domain.AddProductInStock) error
	GetProductInfoById(id int) (domain.ProductInfo, error)
	GetProductList(tag string, limit int) ([]domain.ProductInfo, error)
	GetProductsInStock(productId int) ([]domain.Stock, error)
	Buy(p *domain.Sale) error
	GetSales(sq *domain.SaleQuery) ([]domain.Sale, error)
}

type ProductServiceImpl struct {
	repo repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) *ProductServiceImpl {
	return &ProductServiceImpl{
		repo: repo,
	}
}
//логика добавление продукта в базу
func (u *ProductServiceImpl) AddProduct(p *domain.Product) error {
	tx, err := u.repo.TxBegin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if p.Name == "" {
		return errors.New("product_name cannot be empty")
	}
	productId, err := u.repo.AddProduct(p)
	if err != nil {
		return err
	}
	for _, v := range p.Variants {
		err := u.repo.AddProductVariants(productId, &v)
		if err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}


//логика проверки цены и вставки в базу
func (u *ProductServiceImpl) AddProductPrice(p *domain.ProductPrice) error {
	tx, err := u.repo.TxBegin()
	defer tx.Rollback()
	if err != nil {
		return err
	}
	variantID := strconv.Itoa(p.VariantId)
	if variantID == "" {
		return errors.New("no product with this variant_id")
	}

	if p.Price.IsZero() {
		return errors.New("price cant be zero or empty")
	}

	if p.StartDate == (time.Time{}) {
		return errors.New("date cant be empty")
	}

	isExistsId, err := u.repo.CheckExists(p)
	if err != nil {
		return err
	}
	if p.EndDate != (time.Time{}) {
		if isExistsId > 0 {
			p.EndDate = time.Now()
			err := u.repo.UpdateProductPrice(p, isExistsId)
			if err != nil {
				return err
			}
			tx.Commit()
		} else {
			err := u.repo.AddProductPriceWithEndDate(p)
			if err != nil {
				return err
			}
			tx.Commit()
		}
	} else {
		err := u.repo.AddProductPrice(p)
		if err != nil {
			return err
		}
	}
	return nil
}
//Логика проверка продукта на складе и обновления или добавления на базу
func (u *ProductServiceImpl) AddProductInStock(p *domain.AddProductInStock) error {

	tx, err := u.repo.TxBegin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if p.VariantId == 0 || p.StorageId == 0 || p.Quantity == 0 || p.Added_at == (time.Time{}) {
		return errors.New("variant_id,storage_id,quantity or added_at is empty")
	}
	isExist, err := u.repo.CheckProductsInStock(p)
	if err != nil {
		return err
	}
	if isExist {
		err := u.repo.UpdateProductsInstock(p)
		if err != nil {
			return err
		}
		tx.Commit()
	} else {
		err := u.repo.AddProductInStock(p)
		if err != nil {
			return err
		}
		tx.Commit()
	}

	return nil
}
//Логика получения всей информации о продукте и его вариантах по id 
func (u *ProductServiceImpl) GetProductInfoById(id int) (domain.ProductInfo, error) {
	if id == 0 || id < 0 {
		return domain.ProductInfo{}, errors.New("id cannot be zero or less than 0")
	}

	var productInfo domain.ProductInfo
	productInfo.ProductId = id
	err := u.repo.GetProductInfo(id, &productInfo)
	if err != nil {
		return domain.ProductInfo{}, err
	}
	err = u.repo.GetProductVariants(id, productInfo.Variants, &productInfo)
	if err != nil {
		return domain.ProductInfo{}, nil
	}
	for i := range productInfo.Variants {
		err := u.repo.GetCurrentPrice(&productInfo.Variants[i])
		if err != nil {
			return domain.ProductInfo{}, nil
		}
		err = u.repo.InStorages(&productInfo.Variants[i])
		if err != nil {
			return domain.ProductInfo{}, err
		}
	}
	return productInfo, nil
}


//Логика получения списка продуктов по тегу и лимиту
func (u *ProductServiceImpl) GetProductList(tag string, limit int) ([]domain.ProductInfo, error) {
	if limit == 0 || limit < 0 {
		limit = 3
	}
	if tag != "" {
		products, err := u.repo.GetProductsByTag(tag, limit)
		if err != nil {
			return nil, err
		}
		for i := range products {
			err := u.repo.GetProductVariants(products[i].ProductId, products[i].Variants, &products[i])
			if err != nil {
				return nil, err
			}
			variants := products[i].Variants
			for j := range variants {
				err := u.repo.GetCurrentPrice(&variants[j])
				if err != nil {
					return nil, err
				}
				err = u.repo.InStorages(&variants[j])
				if err != nil {
					return nil, err
				}
			}
		}
		return products, nil
	} else {
		products, err := u.repo.GetProducts(limit)
		if err != nil {
			return nil, err
		}
		for i := range products {
			err := u.repo.GetProductVariants(products[i].ProductId, products[i].Variants, &products[i])
			if err != nil {
				return nil, err
			}
			variants := products[i].Variants
			for j := range variants {
				err := u.repo.GetCurrentPrice(&variants[j])
				if err != nil {
					return nil, err
				}
				err = u.repo.InStorages(&variants[j])
				if err != nil {
					return nil, err
				}
			}
		}
		return products, nil
	}
}


//Логика получения всех складов и продуктов в ней или фильтрация по продукту
func (u *ProductServiceImpl) GetProductsInStock(productId int) ([]domain.Stock, error) {
	if productId < 0 {
		return nil, errors.New("product_id cannot be less than 0")
	}
	if productId == 0 {
		stocks, err := u.repo.GetStocks()
		if err != nil {
			return nil, err
		}
		for i, v := range stocks {
			variants, err := u.repo.GetStocksVariants(v.StorageID)
			if err != nil {
				return nil, err
			}
			stocks[i].ProductVariants = variants
		}

		return stocks, nil
	} else {
		stocks, err := u.repo.GetStocksByProductId(productId)
		if err != nil {
			return nil, err
		}
		for i, v := range stocks {
			variants, err := u.repo.GetStocksVariants(v.StorageID)
			if err != nil {
				return nil, err
			}
			stocks[i].ProductVariants = variants
		}

		return stocks, nil

	}
}
//Лоигка записи о покупке в базу
func (u *ProductServiceImpl) Buy(p *domain.Sale) error {
	tx, err := u.repo.TxBegin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if p.VariantId == 0 || p.StorageId == 0 || p.Quantity == 0 {
		return errors.New("variant_id,storage_id or quantity is empy")
	}
	p.SoldAt = time.Now()
	price, err := u.repo.GetPrice(p.VariantId)
	if err != nil {
		return err
	}
	p.TotalPrice = price.Mul(decimal.NewFromInt(int64(p.Quantity)))
	err = u.repo.Buy(p)
	if err != nil {
		return err
	}
	return nil
}


//Получение списка всех продаж или списка продаж по фильтрам	
func (u *ProductServiceImpl) GetSales(sq *domain.SaleQuery) ([]domain.Sale, error) {
	if sq.Limit == 0 {
		sq.Limit = 3
	}
	if sq.ProductName == "" && sq.StorageId == 0 {
		sales, err := u.repo.GetSales(sq)
		if err != nil {
			return nil, err
		}
		return sales, nil
	} else {
		sales, err := u.repo.GetSalesByFilters(sq)
		if err != nil {
			return nil, err
		}
		return sales, nil
	}

}
