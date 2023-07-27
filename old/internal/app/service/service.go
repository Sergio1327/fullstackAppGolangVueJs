package service

import (
	"database/sql"
	"errors"
	"go-back/internal/app/domain"
	"go-back/internal/app/repository"
	"log"
	"strconv"
	"time"
)

type ProductService interface {
	AddProduct(p domain.Product) (int, error)
	AddProductPrice(pr domain.ProductPrice) (int, error)
	AddProductInStock(p domain.AddProductInStock) (int, error)
	FindProductInfoById(productID int) (domain.ProductInfo, error)
	FindProductList(tag string, limit int) ([]domain.ProductInfo, error)
	FindProductsInStock(productID int) ([]domain.Stock, error)
	Buy(p domain.Sale) (int, error)
	FindSales(sq domain.SaleQuery) ([]domain.Sale, error)
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
func (u *ProductServiceImpl) AddProduct(product domain.Product) (productID int, err error) {
	tx, err := u.repo.TxBegin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// если имя продукта не введено то возвращается ошибка
	if product.Name == "" {
		return 0, errors.New("имя продукта не может быть пустым")
	}

	// добавляется продукт в базу
	productID, err = u.repo.AddProduct(tx, product)
	if err != nil {
		return 0, errors.New("не удалось добавить продукт в базу данных")
	}

	// если пользователь не ввел варианты продукта то данные о продукте просто запишутся в базу
	if product.VariantList == nil {
		err = tx.Commit()
		if err != nil {
			return 0, errors.New("ошибка в добавлении вариантов продукта")
		}
		return productID, nil
	}

	// добавляются варианты продукта
	for _, v := range product.VariantList {
		err := u.repo.AddProductVariantList(tx, productID, v)
		if err != nil {
			return 0, errors.New("не удалось добавить варианты продукта")
		}
	}

	err = tx.Commit()
	return productID, err
}

// AddProductPrice логика проверки цены и вставки в базу
func (u *ProductServiceImpl) AddProductPrice(p domain.ProductPrice) (priceID int, err error) {
	tx, err := u.repo.TxBegin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	variantID := strconv.Itoa(p.VariantID)

	//проверка  id варианта, цены, даты начала цены на нулевые значения
	if variantID == "" {
		return 0, errors.New("нет варианта продукта с таким id")
	}

	if p.Price == 0 {
		return 0, errors.New("цена не может быть пустой или равна 0")
	}

	if p.StartDate == (time.Time{}) {
		return 0, errors.New("дата не может быть пустой")
	}

	// проверка имеется ли запись уже в базе с заданным id продукта и дата начала цены
	isExistsID, err := u.repo.CheckExists(tx, p)
	if err != nil {
		return 0, errors.New("ошибка при проверке цен в базе данных")
	}

	// если запись уже имеется устанавливается дата окончания цены
	if isExistsID > 0 {
		p.EndDate.Scan(time.Now())
		err := u.repo.UpdateProductPrice(tx, p, isExistsID)
		if err != nil {
			return 0, errors.New("не удалось обновить цену")
		}

		priceID = isExistsID
	} else {
		// если записи нет то цена вставляется в базу
		priceID, err = u.repo.AddProductPrice(tx, p)
		if err != nil {
			return 0, errors.New("не удалось добавить цену")
		}
	}

	err = tx.Commit()
	return priceID, err
}

// AddProductInStock логика проверка продукта на складе и обновления или добавления на базу
func (u *ProductServiceImpl) AddProductInStock(p domain.AddProductInStock) (productStockID int, err error) {
	tx, err := u.repo.TxBegin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// проверка запроса на нулевые значения
	err = p.IsNullFields()
	if err != nil {
		return 0, err
	}

	// проверка есть ли уже продукт на складе
	isExist, err := u.repo.CheckProductInStock(tx, p)
	if err != nil {
		return 0, errors.New("ошибка при проверке наличия продукта на складе")
	}

	// если продукт уже имеется в базе обновляется его кол-во
	if isExist {
		productStockID, err = u.repo.UpdateProductInstock(tx, p)
		if err != nil {
			return 0, errors.New("не удалось обновить кол-во продуктов на складе")
		}
	} else {
		// если продукта нет на складе то он просто добавляется на склад
		productStockID, err = u.repo.AddProductInStock(tx, p)
		if err != nil {
			return 0, errors.New("не удалось добавить продукт на склад")
		}
	}

	err = tx.Commit()
	return productStockID, err
}

// FindProductInfoById логика получения всей информации о продукте и его вариантах по id
func (u *ProductServiceImpl) FindProductInfoById(productID int) (product domain.ProductInfo, err error) {
	tx, err := u.repo.TxBegin()
	if err != nil {
		return domain.ProductInfo{}, nil
	}
	defer tx.Rollback()

	// если пользователь не ввел id выводится ошибка
	if productID <= 0 {
		return domain.ProductInfo{}, errors.New("id не может быть меньше или равен 0")
	}

	// поиск продукта по его id
	product, err = u.repo.LoadProductInfo(tx, productID)
	if err != nil {
		return domain.ProductInfo{}, errors.New("не удалось получить информацию о продукте")
	}

	product.VariantList, err = u.repo.FindProductVariantList(tx, product.ProductID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return product, nil
		default:
			return domain.ProductInfo{}, errors.New("не удалось найти варианты продукта")
		}

	}

	for i, v := range product.VariantList {
		price, err := u.repo.FindCurrentPrice(tx, v.VariantID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				continue
			default:
				return domain.ProductInfo{}, errors.New("не удалось найти актуальную цену варианта продукта")
			}
		}

		// получение актуальной цены для каждого варианта продукта
		product.VariantList[i].CurrentPrice = price

		// получение id складов в которых есть этот продукт
		inStorages, err := u.repo.InStorages(tx, v.VariantID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				continue
			default:
				return domain.ProductInfo{}, errors.New("не удалось найти склады в которых есть продукт")
			}
		}
		product.VariantList[i].InStorages = inStorages
	}
	err = tx.Commit()
	return product, err
}

