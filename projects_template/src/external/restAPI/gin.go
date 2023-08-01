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

func (g *GinServer) Run() {
	g.server = gin.Default()

	g.server.POST("/product/add", g.addProduct)
	g.server.POST("/product/price", g.addProductPrice)

	g.server.Run(":8080")
}
