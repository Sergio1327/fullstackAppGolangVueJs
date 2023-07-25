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

// AddProduct добавление товара
func (ph *ProductHandler) AddProduct(c *gin.Context) {
	var product domain.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, "Введены некоректные данные")
		return
	}

	productId, err := ph.productService.AddProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Не удалось добавить продукт")
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, productId)
}

// AddProductPrice добавляет цену продукта
func (ph *ProductHandler) AddProductPrice(c *gin.Context) {
	var productPrice domain.ProductPrice
	if err := c.ShouldBindJSON(&productPrice); err != nil {
		c.String(http.StatusBadRequest, "Введенны некоректные данные")
		return
	}
	err := ph.productService.AddProductPrice(productPrice)
	if err != nil {
		c.String(http.StatusBadRequest, "Не удалось добавить цену продукта")
		return
	}
	c.JSON(http.StatusOK, "the price was added")
}

// AddProductInStock добавляет продукт в склад
func (ph *ProductHandler) AddProductInStock(c *gin.Context) {
	var addProduct domain.AddProductInStock
	if err := c.ShouldBindJSON(&addProduct); err != nil {
		c.String(http.StatusBadRequest, "Введены некорректные данные")
		return
	}
	err := ph.productService.AddProductInStock(addProduct)
	if err != nil {
		c.String(http.StatusInternalServerError, "Не удалось добавить продукт на склад")
		return
	}
	c.JSON(http.StatusOK, "The product was succesfuly added in storage")
}

// FindProductInfoById  выводит данные о продукте по его id
func (ph *ProductHandler) FindProductInfoById(c *gin.Context) {
	id := c.Param("id")
	productId, err := strconv.Atoi(id)
	if err != nil {
		c.String(http.StatusBadRequest, "Неверный или некоректный id")
		return
	}
	productInfo, err := ph.productService.FindProductInfoById(productId)
	if err != nil {
		c.String(http.StatusBadRequest, "Не удалось найти информацию о продукте")
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
		c.String(http.StatusInternalServerError, "Не удалось найти список продуктов")
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
		c.String(http.StatusBadRequest, "Неверный или некоректный id")
		return
	}
	stocks, err := ph.productService.FindProductsInStock(productId)
	if err != nil {
		c.String(http.StatusInternalServerError, "Не удалось найти продукты на складе")
	}
	
	c.JSON(http.StatusOK, stocks)
}

// Buy  запись сделанной продажи в базу
func (ph *ProductHandler) Buy(c *gin.Context) {
	var sale domain.Sale
	if err := c.ShouldBindJSON(&sale); err != nil {
		c.String(http.StatusBadRequest, "Введены неверные или некоректные данные")
		return
	}
	err := ph.productService.Buy(sale)
	if err != nil {
		c.String(http.StatusInternalServerError, "Не удалось добавить продажи в базу данных")
		return
	}

	c.String(http.StatusOK, "the sale was added")
}

// FindSales выводит информацию о продажах по фильтрам или без них
func (ph *ProductHandler) FindSales(c *gin.Context) {
	var salequery domain.SaleQuery
	if err := c.ShouldBindJSON(&salequery); err != nil {
		c.String(http.StatusBadRequest, "Неверные или некоректные фильтры")
		return
	}

	sales, err := ph.productService.FindSales(salequery)
	if err != nil {
		c.String(http.StatusInternalServerError, "Не удалось найти продажи по данным фильтрам")
		return
	}

	c.JSON(http.StatusOK, sales)
}