// LoadProductList логика получения списка продуктов по тегу и лимиту
func (u *ProductServiceImpl) FindProductList(tag string, limit int) (products []domain.ProductInfo, err error) {
	tx, err := u.repo.TxBegin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// если лимит не указан или некорректен то по умолчанию устанавливается 3
	if limit == 0 || limit < 0 {
		limit = 3
	}

	// если пользователь ввел тег продукта произойдет поиск продуктов по данному тегу
	if tag != "" {
		products, err = u.repo.FindProductListByTag(tx, tag, limit)
		if err != nil {
			return nil, errors.New("не удалось найти продукты по данному тегу")
		}

		for i := range products {
			vars, err := u.repo.FindProductVariantList(tx, products[i].ProductID)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
					return products, nil
				default:
					return nil, errors.New("не удалось найти варианты продукта")
				}
			}
			products[i].VariantList = vars
			variantList := products[i].VariantList
			for j := range variantList {
				price, err := u.repo.FindCurrentPrice(tx, variantList[j].VariantID)
				if err != nil {
					switch err {
					case sql.ErrNoRows:
						continue
					default:
						return nil, errors.New("не удалось найти актуальную цену продукта")
					}
				}

				variantList[j].CurrentPrice = price
				inStorages, err := u.repo.InStorages(tx, variantList[j].VariantID)
				if err != nil {
					switch err {
					case sql.ErrNoRows:
						continue
					default:
						return nil, errors.New("не удалось найти склады в которых есть продукт")
					}
				}

				variantList[j].InStorages = inStorages
			}
		}

	} else {
		// если пользователь не ввел тег то просто прозойдет поиск всех продуктов с лимитом вывода
		products, err = u.repo.LoadProductList(tx, limit)
		if err != nil {
			return nil, err
		}

		for i := range products {
			vars, err := u.repo.FindProductVariantList(tx, products[i].ProductID)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
					return products, nil
				default:
					return nil, err
				}
			}

			products[i].VariantList = vars
			variants := products[i].VariantList
			for j := range variants {
				price, err := u.repo.FindCurrentPrice(tx, variants[j].VariantID)
				if err != nil {
					switch err {
					case sql.ErrNoRows:
						continue
					default:
						return nil, err
					}
				}

				variants[j].CurrentPrice = price
				inStorages, err := u.repo.InStorages(tx, variants[j].VariantID)
				if err != nil {
					switch err {
					case sql.ErrNoRows:
						continue
					default:
						return nil, err
					}
				}

				variants[j].InStorages = inStorages
			}
		}

	}

	err = tx.Commit()
	return products, err
}

