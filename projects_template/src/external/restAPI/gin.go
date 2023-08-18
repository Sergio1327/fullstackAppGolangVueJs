package restapi

import (
	"product_storage/uimport"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type GinServer struct {
	server *gin.Engine
	log    *logrus.Logger
	dbLog  *logrus.Logger
	uimport.UsecaseImports
}

func NewGinServer(log, dblog *logrus.Logger, U uimport.UsecaseImports) *GinServer {
	return &GinServer{
		log:            log,
		dbLog:          dblog,
		UsecaseImports: U,
	}
}

func (e *GinServer) Run() {
	e.server = gin.Default()

	e.server.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Замените * на список разрешенных доменов, если это необходимо
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	e.server.POST("/product/add", e.addProduct)
	e.server.POST("/product/price", e.addProductPrice)
	e.server.POST("/product/add/stock", e.addProductInStock)
	e.server.GET("/product/:id", e.findProductInfoById)
	e.server.GET("/product_list", e.findProductList)
	e.server.GET("/stock", e.findProductListInStock)
	e.server.POST("/buy", e.SaveSale)
	e.server.POST("/sales", e.FindSaleList)
	e.server.GET("/stock_list", e.LoadStockList)
	e.server.POST("/stock/add",e.AddStock)

	e.server.Run(":9000")
}
