package main

import (
	"go-back/internal/app/delivery/rest"
	"go-back/internal/app/repository"
	"go-back/internal/app/service"
	"go-back/internal/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	conStr := "dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable"
	db, err := database.NewPostgreSQLdb(conStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewPostgresProductRepository(db)
	useCase := service.NewProductUseCase(repo)
	handler := rest.NewProductHandler(useCase)

	r := gin.Default()

	r.POST("/product/add", handler.AddProduct)
	r.POST("/product/price", handler.AddProductPrice)
	r.POST("/product/add/stock", handler.AddProductInStock)
	r.GET("/product/:id", handler.GetProductInfoById)
	r.GET("/product_list", handler.GetProductList)
	r.GET("/stock", handler.GetProductsInStock)
	r.POST("/buy", handler.Buy)
	r.POST("/sales", handler.GetSales)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
