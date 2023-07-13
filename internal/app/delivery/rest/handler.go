package rest

import (
	"fmt"
	"go-back/internal/app/domain"
	"go-back/internal/app/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (ph *ProductHandler) AddProduct(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.String(400, fmt.Sprintf("%s", err))
		return
	}

	err := ph.service.AddProduct(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, "the product was added")
}
func (ph *ProductHandler) AddProductPrice(c *gin.Context) {
	var productPrice domain.ProductPrice
	if err := c.ShouldBindJSON(&productPrice); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	err := ph.service.AddProductPrice(&productPrice)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "the price was added")
}

func (ph *ProductHandler) AddProductInStock(c *gin.Context) {
	var addProduct domain.AddProductInStock
	if err := c.ShouldBindJSON(&addProduct); err != nil {
		c.String(http.StatusBadRequest, err.Error(), "1")
		return
	}
	err := ph.service.AddProductInStock(&addProduct)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error(), "2")
		return
	}
	c.JSON(http.StatusOK, "The product was succesfuly added in storage")
}

func (ph *ProductHandler) GetProductInfoById(c *gin.Context) {
	id := c.Param("id")
	productId, err := strconv.Atoi(id)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error(), 1)
		return
	}
	productInfo, err := ph.service.GetProductInfoById(productId)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error(), 2)
		return
	}

	c.JSON(http.StatusOK, productInfo)

}

func (ph *ProductHandler) GetProductList(c *gin.Context) {
	tag := c.Query("tag")
	limitstr := c.Query("limit")
	limit, err := strconv.Atoi(limitstr)
	if err != nil {
		limit = 3
	}

	products, err := ph.service.GetProductList(tag, limit)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, products)
}

func (ph *ProductHandler) GetProductsInStock(c *gin.Context) {
	id := c.Query("product_id")
	if id == "" {
		id = "0"
	}
	productId, err := strconv.Atoi(id)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	stocks, err := ph.service.GetProductsInStock(productId)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, stocks)
}

func (ph *ProductHandler) Buy(c *gin.Context) {
	var sale domain.Sale
	if err := c.ShouldBindJSON(&sale); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	err := ph.service.Buy(&sale)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "the sale was added")
}

func (ph *ProductHandler) GetSales(c *gin.Context) {
	var salequery domain.SaleQuery
	if err := c.ShouldBindJSON(&salequery); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	sales, err := ph.service.GetSales(&salequery)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, sales)
}
