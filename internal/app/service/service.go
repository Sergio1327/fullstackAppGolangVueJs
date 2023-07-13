package service

import (
	"errors"
	"go-back/internal/app/domain"
	"go-back/internal/app/repository"
	"strconv"
	"time"
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
	} else {
		tx.Commit()
	}
	for _, v := range p.Variants {
		err := u.repo.AddProductVariants(productId, &v)
		if err != nil {
			return err
		}
	}
	return nil
}

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
		}
	} else {
		err := u.repo.AddProductPrice(p)
		if err != nil {
			return err
		}
	}
	return nil
}

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
	} else {
		err := u.repo.AddProductInStock(p)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *ProductServiceImpl) GetProductInfoById(id int) (domain.ProductInfo, error) {
	if id == 0 || id < 0 {
		return domain.ProductInfo{}, errors.New("id cannot be zero or less than 0")
	}

	var productInfo domain.ProductInfo
	err := u.repo.GetProductInfo(id, &productInfo)
	if err != nil {
		return domain.ProductInfo{}, err
	}
	err = u.repo.GetProductVariants(id, productInfo.Variants,&productInfo)
	if err != nil {
		return domain.ProductInfo{}, nil
	}
	for _, v := range productInfo.Variants {
		err := u.repo.GetCurrentPrice(&v)
		if err != nil {
			return domain.ProductInfo{}, nil
		}
		err = u.repo.InStorages(&v)
		if err != nil {
			return domain.ProductInfo{}, err
		}
	}
	return productInfo, nil
}

func (u *ProductServiceImpl) GetProductList(tag string, limit int) ([]domain.ProductInfo, error) {
	if limit == 0 || limit < 0 {
		limit = 3
	}
	if tag != "" {
		products, err := u.repo.GetProductListByTag(tag, limit)
		if err != nil {
			return nil, err
		}
		return products, nil
	} else {
		products, err := u.repo.GetProductList(limit)
		if err != nil {
			return nil, err
		}
		return products, nil
	}
}

func (u *ProductServiceImpl) GetProductsInStock(productId int) ([]domain.Stock, error) {
	if productId < 0 {
		return nil, errors.New("product_id cannot be less than 0")
	}

	if productId != 0 {
		stocks, err := u.repo.GetProductsInStockById(productId)
		if err != nil {
			return nil, err
		}
		return stocks, nil
	} else {
		stocks, err := u.repo.GetProductsInStock()
		if err != nil {
			return nil, err
		}
		return stocks, nil
	}
}

func (u *ProductServiceImpl) Buy(p *domain.Sale) error {
	tx, err := u.repo.TxBegin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if p.VariantId == 0 || p.StorageId == 0 || p.Quantity == 0 {
		return errors.New("variant_id,storage_id pr quantity is empy")
	}
	err = u.repo.Buy(p)
	if err != nil {
		return err
	} else {
		tx.Commit()
	}
	return nil
}

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
