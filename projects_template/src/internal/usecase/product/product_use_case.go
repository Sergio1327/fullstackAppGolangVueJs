package product

import (
	"database/sql"
	"errors"
	"log"
	"product_storage/internal/entity/product"
	"product_storage/internal/entity/stock"
	"product_storage/rimport"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
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
// AddProduct логика добавление продукта в базу
func (u ProductUseCaseImpl) AddProduct(tx *sqlx.Tx, product product.Product) (productID int, err error) {
	// если имя продукта не введено то возвращается ошибка
	if product.Name == "" {
		return 0, errors.New("имя продукта не может быть пустым")
	}

	// добавляется продукт в базу
	productID, err = u.Repository.Product.AddProduct(tx, product)
	if err != nil {
		return 0, errors.New("не удалось добавить продукт в базу данных")
	}

	// если пользователь не ввел варианты продукта то данные о продукте просто запишутся в базу
	if product.VariantList == nil {
		if err != nil {
			return 0, errors.New("ошибка в добавлении вариантов продукта")
		}

		return productID, nil
	} else {
		// добавляются варианты продукта
		for _, v := range product.VariantList {
			err := u.Repository.Product.AddProductVariantList(tx, productID, v)
			if err != nil {
				return 0, errors.New("не удалось добавить варианты продукта")
			}
		}

		return productID, err
	}
}

// AddProductPrice логика проверки цены и вставки в базу
func (u ProductUseCaseImpl) AddProductPrice(tx *sqlx.Tx, p product.ProductPrice) (priceID int, err error) {
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
	isExistsID, err := u.Repository.Product.CheckExists(tx, p)
	if err != nil {
		return 0, errors.New("ошибка при проверке цен в базе данных")
	}

	// если запись уже имеется устанавливается дата окончания цены
	if isExistsID > 0 {
		p.EndDate.Scan(time.Now())
		err := u.Repository.Product.UpdateProductPrice(tx, p, isExistsID)
		if err != nil {
			return 0, errors.New("не удалось обновить цену")
		}

		priceID = isExistsID
	} else {
		// если записи нет то цена вставляется в базу
		priceID, err = u.Repository.Product.AddProductPrice(tx, p)
		if err != nil {
			return 0, errors.New("не удалось добавить цену")
		}
	}

	return priceID, err
}

// AddProductInStock логика проверка продукта на складе и обновления или добавления на базу
func (u ProductUseCaseImpl) AddProductInStock(tx *sqlx.Tx, p stock.AddProductInStock) (productStockID int, err error) {

	// проверка запроса на нулевые значения
	err = p.IsNullFields()
	if err != nil {
		return 0, err
	}

	// проверка есть ли уже продукт на складе
	isExist, err := u.Repository.Product.CheckProductInStock(tx, p)
	if err != nil {
		return 0, errors.New("ошибка при проверке наличия продукта на складе")
	}

	// если продукт уже имеется в базе обновляется его кол-во
	if isExist {
		productStockID, err = u.Repository.Product.UpdateProductInstock(tx, p)
		if err != nil {
			return 0, errors.New("не удалось обновить кол-во продуктов на складе")
		}
	} else {
		// если продукта нет на складе то он просто добавляется на склад
		productStockID, err = u.Repository.Product.AddProductInStock(tx, p)
		if err != nil {
			return 0, errors.New("не удалось добавить продукт на склад")
		}
	}

	return productStockID, err
}

// FindProductInfoById логика получения всей информации о продукте и его вариантах по id
func (u ProductUseCaseImpl) FindProductInfoById(tx *sqlx.Tx, productID int) (productInfo product.ProductInfo, err error) {

	// если пользователь не ввел id выводится ошибка
	if productID <= 0 {
		return product.ProductInfo{}, errors.New("id не может быть меньше или равен 0")
	}

	// поиск продукта по его id
	productInfo, err = u.Repository.Product.LoadProductInfo(tx, productID)
	if err != nil {
		return product.ProductInfo{}, errors.New("не удалось получить информацию о продукте")
	}

	productInfo.VariantList, err = u.Repository.Product.FindProductVariantList(tx, productInfo.ProductID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return productInfo, nil
		default:
			return product.ProductInfo{}, errors.New("не удалось найти варианты продукта")
		}

	}

	for i, v := range productInfo.VariantList {
		price, err := u.Repository.Product.FindCurrentPrice(tx, v.VariantID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				continue
			default:
				return product.ProductInfo{}, errors.New("не удалось найти актуальную цену варианта продукта")
			}
		}

		// получение актуальной цены для каждого варианта продукта
		productInfo.VariantList[i].CurrentPrice = price

		// получение id складов в которых есть этот продукт
		inStorages, err := u.Repository.Product.InStorages(tx, v.VariantID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				continue
			default:
				return product.ProductInfo{}, errors.New("не удалось найти склады в которых есть продукт")
			}
		}
		productInfo.VariantList[i].InStorages = inStorages
	}

	return productInfo, err
}

