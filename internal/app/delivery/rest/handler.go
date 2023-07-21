package rest

import (
	"go-back/internal/app/domain"
	"go-back/internal/app/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: service,
	}
}

// AddProduct обработчик который принимает json парсит данные в структуру и отправляет данные в метод AddProduct в слой сервиса
func (ph *ProductHandler) AddProduct(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	err := ph.productService.AddProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, "the product was added")
}

// AddProductPrice обработчик который принимает json парсит данные в структуру и отправляет данные в
// метод productService.AddProductPrice
func (ph *ProductHandler) AddProductPrice(c *gin.Context) {
	var productPrice domain.ProductPrice
	if err := c.ShouldBindJSON(&productPrice); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	err := ph.productService.AddProductPrice(productPrice)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, "the price was added")
}

// AddProductInStock парсит данные в структуру и добавляет продукт в склад через product.service
func (ph *ProductHandler) AddProductInStock(c *gin.Context) {
	var addProduct domain.AddProductInStock
	if err := c.ShouldBindJSON(&addProduct); err != nil {
		c.String(http.StatusBadRequest, err.Error(), "1")
		return
	}
	err := ph.productService.AddProductInStock(addProduct)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error(), "2")
		return
	}
	c.JSON(http.StatusOK, "The product was succesfuly added in storage")
}

// FindProductInfoById  парсит данные в структуру и выводит ионформацию о продукте по его id
func (ph *ProductHandler) FindProductInfoById(c *gin.Context) {
	id := c.Param("id")
	productId, err := strconv.Atoi(id)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error(), 1)
		return
	}
	productInfo, err := ph.productService.FindProductInfoById(productId)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error(), 2)
		return
	}

	c.JSON(http.StatusOK, productInfo)
}

// LoadProductList выводит список продуктов по тегам и лимитам
func (ph *ProductHandler) LoadProductList(c *gin.Context) {
	tag := c.Query("tag")
	limitstr := c.Query("limit")
	limit, err := strconv.Atoi(limitstr)
	if err != nil {
		limit = 3
	}

	products, err := ph.productService.FindProductList(tag, limit)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, products)
}

// LoadProductsInStock выводит информацию о складах и продуктах в них
func (ph *ProductHandler) LoadProductsInStock(c *gin.Context) {
	id := c.Query("product_id")
	if id == "" {
		id = "0"
	}
	productId, err := strconv.Atoi(id)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	stocks, err := ph.productService.FindProductsInStock(productId)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, stocks)
}

// Buy  запись сделанной продажи в базу
func (ph *ProductHandler) Buy(c *gin.Context) {
	var sale domain.Sale
	if err := c.ShouldBindJSON(&sale); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	err := ph.productService.Buy(sale)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "the sale was added")
}

// FindSales выводит информацию о продажах по фильтрам или без них
func (ph *ProductHandler) FindSales(c *gin.Context) {
	var salequery domain.SaleQuery
	if err := c.ShouldBindJSON(&salequery); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	sales, err := ph.productService.FindSales(salequery)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, sales)
}
