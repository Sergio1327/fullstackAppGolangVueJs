package restapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"product_storage/internal/entity/product"
	"product_storage/internal/entity/stock"
	"product_storage/tools/response"
	"strconv"
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
	var product product.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	productID, err := e.Usecase.ProdcutUsecase.AddProduct(ts, product)
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

	var productPrice product.ProductPrice

	if err := c.ShouldBindJSON(&productPrice); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	priceID, err := e.Usecase.ProdcutUsecase.AddProductPrice(ts, productPrice)
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

	var addProduct stock.AddProductInStock

	if err := c.ShouldBindJSON(&addProduct); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	productStockID, err := e.Usecase.ProdcutUsecase.AddProductInStock(ts, addProduct)
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

	productInfo, err := e.Usecase.ProdcutUsecase.FindProductInfoById(ts, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	if err := ts.Commit(); err != nil {
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

	productList, err := e.Usecase.ProdcutUsecase.FindProductList(ts, tag, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	if err := ts.Commit(); err != nil {
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

	stockList, err := e.Usecase.ProdcutUsecase.FindProductsInStock(ts, productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	if err := ts.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(stockList, "stock_list"))
}

// buy запись сделанной продажи в базу
func (e *GinServer) buy(c *gin.Context) {
	ts := e.SessionManager.CreateSession()
	err := ts.Start()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}
	defer ts.Rollback()

	var sale product.Sale

	if err := c.ShouldBindJSON(&sale); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	saleID, err := e.Usecase.ProdcutUsecase.Buy(ts, sale)
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

	var saleQuery product.SaleQuery

	if err := c.ShouldBindJSON(&saleQuery); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	saleList, err := e.Usecase.ProdcutUsecase.FindSaleList(ts, saleQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	if err := ts.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}
	
	c.JSON(http.StatusOK, response.NewSuccessResponse(saleList, "sale_list"))
}
