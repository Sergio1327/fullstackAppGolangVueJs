package service

import (
	"errors"
	"go-back/internal/app/domain"
	"go-back/internal/app/repository"
	"strconv"
	"time"
)

type ProductService interface {
	AddProduct(p domain.Product) (int, error)
	AddProductPrice(pr domain.ProductPrice) error
	AddProductInStock(p domain.AddProductInStock) error
	FindProductInfoById(id int) (domain.ProductInfo, error)
	FindProductList(tag string, limit int) ([]domain.ProductInfo, error)
	FindProductsInStock(productId int) ([]domain.Stock, error)
	Buy(p domain.Sale) error
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
func (u *ProductServiceImpl) AddProduct(product domain.Product) (productId int, err error) {
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
	productId, err = u.repo.AddProduct(tx, product)
	if err != nil {
		return 0, errors.New("не удалось добавить продукт в базу данных")
	}

	// если пользователь не ввел варианты продукта то данные о продукте просто запишутся в базу
	if product.Variants == nil {
		err = tx.Commit()
		if err != nil {
			return 0, errors.New("ошибка в добавлении вариантов продукта")
		}
		return productId, nil
	}

	// добавляются варианты продукта
	for _, v := range product.Variants {
		err := u.repo.AddProductVariants(tx, productId, v)
		if err != nil {
			return 0, errors.New("не удалось добавить варианты продукта")
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, errors.New("не удалось добавить продукт и его варианты в базу")
	}

	return productId, nil
}

// AddProductPrice  логика проверки цены и вставки в базу
func (u *ProductServiceImpl) AddProductPrice(p domain.ProductPrice) error {
	tx, err := u.repo.TxBegin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	variantID := strconv.Itoa(p.VariantId)

	//проверка  id варианта, цены, даты начала цены на нулевые значения
	if variantID == "" {
		return errors.New("нет варианта продукта с таким id")
	}

	if p.Price == 0 {
		return errors.New("цена не может быть пустой или равна 0")
	}

	if p.StartDate == (time.Time{}) {
		return errors.New("дата не может быть пустой")
	}
	// проверка имеется ли запись уже в базе с заданным id продукта и дата начала цены
	isExistsId, err := u.repo.CheckExists(tx, p)
	if err != nil {
		return errors.New("ошибка при проверке цен в базе данных")
	}
	// если пользователь ввел дату окончания цены то
	// прооисходит проверка есть ли записи уже в базе
	if p.EndDate.Valid {
		// если записи есть то вставляется дата окончания цены
		if isExistsId > 0 {
			p.EndDate.Time = time.Now()
			err := u.repo.UpdateProductPrice(tx, p, isExistsId)
			if err != nil {
				return errors.New("не удалось обновить цену")
			}
		} else {
			// если же нет то просто добавляется запись в базу
			err := u.repo.AddProductPriceWithEndDate(tx, p)
			if err != nil {
				return errors.New("не удалось добавить цену")
			}

		}
		// если пользователь не ввел дату окончания то просто вставляется новая запись в базу
	} else {
		err := u.repo.AddProductPrice(tx, p)
		if err != nil {
			return errors.New("не удалось добавить цену")
		}

	}
	err = tx.Commit()
	return err
}

// AddProductInStock логика проверка продукта на складе и обновления или добавления на базу
func (u *ProductServiceImpl) AddProductInStock(p domain.AddProductInStock) error {
	tx, err := u.repo.TxBegin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// проверка запроса на нулевые значения
	err = p.IsNullFields()
	if err != nil {
		return err
	}

	// проверка есть ли уже продукт на складе
	isExist, err := u.repo.CheckProductsInStock(tx, p)
	if err != nil {
		return errors.New("ошибка при проверке наличия продукта на складе")
	}

	// если продукт уже имеется в базе обновляется его кол-во
	if isExist {
		err := u.repo.UpdateProductsInstock(tx, p)
		if err != nil {
			return errors.New("не удалось обновить кол-во продуктов на складе")
		}
	} else {
		// если продукта нет на складе то он просто добавляется на склад
		err := u.repo.AddProductInStock(tx, p)
		if err != nil {
			return errors.New("не удалось добавить продукт на склад")
		}
	}

	err = tx.Commit()
	return err
}

// FindProductInfoById  логика получения всей информации о продукте и его вариантах по id
func (u *ProductServiceImpl) FindProductInfoById(id int) (product domain.ProductInfo, err error) {
	tx, err := u.repo.TxBegin()
	if err != nil {
		return domain.ProductInfo{}, nil
	}
	defer tx.Rollback()

	// если пользователь не ввел id выводится ошибка
	if id == 0 || id < 0 {
		return domain.ProductInfo{}, errors.New("id не может быть меньше или равен 0")
	}

	// поиск продукта по его id
	product, err = u.repo.LoadProductInfo(tx, id)
	if err != nil {
		return domain.ProductInfo{}, errors.New("не удалось получить информацию о продукте")
	}

	// проверка имеется ли на данный моммент у продукта его варианты
	isExists, err := u.repo.AreExistsVariants(tx, id)
	if err != nil {
		return domain.ProductInfo{}, errors.New("ошибка при проверке наличия вариантов продукта")
	}

	// если вариантов нет  то выведется информация о продукте без вариантов
	if !isExists {
		emptyVariants := []domain.Variant{}
		product.Variants = emptyVariants
		return product, nil
	}

	// если варианты есть то происходит поиск всех вариантов
	product.Variants, err = u.repo.FindProductVariants(tx, product.ProductId)
	if err != nil {
		return domain.ProductInfo{}, errors.New("не удалось найти варианты продукта")
	}

	for i, v := range product.Variants {
		price, err := u.repo.FindCurrentPrice(tx, v.VariantId)
		if err != nil {
			return domain.ProductInfo{}, errors.New("не удалось найти актуальную цену варианта продукта")
		}

		// получение актуальной цены для каждого варианта продукта
		product.Variants[i].CurrentPrice = price

		// получение id складов в которых есть этот продукт
		inStorages, err := u.repo.InStorages(tx, v.VariantId)
		if err != nil {
			return domain.ProductInfo{}, errors.New("не удалось найти склады в которых есть продукт")
		}
		product.Variants[i].InStorages = inStorages
	}

	return product, nil
}

// LoadProductList  логика получения списка продуктов по тегу и лимиту
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
		products, err = u.repo.FindProductsByTag(tx, tag, limit)
		if err != nil {
			return nil, errors.New("не удалось найти продукты по данному тегу")
		}

		for i := range products {
			vars, err := u.repo.FindProductVariants(tx, products[i].ProductId)
			if err != nil {
				return nil, errors.New("не удалось найти варианты продукта")
			}
			products[i].Variants = vars
			variants := products[i].Variants
			for j := range variants {
				price, err := u.repo.FindCurrentPrice(tx, variants[j].VariantId)
				if err != nil {
					return nil, errors.New("не удалось найти актуальную цену продукта")
				}

				variants[j].CurrentPrice = price
				inStorages, err := u.repo.InStorages(tx, variants[j].VariantId)
				if err != nil {
					return nil, errors.New("не удалось найти склады в которых есть продукт")
				}

				variants[j].InStorages = inStorages
			}
		}

		return products, nil
	} else {
		// если пользователь не ввел тег то просто прозойдет поиск всех продуктов с лимитом вывода
		products, err = u.repo.LoadProducts(tx, limit)
		if err != nil {
			return nil, err
		}

		for i := range products {
			vars, err := u.repo.FindProductVariants(tx, products[i].ProductId)
			if err != nil {
				return nil, err
			}

			products[i].Variants = vars
			variants := products[i].Variants
			for j := range variants {
				price, err := u.repo.FindCurrentPrice(tx, variants[j].VariantId)
				if err != nil {
					return nil, err
				}

				variants[j].CurrentPrice = price
				inStorages, err := u.repo.InStorages(tx, variants[j].VariantId)
				if err != nil {
					return nil, err
				}

				variants[j].InStorages = inStorages
			}
		}

		return products, nil
	}
}

