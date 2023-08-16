package rest

import (
	"go-back/internal/app/domain"
	"go-back/internal/app/service"
	"go-back/internal/tools/response"
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
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	productID, err := ph.productService.AddProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(productID, "product_id"))
}

// AddProductPrice добавляет цену продукта
func (ph *ProductHandler) AddProductPrice(c *gin.Context) {
	var productPrice domain.ProductPrice

	if err := c.ShouldBindJSON(&productPrice); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	priceID, err := ph.productService.AddProductPrice(productPrice)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(priceID, "price_id"))
}

// AddProductInStock добавляет продукт в склад
func (ph *ProductHandler) AddProductInStock(c *gin.Context) {
	var addProduct domain.AddProductInStock

	if err := c.ShouldBindJSON(&addProduct); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	productStockID, err := ph.productService.AddProductInStock(addProduct)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(productStockID, "product_in_stock_ID"))
}

// FindProductInfoById выводит данные о продукте по его id
func (ph *ProductHandler) FindProductInfoById(c *gin.Context) {
	id := c.Param("id")
	productID, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	productInfo, err := ph.productService.FindProductInfoById(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(productInfo, "product_info"))
}

// LoadProductList выводит список продуктов по тегам и лимитам
func (ph *ProductHandler) FindProductList(c *gin.Context) {
	tag := c.Query("tag")
	limitstr := c.Query("limit")
	limit, err := strconv.Atoi(limitstr)
	if err != nil {
		limit = 3
	}

	productList, err := ph.productService.FindProductList(tag, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(productList, "product_list"))
}

// FindProductListInStock выводит информацию о складах и продуктах в них
func (ph *ProductHandler) FindProductListInStock(c *gin.Context) {
	id := c.Query("product_id")
	if id == "" {
		id = "0"
	}
	productId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	stockList, err := ph.productService.FindProductsInStock(productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(stockList, "stock_list"))
}

// Buy запись сделанной продажи в базу
func (ph *ProductHandler) Buy(c *gin.Context) {
	var sale domain.Sale

	if err := c.ShouldBindJSON(&sale); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	saleID, err := ph.productService.Buy(sale)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(saleID, "sale_id"))
}

// FindSales выводит информацию о продажах по фильтрам или без них
func (ph *ProductHandler) FindSaleList(c *gin.Context) {
	var salequery domain.SaleQuery

	if err := c.ShouldBindJSON(&salequery); err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(err))
		return
	}

	saleList, err := ph.productService.FindSales(salequery)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, response.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, response.NewSuccessResponse(saleList, "sale_list"))
}
