package usecase

import (
	"errors"
	"fmt"
	"product_storage/internal/entity/global"
	"product_storage/internal/entity/product"
	"product_storage/internal/entity/stock"
	"product_storage/internal/transaction"
	"product_storage/rimport"
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
func (u *ProductUseCase) AddProduct(ts transaction.Session, product product.ProductParams) (productID int, err error) {
	lf := product.Log()
	lf["product_params"] = product
	// если имя продукта не введено то возвращается ошибка
	if product.Name == "" {
		err = errors.New("имя продукта не может быть пустым")
		return
	}
	product.AddetAt = time.Now()
	// добавляется продукт в базу
	productID, err = u.Repository.Product.AddProduct(ts, product)
	if err != nil {
		u.log.WithFields(lf).Error("не удалось добавить продукт", err)
		err = global.ErrInternalError
		return
	}

	lf["product_ID"] = productID

	// если пользователь не ввел варианты продукта то данные о продукте просто запишутся в базу
	if product.VariantList == nil {
		u.log.WithFields(lf).Info("продукт успешно добавлен в базу данных")

		return productID, nil
	} else {
		// добавляются варианты продукта
		for _, v := range product.VariantList {
			err = u.Repository.Product.AddProductVariantList(ts, productID, v)
			if err != nil {
				u.log.WithFields(lf).Error("не удалось добавить варианты продукта", err)
				err = global.ErrInternalError
				return
			}
		}

		u.log.WithFields(lf).Info("продукт успешно добавлен в базу данных")
		return productID, err
	}
}

// AddProductPrice логика проверки цены и вставки в базу
func (u *ProductUseCase) AddProductPrice(ts transaction.Session, p product.ProductPriceParams) (priceID int, err error) {
	lf := p.Log()
	//проверка  id варианта, цены, даты начала цены на нулевые значения
	if err := p.IsNullFields(); err != nil {
		return 0, err
	}
	p.StartDate = time.Now()
	// проверка имеется ли запись уже в базе с заданным id продукта и дата начала цены
	isExistsID, err := u.Repository.Product.CheckExists(ts, p)
	switch err {
	case global.ErrNoData:
		// если записи нет то цена вставляется в базу
		priceID, err = u.Repository.Product.AddProductPrice(ts, p)
		if err != nil {
			u.log.WithFields(lf).Error("не удалось добавить цену продукта в базу данных", err)
			err = global.ErrInternalError
			return
		}

	case nil:
		// если запись уже имеется устанавливается дата окончания цены
		p.EndDate.Scan(time.Now())
		err = u.Repository.Product.UpdateProductPrice(ts, p, isExistsID)
		if err != nil {
			u.log.WithFields(lf).Error("не удалось обновить цену", err)
			err = global.ErrInternalError
			return
		}

		priceID = isExistsID

	default:
		u.log.Error("не удалось проверить наличие цен в базе данных", err)
		return 0, global.ErrInternalError
	}
	lf["price_ID"] = priceID

	u.log.WithFields(lf).Info("цена продукта успешно добавлена в базу данных")
	return priceID, err
}

// AddProductInStock логика проверка продукта на складе и обновления или добавления на базу
func (u *ProductUseCase) AddProductInStock(ts transaction.Session, p stock.ProductInStockParams) (productStockID int, err error) {
	lf := p.Log()
	lf["product_in_stock_params"] = p
	// проверка запроса на нулевые значения
	if err := p.IsNullFields(); err != nil {
		return 0, err
	}
	p.AddedAt = time.Now()
	// проверка есть ли уже продукт на складе
	isExist, err := u.Repository.Product.CheckProductInStock(ts, p)
	if err != nil {
		u.log.WithFields(lf).Error("не удалось проверить наличие продуктов на складе", err)
		err = global.ErrInternalError
		return
	}

	// если продукт уже имеется в базе обновляется его кол-во
	if isExist {
		productStockID, err = u.Repository.Product.UpdateProductInstock(ts, p)
		if err != nil {
			u.log.WithFields(lf).Error("не удалось обновить кол-во продуктов на складе", err)
			err = global.ErrInternalError
			return
		}
	} else {
		// если продукта нет на складе то он просто добавляется на склад
		productStockID, err = u.Repository.Product.AddProductInStock(ts, p)
		if err != nil {
			u.log.WithFields(lf).Error("не удалось добавить продукт на склад", err)
			err = global.ErrInternalError
			return
		}
	}
	lf["product_in_stock_ID"] = productStockID

	u.log.WithFields(lf).Info("продукт успешно добавлен на склад")
	return productStockID, err
}

// FindProductInfoById логика получения всей информации о продукте и его вариантах по id
func (u *ProductUseCase) FindProductInfoById(ts transaction.Session, productID int) (productInfo product.ProductInfo, err error) {
	lf := logrus.Fields{"product_ID": productID}
	// если пользователь не ввел id выводится ошибка
	if productID <= 0 {
		err = errors.New("id не может быть меньше или равен 0")
		return
	}

	// поиск продукта по его id
	productInfo, err = u.Repository.Product.LoadProductInfo(ts, productID)
	if err != nil {
		u.log.WithFields(lf).Error("не удалось найти информацию о продукте", err)
		err = global.ErrInternalError
		return
	}

	productInfo.VariantList, err = u.Repository.Product.FindProductVariantList(ts, productInfo.ProductID)
	switch err {
	case nil:
	case global.ErrNoData:
		return productInfo, nil
	default:
		u.log.WithFields(lf).Error("не удалось найти варианты продукта", err)
		return product.ProductInfo{}, global.ErrInternalError
	}

	for i, v := range productInfo.VariantList {
		price, err := u.Repository.Product.FindCurrentPrice(ts, v.VariantID)
		switch err {
		case nil:
		case global.ErrNoData:
			continue
		default:
			u.log.WithFields(lf).Error("не удалось найти актуальную цену варианта продукта", err)
			return product.ProductInfo{}, global.ErrInternalError
		}

		// получение актуальной цены для каждого варианта продукта
		productInfo.VariantList[i].CurrentPrice = price

		// получение id складов в которых есть этот продукт
		inStorages, err := u.Repository.Product.InStorages(ts, v.VariantID)
		switch err {
		case nil:
		case global.ErrNoData:
			continue
		default:
			u.log.WithFields(lf).Error("не удалось найти склады в которых есть продукт", err)
			return product.ProductInfo{}, global.ErrInternalError
		}

		productInfo.VariantList[i].InStorages = inStorages
	}

	return productInfo, err
}

// FindProductList логика получения списка продуктов по тегу и лимиту
func (u *ProductUseCase) FindProductList(ts transaction.Session, tag, name string, limit int) (products []product.ProductInfo, err error) {
	lf := logrus.Fields{"tag": tag, "limit": limit, "productName": name}
	// если лимит не указан или некорректен то по умолчанию устанавливается 3
	if limit == 0 || limit < 0 {
		limit = 3
	}

	// если пользователь ввел тег продукта произойдет поиск продуктов по данному тегу
	if tag != "" || name != "" {

		if tag != "" && name == "" {
			products, err = u.Repository.Product.FindProductListByTag(ts, tag, limit)
			if err != nil {
				u.log.WithFields(lf).Error("не удалось найти продукты по данному тегу", err)
				err = global.ErrInternalError
				return
			}
		} else if name != "" && tag == "" {
			products, err = u.Repository.Product.FindProductListByName(ts, name, limit)
			if err != nil {
				u.log.WithFields(lf).Error("не удалось найти продукты по данному тегу", err)
				err = global.ErrInternalError
				return
			}
		} else {
			products, err = u.Repository.Product.FindProductListByTagAndName(ts, tag, name, limit)
			if err != nil {
				u.log.WithFields(lf).Error("не удалось найти продукты по данному тегу", err)
				err = global.ErrInternalError
				return
			}
		}

		// поиск вариантов продука
		for i := range products {
			vars, err := u.Repository.Product.FindProductVariantList(ts, products[i].ProductID)
			switch err {
			case nil:
			case global.ErrNoData:
				return products, nil
			default:
				u.log.WithFields(lf).Error("не удалось найти варианты продукта", err)
				return nil, global.ErrInternalError
			}

			products[i].VariantList = vars
			variantList := products[i].VariantList
			for j := range variantList {
				price, err := u.Repository.Product.FindCurrentPrice(ts, variantList[j].VariantID)
				switch err {
				case nil:
				case global.ErrNoData:
					continue
				default:
					u.log.WithFields(lf).Error("не удалось найти актуальную цену продукта", err)
					return nil, global.ErrInternalError
				}

				variantList[j].CurrentPrice = price
				inStorages, err := u.Repository.Product.InStorages(ts, variantList[j].VariantID)
				switch err {
				case nil:
				case global.ErrNoData:
					continue
				default:
					u.log.Error("не удалось найти склады в которых есть продукт", err)
					return nil, global.ErrInternalError
				}

				variantList[j].InStorages = inStorages
			}
		}
	} else {
		// если пользователь не ввел тег то просто прозойдет поиск всех продуктов с лимитом вывода
		products, err = u.Repository.Product.LoadProductList(ts, limit)
		if err != nil {
			u.log.WithFields(lf).Error("не удалось найти список продуктов", err)
			err = global.ErrInternalError
			return
		}

		for i := range products {
			vars, err := u.Repository.Product.FindProductVariantList(ts, products[i].ProductID)
			switch err {
			case nil:
			case global.ErrNoData:
				return products, nil
			default:
				u.log.WithFields(lf).Error("не удалось найти варианты продукта", err)
				return nil, global.ErrInternalError
			}

			products[i].VariantList = vars
			variants := products[i].VariantList
			for j := range variants {
				price, err := u.Repository.Product.FindCurrentPrice(ts, variants[j].VariantID)
				switch err {
				case nil:
				case global.ErrNoData:
					continue
				default:
					u.log.WithFields(lf).Error("не удалось найти актуальную цену продукта", err)
					return nil, global.ErrInternalError
				}

				variants[j].CurrentPrice = price
				inStorages, err := u.Repository.Product.InStorages(ts, variants[j].VariantID)
				switch err {
				case nil:
				case global.ErrNoData:
					continue
				default:
					u.log.WithFields(lf).Error("не удалось найти склады в которых есть продукт", err)
					return nil, global.ErrInternalError
				}

				variants[j].InStorages = inStorages
			}
		}
	}

	return products, err
}

// FindProductsInStock логика получения всех складов и продуктов в ней или фильтрация по продукту
func (u *ProductUseCase) FindProductsInStock(ts transaction.Session, productID int) (stockList []stock.Stock, err error) {
	lf := logrus.Fields{"product_ID": productID}

	if productID < 0 {
		err = errors.New("id продукта не может быть меньше нуля")
		return
	}

	// если пользователь не ввел id продукта то будет выполнен поиск всех складов
	if productID == 0 {
		stockList, err = u.Repository.Product.LoadStockList(ts)
		if err != nil {
			u.log.WithFields(lf).Error("не удалось найти список складов", err)
			err = global.ErrInternalError
			return
		}

		for i, v := range stockList {
			variants, err := u.Repository.Product.FindStocksVariantList(ts, v.StorageID)
			switch err {
			case nil:
			case global.ErrNoData:
				return stockList, nil
			default:
				u.log.WithFields(lf).Error("не удалось найти варианты продукта на складе", err)
				return nil, global.ErrInternalError
			}

			stockList[i].ProductVariantList = variants
		}
	} else {

		// если же пользователь ввел id продукта то произойдет фильтрация складов по id продукта
		stockList, err = u.Repository.Product.FindStockListByProductId(ts, productID)
		if err != nil {
			u.log.WithFields(lf).Error("не удалось найти склады с продуктами по данному id", err)
			err = global.ErrInternalError
			return
		}
		for i, v := range stockList {
			variants, err := u.Repository.Product.FindStocksVariantList(ts, v.StorageID)
			switch err {
			case nil:
			case global.ErrNoData:
				return stockList, nil
			default:
				u.log.WithFields(lf).Error("не удалось найти варианты продукта на складе", err)
				return nil, global.ErrInternalError
			}

			stockList[i].ProductVariantList = variants
		}
	}
	lf["stock_list"] = stockList

	u.log.WithFields(lf).Info("успешно найдены склады и продукты в них")
	return stockList, err
}

// SaveSale логuка записи о покупке в базу
func (u *ProductUseCase) SaveSale(ts transaction.Session, p product.SaleParams) (saleID int, err error) {
	lf := p.Log()
	lf["sale_params"] = p

	// проверка фильтров на нулевые значения ,которые ввел пользователь
	if err := p.IsNullFields(); err != nil {
		return 0, err
	}

	// получение цены варианта
	price, err := u.Repository.Product.FindPrice(ts, p.VariantID)
	if err != nil {
		u.log.WithFields(lf).Error("не удалось найти цену варианта продукта ", err)
		err = global.ErrInternalError
		return
	}

	// подсчет общей цены продажи
	p.TotalPrice = price * float64(p.Quantity)
	if p.TotalPrice == 0 {
		u.log.WithFields(lf).Error("общая цена не может быть равна 0")
		err = global.ErrInternalError
		return
	}
	// запись продажи в базу
	saleID, err = u.Repository.Product.SaveSale(ts, p)
	if err != nil {
		u.log.WithFields(lf).Error("не удалось записать продажу в базу", err)
		err = global.ErrInternalError
		return
	}

	lf["sale_ID"] = saleID

	u.log.WithFields(lf).Info("продажа успешно добавлена в базу данных")
	return saleID, err
}

// FindSales получение списка всех продаж или списка продаж по фильтрам
func (u *ProductUseCase) FindSaleList(ts transaction.Session, sq product.SaleQueryParam) (saleList []product.Sale, err error) {
	lf := sq.Log()

	// если лимит не указан то по умолчанию устанавливается 3
	if !sq.Limit.Valid {
		sq.Limit.Scan(3)
	}

	if sq.StorageID.Int64 == 0 {
		sq.StorageID.Valid = false
	}

	// если не указано имя продукта или id склада то произойдет фильтрация только по датам
	if !sq.ProductName.Valid && !sq.StorageID.Valid {
		s := product.SaleQueryOnlyBySoldDateParam{
			StartDate: sq.StartDate,
			EndDate:   sq.EndDate,
			Limit:     sq.Limit,
		}
		lf = s.Log()

		saleList, err = u.Repository.Product.FindSaleListOnlyBySoldDate(ts, s)
		if err != nil {
			u.log.WithFields(lf).Error("не удалось найти продажи", err)
			err = global.ErrInternalError
			return
		}
	} else {
		//  если имя продукта или id склада указан то произойдет фильтрация по этим параметрам
		saleList, err = u.Repository.Product.FindSaleListByFilters(ts, sq)
		if err != nil {
			u.log.WithFields(lf).Error("не удалось найти продажи по данным фильтрам", err)
			err = global.ErrInternalError
			return
		}
	}

	return saleList, err
}

func (u *ProductUseCase) LoadStockList(ts transaction.Session) (stockList []stock.Stock, err error) {
	stockList, err = u.Repository.Product.LoadStockList(ts)
	lf := logrus.Fields{"stockList": stockList}

	switch err {
	case nil:
	case global.ErrNoData:
		return stockList, nil
	default:
		u.log.WithFields(lf).Error("не удалось найти склады ", err)
		return nil, global.ErrInternalError
	}

	return stockList, nil
}

func (u *ProductUseCase) AddStock(ts transaction.Session, storage stock.StockParams) (stockID int, err error) {
	lf := storage.Log()

	if storage.StorageName == "" {
		err = errors.New("название склада не может быть пустым")
		return
	}

	storage.Added_at.Scan(time.Now())

	stockID, err = u.Repository.Product.AddStock(ts, storage)
	if err != nil {
		u.log.WithFields(lf).Error("не удалось добавить склад ", err)
		err = global.ErrInternalError
		return
	}

	lf["stockID"] = stockID

	u.log.WithFields(lf).Info("склад успешно добавлен")
	return stockID, err
}

func (u *ProductUseCase) DeleteStock(ts transaction.Session, storage stock.StockParams) (err error) {
	lf := storage.Log()

	err = u.Repository.Product.DeleteStock(ts, storage)
	if err != nil {
		fmt.Println(1111)
		u.log.WithFields(lf).Error("не удалось удалить склад ", err)
		err = global.ErrInternalError
		return
	}

	u.log.WithFields(lf).Info("склад успешно удален")
	return err
}
