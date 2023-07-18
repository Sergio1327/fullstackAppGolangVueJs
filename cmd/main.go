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

	//инициализация слоев 
	repo := repository.NewPostgresProductRepository(db)
	useCase := service.NewProductUseCase(repo)
	handler := rest.NewProductHandler(useCase)

	r := gin.Default()

//инициализация эндпоинтов и обработчиков
	r.POST("/product/add", handler.AddProduct)
	r.POST("/product/price", handler.AddProductPrice)
	r.POST("/product/add/stock", handler.AddProductInStock)
	r.GET("/product/:id", handler.FindProductInfoById)
	r.GET("/product_list", handler.LoadProductList)
	r.GET("/stock", handler.LoadProductsInStock)
	r.POST("/buy", handler.Buy)
	r.POST("/sales", handler.LoadSales)

	//запуск сервера
	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