// FindProductList логика получения списка продуктов по тегу и лимиту
func (u ProductUseCaseImpl) FindProductList(tx *sqlx.Tx, tag string, limit int) (products []product.ProductInfo, err error) {

	// если лимит не указан или некорректен то по умолчанию устанавливается 3
	if limit == 0 || limit < 0 {
		limit = 3
	}

	// если пользователь ввел тег продукта произойдет поиск продуктов по данному тегу
	if tag != "" {
		products, err = u.Repository.Product.FindProductListByTag(tx, tag, limit)
		if err != nil {
			return nil, errors.New("не удалось найти продукты по данному тегу")
		}

		for i := range products {
			vars, err := u.Repository.Product.FindProductVariantList(tx, products[i].ProductID)
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
				price, err := u.Repository.Product.FindCurrentPrice(tx, variantList[j].VariantID)
				if err != nil {
					switch err {
					case sql.ErrNoRows:
						continue
					default:
						return nil, errors.New("не удалось найти актуальную цену продукта")
					}
				}

				variantList[j].CurrentPrice = price
				inStorages, err := u.Repository.Product.InStorages(tx, variantList[j].VariantID)
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
		products, err = u.Repository.Product.LoadProductList(tx, limit)
		if err != nil {
			return nil, err
		}

		for i := range products {
			vars, err := u.Repository.Product.FindProductVariantList(tx, products[i].ProductID)
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
				price, err := u.Repository.Product.FindCurrentPrice(tx, variants[j].VariantID)
				if err != nil {
					switch err {
					case sql.ErrNoRows:
						continue
					default:
						return nil, err
					}
				}

				variants[j].CurrentPrice = price
				inStorages, err := u.Repository.Product.InStorages(tx, variants[j].VariantID)
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

	return products, err
}

// FindProductsInStock логика получения всех складов и продуктов в ней или фильтрация по продукту
func (u ProductUseCaseImpl) FindProductsInStock(tx *sqlx.Tx, productID int) (stocks []stock.Stock, err error) {

	if productID < 0 {
		return nil, errors.New("id продукта не может быть меньше нуля")
	}

	// если пользователь не ввел id продукта то будет выполнен поиск всех складов
	if productID == 0 {
		stocks, err = u.Repository.Product.LoadStockList(tx)
		if err != nil {
			return nil, errors.New("не удалось найти склады")
		}
		for i, v := range stocks {
			variants, err := u.Repository.Product.FindStocksVariantList(tx, v.StorageID)
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
		stocks, err = u.Repository.Product.FindStockListByProductId(tx, productID)
		if err != nil {
			return nil, errors.New("не удалось найти склады с продуктами по данному id ")
		}
		for i, v := range stocks {
			variants, err := u.Repository.Product.FindStocksVariantList(tx, v.StorageID)
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

	return stocks, err
}

// Buy логuка записи о покупке в базу
func (u ProductUseCaseImpl) Buy(tx *sqlx.Tx, p product.Sale) (saleID int, err error) {

	// проверка фильтров на нулевые значения ,которые ввел пользователь
	err = p.IsNullFields()
	if err != nil {
		return 0, err
	}

	// устанавливаем текущую дату как дату продажи
	p.SoldAt = time.Now()

	// получение цены варианта
	price, err := u.Repository.Product.FindPrice(tx, p.VariantID)
	if err != nil {
		return 0, errors.New("не удалось найти цену продукта")
	}

	// подсчет общей цены продажи
	p.TotalPrice = price * float64(p.Quantity)

	// запись продажи в базу
	saleID, err = u.Repository.Product.Buy(tx, p)
	if err != nil {
		return 0, errors.New("не удалось записать продажу в базу")
	}

	return saleID, err
}

// FindSales получение списка всех продаж или списка продаж по фильтрам
func (u ProductUseCaseImpl) FindSaleList(tx *sqlx.Tx, sq product.SaleQuery) (sales []product.Sale, err error) {

	// если лимит не указан то по умолчанию устанавливается 3
	if !sq.Limit.Valid {
		sq.Limit.Scan(3)
	}

	// если не указано имя продукта или id склада то произойдет фильтрация только по датам
	if !sq.ProductName.Valid && !sq.StorageId.Valid {
		s := product.SaleQueryOnlyBySoldDate{
			StartDate: sq.StartDate,
			EndDate:   sq.EndDate,
			Limit:     sq.Limit,
		}

		sales, err = u.Repository.Product.FindSaleListOnlyBySoldDate(tx, s)
		if err != nil {
			return nil, errors.New("не удалось найти продажи")
		}
	} else {
		//  если имя продукта или id склада указан то произойдет фильтрация по этим параметрам
		sales, err = u.Repository.Product.FindSaleListByFilters(tx, sq)
		if err != nil {
			log.Println(err)
			return nil, errors.New("не удалось найти продажи по данным фильтрам")
		}
	}

	return sales, err
}
