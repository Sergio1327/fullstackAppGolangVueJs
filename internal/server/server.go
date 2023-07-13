package server

import (
	"backend/db"
	"backend/internal/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func Run(db *db.Db) error {
	r := mux.NewRouter()

	r.HandleFunc("/product/add", handlers.AddProduct(db))
	r.HandleFunc("/product/price", handlers.AddPrice(db))
	r.HandleFunc("/product/add/stock", handlers.AddInStock(db))
	r.HandleFunc("/product/{id}", handlers.GetProductInfo(db))
	r.HandleFunc("/product_list", handlers.GetProductList(db))
	r.HandleFunc("/stock", handlers.GetStockList(db))
	r.HandleFunc("/buy", handlers.BuyHandler(db))
	r.HandleFunc("/sales", handlers.GetSales(db))

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		return err
	}

	return nil
}
