package restapi

import (
	"log"
	"net/http"
	"product_storage/internal/entity/product"
	"product_storage/internal/repository/postgresql"
	"product_storage/tools/response"

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