// FindProductsInStock  логика получения всех складов и продуктов в ней или фильтрация по продукту
func (u *ProductServiceImpl) FindProductsInStock(productId int) (stocks []domain.Stock, err error) {
	tx, err := u.repo.TxBegin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if productId < 0 {
		return nil, errors.New("id продукта не может быть меньше нуля")
	}

	// если пользователь не ввел id продукта то будет выполнен поиск всех складов
	if productId == 0 {
		stocks, err = u.repo.LoadStocks(tx)
		if err != nil {
			return nil, errors.New("не удалось найти склады")
		}
		for i, v := range stocks {
			variants, err := u.repo.FindStocksVariants(tx, v.StorageID)
			if err != nil {
				return nil, errors.New("не удалось найти варианты продукта на складе")
			}
			stocks[i].ProductVariants = variants
		}

		return stocks, nil
	} else {

		//Если же пользователь ввел id продукта то произойдет фильтрация складов по id продукта
		stocks, err = u.repo.FindStocksByProductId(tx, productId)
		if err != nil {
			return nil, errors.New("не удалось найти склады с продуктами по данному id ")
		}
		for i, v := range stocks {
			variants, err := u.repo.FindStocksVariants(tx, v.StorageID)
			if err != nil {
				return nil, errors.New("не удалось найти варианты продукта на складе")
			}
			stocks[i].ProductVariants = variants
		}

		return stocks, nil

	}
}

// Buy  логuка записи о покупке в базу
func (u *ProductServiceImpl) Buy(p domain.Sale) error {
	tx, err := u.repo.TxBegin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// проверка фильтров на нулевые значения ,которые ввел пользователь
	err = p.IsNullFields()
	if err != nil {
		return err
	}

	// устанавливаем текущую дату как дату продажи
	p.SoldAt = time.Now()

	// получение цены варианта
	price, err := u.repo.FindPrice(tx, p.VariantId)
	if err != nil {
		return errors.New("не удалось найти цену продукта")
	}

	// подсчет общей цены продажи
	p.TotalPrice = price * float64(p.Quantity)

	// запись продажи в базу
	err = u.repo.Buy(tx, p)
	if err != nil {
		return errors.New("не удалось записать продажу в базу")
	}
	err = tx.Commit()
	return err
}

// LoadSales  получение списка всех продаж или списка продаж по фильтрам
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
		s := domain.SaleQueryWithoutFilters{
			StartDate: sq.StartDate,
			EndDate:   sq.EndDate,
			Limit:     sq.Limit,
		}

		sales, err = u.repo.FindSales(tx, s)
		if err != nil {
			return nil, errors.New("не удалось найти продажи")
		}

		return sales, nil
	} else {
		//  если имя продукта или id склада указан то произойдет фильтрация по этим параметрам
		sales, err = u.repo.FindSalesByFilters(tx, sq)
		if err != nil {
			return nil, errors.New("не удалось найти продажи по данным фильтрам")
		}

		return sales, nil
	}
}
