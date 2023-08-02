package restapi

import (
	"log"
	"net/http"
	"product_storage/internal/entity/product"
	"product_storage/internal/entity/stock"
	"product_storage/internal/repository/postgresql"
	"product_storage/tools/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddProduct добавление товара
func (g *GinServer) addProduct(c *gin.Context) {
	ts := g.SessionManager.CreateSession()
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

	productID, err := g.Usecase.ProdcutUsecase.AddProduct(postgresql.SqlxTx(ts), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(productID, "product_id"))

	err = ts.Commit()
	if err != nil {
		log.Println(1)
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}
}

// AddProductPrice добавляет цену продукта
func (g *GinServer) addProductPrice(c *gin.Context) {
	ts := g.SessionManager.CreateSession()
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

	priceID, err := g.Usecase.ProdcutUsecase.AddProductPrice(postgresql.SqlxTx(ts), productPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(priceID, "price_id"))

	err = ts.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}
}

func (g *GinServer) addProductInStock(c *gin.Context) {
	ts := g.SessionManager.CreateSession()
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

	productStockID, err := g.Usecase.ProdcutUsecase.AddProductInStock(postgresql.SqlxTx(ts), addProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	err = ts.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(productStockID, "product_stock_ID"))
}

func (g *GinServer) findProductInfoById(c *gin.Context) {
	ts := g.SessionManager.CreateSession()
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

	productInfo, err := g.Usecase.ProdcutUsecase.FindProductInfoById(postgresql.SqlxTx(ts), productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	err = ts.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(productInfo, "product_info"))
}

func (g *GinServer) findProductList(c *gin.Context) {
	ts := g.SessionManager.CreateSession()
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

	productList, err := g.Usecase.ProdcutUsecase.FindProductList(postgresql.SqlxTx(ts), tag, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	err = ts.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(productList, "product_list"))
}

func (g *GinServer) findProductListInStock(c *gin.Context) {
	ts := g.SessionManager.CreateSession()
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

	stockList, err := g.Usecase.ProdcutUsecase.FindProductsInStock(postgresql.SqlxTx(ts), productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	err = ts.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(stockList, "stock_list"))
}

func (g *GinServer) buy(c *gin.Context) {
	ts := g.SessionManager.CreateSession()
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

	saleID, err := g.Usecase.ProdcutUsecase.Buy(postgresql.SqlxTx(ts), sale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	err = ts.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(saleID, "sale_id"))
}

func (g *GinServer) FindSaleList(c *gin.Context) {
	ts := g.SessionManager.CreateSession()
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

	saleList, err := g.Usecase.ProdcutUsecase.FindSaleList(postgresql.SqlxTx(ts), saleQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	err = ts.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(saleList, "sale_list"))
}
