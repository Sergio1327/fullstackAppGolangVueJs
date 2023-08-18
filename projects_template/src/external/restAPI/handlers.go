package restapi

import (
	"log"
	"net/http"
	"product_storage/internal/entity/product"
	"product_storage/internal/entity/stock"
	"product_storage/tools/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// addProduct добавление товара
func (e *GinServer) addProduct(c *gin.Context) {
	ts := e.SessionManager.CreateSession()
	err := ts.Start()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}
	defer ts.Rollback()
	var product product.ProductParams

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	productID, err := e.Usecase.Product.AddProduct(ts, product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	if err := ts.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(productID, "product_id"))
}

// addProductPrice добавляет цену продукта
func (e *GinServer) addProductPrice(c *gin.Context) {
	ts := e.SessionManager.CreateSession()
	err := ts.Start()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}
	defer ts.Rollback()

	var productPrice product.ProductPriceParams

	if err := c.ShouldBindJSON(&productPrice); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	priceID, err := e.Usecase.Product.AddProductPrice(ts, productPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	if err := ts.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(priceID, "price_id"))
}

// addProductInStock добавляет продукт в склад
func (e *GinServer) addProductInStock(c *gin.Context) {
	ts := e.SessionManager.CreateSession()
	err := ts.Start()
	if err != nil {
		c.JSON(http.StatusOK, response.NewErrorResponse(err))
		return
	}
	defer ts.Rollback()

	var addProduct stock.ProductInStockParams

	if err := c.ShouldBindJSON(&addProduct); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	productStockID, err := e.Usecase.Product.AddProductInStock(ts, addProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	if err := ts.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(productStockID, "product_stock_ID"))
}

// findProductInfoById выводит данные о продукте по его id
func (e *GinServer) findProductInfoById(c *gin.Context) {
	ts := e.SessionManager.CreateSession()
	err := ts.Start()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
	}
	defer ts.Rollback()

	id := c.Param("id")
	productID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	productInfo, err := e.Usecase.Product.FindProductInfoById(ts, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(productInfo, "product_info"))
}

// findProductList выводит список продуктов по тегам и лимитам
func (e *GinServer) findProductList(c *gin.Context) {
	ts := e.SessionManager.CreateSession()
	err := ts.Start()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}
	defer ts.Rollback()

	tag := c.Query("tag")
	limitstr := c.Query("limit")
	limit, err := strconv.Atoi(limitstr)
	if err != nil {
		limit = 3
	}

	productList, err := e.Usecase.Product.FindProductList(ts, tag, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(productList, "product_list"))
}

// findProductListInStock выводит информацию о складах и продуктах в них
func (e *GinServer) findProductListInStock(c *gin.Context) {
	ts := e.SessionManager.CreateSession()
	err := ts.Start()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}
	defer ts.Rollback()

	id := c.Query("product_id")
	if id == "" {
		id = "0"
	}

	productId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	stockList, err := e.Usecase.Product.FindProductsInStock(ts, productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(stockList, "stock_list"))
}

// buy запись сделанной продажи в базу
func (e *GinServer) SaveSale(c *gin.Context) {
	ts := e.SessionManager.CreateSession()
	err := ts.Start()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}
	defer ts.Rollback()

	var sale product.SaleParams

	if err := c.ShouldBindJSON(&sale); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}
	sale.SoldAt = time.Now()

	saleID, err := e.Usecase.Product.SaveSale(ts, sale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	if err := ts.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(saleID, "sale_id"))
}

// findSales выводит информацию о продажах по фильтрам или без них
func (e *GinServer) FindSaleList(c *gin.Context) {
	ts := e.SessionManager.CreateSession()
	err := ts.Start()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}
	defer ts.Rollback()

	var saleQuery product.SaleQueryParam

	if err := c.ShouldBindJSON(&saleQuery); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}
	log.Println(saleQuery)

	if saleQuery.ProductName.String == "" {
		saleQuery.ProductName.Valid = false
	}

	saleList, err := e.Usecase.Product.FindSaleList(ts, saleQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(saleList, "sale_list"))
}

func (e *GinServer) LoadStockList(c *gin.Context) {
	ts := e.SessionManager.CreateSession()
	err := ts.Start()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}
	defer ts.Rollback()

	stockList, err := e.Usecase.Product.LoadStockList(ts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(stockList, "stock_list"))
}
