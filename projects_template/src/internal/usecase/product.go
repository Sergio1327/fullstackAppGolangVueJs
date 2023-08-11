package usecase

import (
	"errors"
	"product_storage/internal/entity/global"
	"product_storage/internal/entity/product"
	"product_storage/internal/entity/stock"
	"product_storage/internal/transaction"
	"product_storage/rimport"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type ProductUseCase struct {
	log   *logrus.Logger
	dbLog *logrus.Logger
	rimport.RepositoryImports
}

func NewProduct(log, dblog *logrus.Logger, ri rimport.RepositoryImports) *ProductUseCase {
	return &ProductUseCase{
		log:               log,
		dbLog:             dblog,
		RepositoryImports: ri,
	}
}

// AddProduct логика добавление продукта в базу
func (u *ProductUseCase) AddProduct(ts transaction.Session, product product.Product) (productID int, err error) {
	// если имя продукта не введено то возвращается ошибка
	if product.Name == "" {
		return 0, errors.New("имя продукта не может быть пустым")
	}

	// добавляется продукт в базу
	productID, err = u.Repository.Product.AddProduct(ts, product)
	if err != nil {
		u.log.Error("не удалось добавить продукт", err)
		return 0, global.ErrInternalError
	}

	// если пользователь не ввел варианты продукта то данные о продукте просто запишутся в базу
	if product.VariantList == nil {
		u.log.WithFields(logrus.Fields{"product_id": productID}).Info("продукт успешно добавлен в базу данных")

		return productID, nil
	} else {
		// добавляются варианты продукта
		for _, v := range product.VariantList {
			err := u.Repository.Product.AddProductVariantList(ts, productID, v)
			if err != nil {
				u.log.Error("не удалось добавить варианты продукта", err)
				return 0, global.ErrInternalError
			}
		}

		u.log.WithFields(logrus.Fields{"product_id": productID}).Info("продукт успешно добавлен в базу данных")
		return productID, err
	}
}

// AddProductPrice логика проверки цены и вставки в базу
func (u *ProductUseCase) AddProductPrice(ts transaction.Session, p product.ProductPrice) (priceID int, err error) {
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
	isExistsID, err := u.Repository.Product.CheckExists(ts, p)
	if err != nil {
		switch err {
		case global.ErrNoData:
			isExistsID = 0
		default:
			u.log.Error("не удалось проверить наличие цен в базе данных", err)
			return 0, global.ErrInternalError
		}
	}

	// если запись уже имеется устанавливается дата окончания цены
	if isExistsID > 0 {
		p.EndDate.Scan(time.Now())
		err := u.Repository.Product.UpdateProductPrice(ts, p, isExistsID)
		if err != nil {
			u.log.Error("не удалось обновить цену", err)
			return 0, global.ErrInternalError
		}

		priceID = isExistsID
	} else {
		// если записи нет то цена вставляется в базу
		priceID, err = u.Repository.Product.AddProductPrice(ts, p)
		if err != nil {
			u.log.Error("не удалось добавить цену продукта в базу данных", err)
			return 0, global.ErrInternalError
		}
	}

	u.log.WithFields(logrus.Fields{"price_id": priceID}).Info("цена продукта успешно добавлена в базу данных")
	return priceID, err
}

// AddProductInStock логика проверка продукта на складе и обновления или добавления на базу
func (u *ProductUseCase) AddProductInStock(ts transaction.Session, p stock.AddProductInStock) (productStockID int, err error) {

	// проверка запроса на нулевые значения
	err = p.IsNullFields()
	if err != nil {
		return 0, err
	}

	// проверка есть ли уже продукт на складе
	isExist, err := u.Repository.Product.CheckProductInStock(ts, p)
	if err != nil {
		u.log.Error("не удалось проверить наличие продуктов на складе", err)
		return 0, global.ErrInternalError
	}

	// если продукт уже имеется в базе обновляется его кол-во
	if isExist {
		productStockID, err = u.Repository.Product.UpdateProductInstock(ts, p)
		if err != nil {
			u.log.Error("не удалось обновить кол-во продуктов на складе", err)
			return 0, global.ErrInternalError
		}
	} else {
		// если продукта нет на складе то он просто добавляется на склад
		productStockID, err = u.Repository.Product.AddProductInStock(ts, p)
		if err != nil {
			u.log.Error("не удалось добавить продукт на склад", err)
			return 0, global.ErrInternalError
		}
	}

	u.log.WithFields(logrus.Fields{"product_stock_id": productStockID}).Info("продукт успешно добавлен на склад")
	return productStockID, err
}

// FindProductInfoById логика получения всей информации о продукте и его вариантах по id
func (u *ProductUseCase) FindProductInfoById(ts transaction.Session, productID int) (productInfo product.ProductInfo, err error) {

	// если пользователь не ввел id выводится ошибка
	if productID <= 0 {
		return product.ProductInfo{}, errors.New("id не может быть меньше или равен 0")
	}

	// поиск продукта по его id
	productInfo, err = u.Repository.Product.LoadProductInfo(ts, productID)
	if err != nil {
		u.log.Error("не удалось найти информацию о продукте", err)
		return product.ProductInfo{}, global.ErrInternalError
	}

	productInfo.VariantList, err = u.Repository.Product.FindProductVariantList(ts, productInfo.ProductID)
	if err != nil {
		switch err {
		case global.ErrNoData:
			return productInfo, nil
		default:
			u.log.Error("не удалось найти варианты продукта", err)
			return product.ProductInfo{}, global.ErrInternalError
		}

	}

	for i, v := range productInfo.VariantList {
		price, err := u.Repository.Product.FindCurrentPrice(ts, v.VariantID)
		if err != nil {
			switch err {
			case global.ErrNoData:
				continue
			default:
				u.log.Error("не удалось найти актуальную цену варианта продукта", err)
				return product.ProductInfo{}, global.ErrInternalError
			}
		}

		// получение актуальной цены для каждого варианта продукта
		productInfo.VariantList[i].CurrentPrice = price

		// получение id складов в которых есть этот продукт
		inStorages, err := u.Repository.Product.InStorages(ts, v.VariantID)
		if err != nil {
			switch err {
			case global.ErrNoData:
				continue
			default:
				u.log.Error("не удалось найти склады в которых есть продукт", err)
				return product.ProductInfo{}, global.ErrInternalError
			}
		}

		productInfo.VariantList[i].InStorages = inStorages
	}

	u.log.WithFields(logrus.Fields{"product_info": productInfo}).Info("успешно получена информация о продукте")
	return productInfo, err
}

// FindProductList логика получения списка продуктов по тегу и лимиту
func (u *ProductUseCase) FindProductList(ts transaction.Session, tag string, limit int) (products []product.ProductInfo, err error) {

	// если лимит не указан или некорректен то по умолчанию устанавливается 3
	if limit == 0 || limit < 0 {
		limit = 3
	}

	// если пользователь ввел тег продукта произойдет поиск продуктов по данному тегу
	if tag != "" {
		products, err = u.Repository.Product.FindProductListByTag(ts, tag, limit)
		if err != nil {
			u.log.Error("не удалось найти продукты по данному тегу", err)
			return nil, global.ErrInternalError
		}

		for i := range products {
			vars, err := u.Repository.Product.FindProductVariantList(ts, products[i].ProductID)
			if err != nil {
				switch err {
				case global.ErrNoData:
					return products, nil
				default:
					u.log.Error("не удалось найти варианты продукта", err)
					return nil, global.ErrInternalError
				}
			}
			products[i].VariantList = vars
			variantList := products[i].VariantList
			for j := range variantList {
				price, err := u.Repository.Product.FindCurrentPrice(ts, variantList[j].VariantID)
				if err != nil {
					switch err {
					case global.ErrNoData:
						continue
					default:
						u.log.Error("не удалось найти актуальную цену продукта", err)
						return nil, global.ErrInternalError
					}
				}

				variantList[j].CurrentPrice = price
				inStorages, err := u.Repository.Product.InStorages(ts, variantList[j].VariantID)
				if err != nil {
					switch err {
					case global.ErrNoData:
						continue
					default:
						u.log.Error("не удалось найти склады в которых есть продукт", err)
						return nil, global.ErrInternalError
					}
				}

				variantList[j].InStorages = inStorages
			}
		}
	} else {
		// если пользователь не ввел тег то просто прозойдет поиск всех продуктов с лимитом вывода
		products, err = u.Repository.Product.LoadProductList(ts, limit)
		if err != nil {
			u.log.Error("не удалось найти список продуктов", err)
			return nil, global.ErrInternalError
		}

		for i := range products {
			vars, err := u.Repository.Product.FindProductVariantList(ts, products[i].ProductID)
			if err != nil {
				switch err {
				case global.ErrNoData:
					return products, nil
				default:
					u.log.Error("не удалось найти варианты продукта", err)
					return nil, global.ErrInternalError
				}
			}

			products[i].VariantList = vars
			variants := products[i].VariantList
			for j := range variants {
				price, err := u.Repository.Product.FindCurrentPrice(ts, variants[j].VariantID)
				if err != nil {
					switch err {
					case global.ErrNoData:
						continue
					default:
						u.log.Error("не удалось найти актуальную цену продукта", err)
						return nil, global.ErrInternalError
					}
				}

				variants[j].CurrentPrice = price
				inStorages, err := u.Repository.Product.InStorages(ts, variants[j].VariantID)
				if err != nil {
					switch err {
					case global.ErrNoData:
						continue
					default:
						u.log.Error("не удалось найти склады в которых есть продукт", err)
						return nil, global.ErrInternalError
					}
				}

				variants[j].InStorages = inStorages
			}
		}
	}

	u.log.WithFields(logrus.Fields{"product_list": products}).Info("успешно получен список продуктов")
	return products, err
}

// FindProductsInStock логика получения всех складов и продуктов в ней или фильтрация по продукту
func (u *ProductUseCase) FindProductsInStock(ts transaction.Session, productID int) (stocks []stock.Stock, err error) {

	if productID < 0 {
		return nil, errors.New("id продукта не может быть меньше нуля")
	}

	// если пользователь не ввел id продукта то будет выполнен поиск всех складов
	if productID == 0 {
		stocks, err = u.Repository.Product.LoadStockList(ts)
		if err != nil {
			u.log.Error("не удалось найти список складов", err)
			return nil, global.ErrInternalError
		}
		for i, v := range stocks {
			variants, err := u.Repository.Product.FindStocksVariantList(ts, v.StorageID)
			if err != nil {
				switch err {
				case global.ErrNoData:
					return stocks, nil
				default:
					u.log.Error("не удалось найти варианты продукта на складе", err)
					return nil, global.ErrInternalError
				}
			}

			stocks[i].ProductVariantList = variants
		}
	} else {

		//Если же пользователь ввел id продукта то произойдет фильтрация складов по id продукта
		stocks, err = u.Repository.Product.FindStockListByProductId(ts, productID)
		if err != nil {
			u.log.Error("не удалось найти склады с продуктами по данному id", err)
			return nil, global.ErrInternalError
		}
		for i, v := range stocks {
			variants, err := u.Repository.Product.FindStocksVariantList(ts, v.StorageID)
			if err != nil {
				switch err {
				case global.ErrNoData:
					return stocks, nil
				default:
					u.log.Error("не удалось найти варианты продукта на складе", err)
					return nil, global.ErrInternalError
				}
			}

			stocks[i].ProductVariantList = variants
		}
	}

	u.log.WithFields(logrus.Fields{"stock_list": stocks}).Info("успешно найдены склады и продукты в них")
	return stocks, err
}

// Buy логuка записи о покупке в базу
func (u *ProductUseCase) Buy(ts transaction.Session, p product.Sale) (saleID int, err error) {

	// проверка фильтров на нулевые значения ,которые ввел пользователь
	err = p.IsNullFields()
	if err != nil {
		return 0, err
	}

	// получение цены варианта
	price, err := u.Repository.Product.FindPrice(ts, p.VariantID)
	if err != nil {
		u.log.Error("не удалось найти цену варианта продукта", err)
		err = global.ErrInternalError
		return 0, err
	}

	// подсчет общей цены продажи
	p.TotalPrice = u.Repository.Product.CalculateTotalPrice(price, p.Quantity)
	// запись продажи в базу
	saleID, err = u.Repository.Product.Buy(ts, p)
	if err != nil {
		u.log.Error("не удалось записать продажу в базу", err)
		return 0, global.ErrInternalError
	}

	u.log.WithFields(logrus.Fields{"sale_id": saleID}).Info("продажа успешно добавлена в базу данных")
	return saleID, err
}

// FindSales получение списка всех продаж или списка продаж по фильтрам
func (u *ProductUseCase) FindSaleList(ts transaction.Session, sq product.SaleQueryParam) (sales []product.Sale, err error) {

	// если лимит не указан то по умолчанию устанавливается 3
	if !sq.Limit.Valid {
		sq.Limit.Scan(3)
	}

	// если не указано имя продукта или id склада то произойдет фильтрация только по датам
	if !sq.ProductName.Valid && !sq.StorageID.Valid {
		s := product.SaleQueryOnlyBySoldDateParam{
			StartDate: sq.StartDate,
			EndDate:   sq.EndDate,
			Limit:     sq.Limit,
		}

		sales, err = u.Repository.Product.FindSaleListOnlyBySoldDate(ts, s)
		if err != nil {
			u.log.Error("не удалось найти продажи", err)
			return nil, global.ErrInternalError
		}
	} else {
		//  если имя продукта или id склада указан то произойдет фильтрация по этим параметрам
		sales, err = u.Repository.Product.FindSaleListByFilters(ts, sq)
		if err != nil {
			u.log.Error("не удалось найти продажи по данным фильтрам", err)
			return nil, global.ErrInternalError
		}
	}

	u.log.WithFields(logrus.Fields{"sale_list": sales}).Info("успешно получены продажи по заданным фильтрам")
	return sales, err
}