// FindProductsInStock логика получения всех складов и продуктов в ней или фильтрация по продукту
func (u *ProductServiceImpl) FindProductsInStock(productID int) (stocks []domain.Stock, err error) {
	tx, err := u.repo.TxBegin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if productID < 0 {
		return nil, errors.New("id продукта не может быть меньше нуля")
	}

	// если пользователь не ввел id продукта то будет выполнен поиск всех складов
	if productID == 0 {
		stocks, err = u.repo.LoadStockList(tx)
		if err != nil {
			return nil, errors.New("не удалось найти склады")
		}
		for i, v := range stocks {
			variants, err := u.repo.FindStocksVariantList(tx, v.StorageID)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
					return stocks, nil
				default:
					return nil, errors.New("не удалось найти варианты продукта на складе")
				}
			}
			stocks[i].ProductVariantList = variants
		}
	} else {

		//Если же пользователь ввел id продукта то произойдет фильтрация складов по id продукта
		stocks, err = u.repo.FindStockListByProductId(tx, productID)
		if err != nil {
			return nil, errors.New("не удалось найти склады с продуктами по данному id ")
		}
		for i, v := range stocks {
			variants, err := u.repo.FindStocksVariantList(tx, v.StorageID)
			if err != nil {
				switch err {
				case sql.ErrNoRows:
					return stocks, nil
				default:
					return nil, errors.New("не удалось найти варианты продукта на складе")
				}
			}
			stocks[i].ProductVariantList = variants
		}
	}
	err = tx.Commit()
	return stocks, err
}

// Buy логuка записи о покупке в базу
func (u *ProductServiceImpl) Buy(p domain.Sale) (saleID int, err error) {
	tx, err := u.repo.TxBegin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// проверка фильтров на нулевые значения ,которые ввел пользователь
	err = p.IsNullFields()
	if err != nil {
		return 0, err
	}

	// устанавливаем текущую дату как дату продажи
	p.SoldAt = time.Now()

	// получение цены варианта
	price, err := u.repo.FindPrice(tx, p.VariantID)
	if err != nil {
		return 0, errors.New("не удалось найти цену продукта")
	}

	// подсчет общей цены продажи
	p.TotalPrice = price * float64(p.Quantity)

	// запись продажи в базу
	saleID, err = u.repo.Buy(tx, p)
	if err != nil {
		return 0, errors.New("не удалось записать продажу в базу")
	}

	err = tx.Commit()
	return saleID, err
}

// LoadSales получение списка всех продаж или списка продаж по фильтрам
func (u *ProductServiceImpl) FindSales(sq domain.SaleQuery) (sales []domain.Sale, err error) {
	tx, err := u.repo.TxBegin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// если лимит не указан то по умолчанию устанавливается 3
	if !sq.Limit.Valid {
		sq.Limit.Scan(3)
	}

	// если не указано имя продукта или id склада то произойдет фильтрация только по датам
	if !sq.ProductName.Valid && !sq.StorageId.Valid {
		s := domain.SaleQueryOnlyBySoldDate{
			StartDate: sq.StartDate,
			EndDate:   sq.EndDate,
			Limit:     sq.Limit,
		}

		sales, err = u.repo.FindSaleListOnlyBySoldDate(tx, s)
		if err != nil {
			return nil, errors.New("не удалось найти продажи")
		}
	} else {
		//  если имя продукта или id склада указан то произойдет фильтрация по этим параметрам
		sales, err = u.repo.FindSaleListByFilters(tx, sq)
		if err != nil {
			log.Println(err)
			return nil, errors.New("не удалось найти продажи по данным фильтрам")
		}
	}

	err = tx.Commit()
	return sales, err
}