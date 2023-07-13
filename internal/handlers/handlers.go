package handlers

import (
	"backend/db"
	"backend/entities"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

func AddProduct(db *db.Db) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var Product entities.Product
			err := json.NewDecoder(r.Body).Decode(&Product)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			tags, err := json.Marshal(Product.Tags)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			Product.Added_at = time.Now()
			var productId int
			err = db.QueryRow("insert into products(name,description,tags) values($1,$2,$3) returning product_id",
				Product.Name, Product.Description, tags).Scan(&productId)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			for _, variant := range Product.Variants {
				_, err := db.Exec("insert into product_variants(product_id,weight,unit,added_at) values($1,$2,$3,$4)",
					productId, variant.Weight, variant.Unit, Product.Added_at)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			var addedProductId int
			var addedVariantsIds []int
			err = db.QueryRow("select product_id from products where name=$1 and description=$2",
				Product.Name, Product.Description).Scan(&addedProductId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			rows, err := db.Query("select variant_id from product_variants where product_id=$1", addedProductId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			for rows.Next() {
				var addedVarId int
				err := rows.Scan(&addedVarId)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				addedVariantsIds = append(addedVariantsIds, addedVarId)
			}

			fmt.Fprintf(w, "ID продукта = %d, ID его вариантов = %d", addedProductId, addedVariantsIds)

		} else {
			http.Error(w, "err,enter correct method", http.StatusMethodNotAllowed)
		}
		w.Header().Set("Content-Type", "application/json")

	}
}

func AddPrice(db *db.Db) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var ProductPrice entities.ProductPrice
			err := json.NewDecoder(r.Body).Decode(&ProductPrice)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var isExistId int
			err = db.QueryRow("select price_id from product_prices where variant_id=$1 and start_date=$2 and (end_date>$3 or end_date is null)",
				ProductPrice.VariantId, ProductPrice.StartDate, ProductPrice.StartDate).Scan(&isExistId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			end_date := time.Now()
			if isExistId > 0 {
				_, err = db.Exec("update product_prices set end_date=$ where price_id=$2", end_date, isExistId)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			_, err = db.Exec("insert into product_prices(variant_id,price,start_date,end_date)values($1,$2,$3,$4)",
				ProductPrice.VariantId, ProductPrice.Price, ProductPrice.StartDate, ProductPrice.EndDate)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Fprint(w, ProductPrice)
		} else {
			http.Error(w, "error pls enter correct method", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")

	}
}

func AddInStock(db *db.Db) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var stock entities.Product_in_stock

		err := json.NewDecoder(r.Body).Decode(&stock)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var isExist bool
		err = db.QueryRow("select exists(select 1 from products_in_storage where variant_id=$1)").Scan(&isExist)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if isExist {
			_, err := db.Exec("update products_in_storage set quantity=$1 where variant_id=$2 and storage_id=$3", stock.Quantity, stock.VariantId, stock.StockId)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		_, err = db.Exec("insert into products_in_storage(variant_id,quantity,storage_id) values($1,$2,$3)", stock.VariantId, stock.Quantity, stock.StockId)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
	}
}

func GetProductInfo(db *db.Db) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			params := mux.Vars(r)
			productId := params["id"]
			var tags []byte
			Prodcut := entities.Product{}
			err := db.QueryRow("select product_id,name,description,tags from products where product_id=$1", productId).Scan(
				&Prodcut.ProductId, &Prodcut.Name, &Prodcut.Description, &tags)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = json.Unmarshal(tags, &Prodcut.Tags)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			rows, err := db.Query("select pv.variant_id, pv.weight, pv.unit from product_variants pv where pv.product_id = $1", productId)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer rows.Close()

			var variants []entities.ProductVariants
			for rows.Next() {
				var variant entities.ProductVariants
				err := rows.Scan(&variant.VariantId, &variant.Weight, &variant.Unit)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				variants = append(variants, variant)
			}
			Prodcut.Variants = variants

			for i := range variants {
				variant := &variants[i]
				err := db.QueryRow("select price from product_prices where variant_id = $1 and start_date <= NOW() and (end_date IS NULL OR end_date > NOW()) ORDER BY start_date DESC LIMIT 1", variant.VariantId).Scan(&variant.CurrentPrice)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}

			query := `
				select DISTINCT s.storage_id
				from storages s
				inner join products_in_storage pis ON s.storage_id = pis.storage_id
				inner join product_variants pv ON pis.variant_id = pv.variant_id
				where pv.product_id = $1
				`

			rows, err = db.Query(query, productId)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer rows.Close()

			var storageIDs []int
			for rows.Next() {
				var storageID int
				err := rows.Scan(&storageID)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				storageIDs = append(storageIDs, storageID)
			}

			var ProductInfo entities.ProductInfo
			ProductInfo.Product = Prodcut
			ProductInfo.StoragesIds = storageIDs
			result, err := json.Marshal(ProductInfo)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			err = json.NewEncoder(w).Encode(result)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

		} else {
			log.Println("enter another http method")
			http.Error(w, "enter another method", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
	}
}

func GetProductList(db *db.Db) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			var err error
			limitStr := r.URL.Query().Get("limit")
			tags := r.URL.Query().Get("tag")
			limit := 3
			if limitStr == "" {
				limit, err = strconv.Atoi(limitStr)
				if err != nil {
					log.Println(err)
					http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
					return
				}
			}
			query := "select * from products"
			if tags != "" {
				query += "where tags=$1"
			}
			query += "limit $2"

			rows, err := db.Query(query, strings.ReplaceAll(tags, "'", `"`), limit)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer rows.Close()

			products := []entities.Product{}
			for rows.Next() {
				product := entities.Product{}
				err := rows.Scan(&product.ProductId, &product.Name, &product.Description, &product.Tags)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				products = append(products, product)
			}
			if err = rows.Err(); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = json.NewEncoder(w).Encode(products)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
		} else {
			http.Error(w, "incorrect method", http.StatusMethodNotAllowed)
			return
		}
	}
}

func GetStockList(db *db.Db) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var variantId int
		variantIdstr := r.URL.Query().Get("variant_id")
		if variantIdstr != "" {
			var err error
			variantId, err = strconv.Atoi(variantIdstr)
			if err != nil {
				log.Println(err)
				http.Error(w, "Invalid variant ID", http.StatusBadRequest)
				return
			}
		}

		query := `
			select
				s.storage_id,
				s.name AS storage_name,
				pis.variant_id,
				pis.quantity
			from
				storages AS s
			inner join
				products_in_storage AS pis ON s.storage_id = pis.storage_id`
		args := []interface{}{}

		if variantId != 0 {
			query += `
			WHERE
				pis.variant_id = $1`
			args = append(args, variantId)
		}

		rows, err := db.Query(query, args...)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var stocks []entities.Stock

		for rows.Next() {
			var stock entities.Stock
			err := rows.Scan(&stock.StorageID, &stock.StorageName, &stock.VariantId, &stock.ProductCount)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			stocks = append(stocks, stock)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(stocks); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func BuyHandler(db *db.Db) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {

			var purchase []entities.Purchase
			err := json.NewDecoder(r.Body).Decode(&purchase)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			tx, err := db.DB.Begin()
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			for _, p := range purchase {
				var storage_id int
				err := db.QueryRow("select storage_id from products_in_storage where variant_id=$1 order by added_at asc limit 1", p.VariantId).Scan(&storage_id)
				if err != nil {
					log.Println(err)
					tx.Rollback()
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				var price decimal.Decimal

				err = db.QueryRow(`
				SELECT price FROM product_prices
				WHERE variant_id = $1
				AND start_date <= $2
				AND (end_date > $2 OR end_date IS NULL)
				LIMIT 1
			`, p.VariantId, time.Now()).Scan(&price)

				if err != nil {
					tx.Rollback()
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Product price not found")
					return
				}
				totalPrice := price.Mul(decimal.NewFromInt(int64(p.Quantity)))
				_, err = db.Exec("insert into sales(variant_id,storage_id,sold_at,quantity,total_price) values($1,$2,$3,$4,$5)",
					p.VariantId, storage_id, time.Now(), p.Quantity, totalPrice)
				if err != nil {
					log.Println(err)
					tx.Rollback()
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Product price not found")
					return
				}
				err = tx.Commit()
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, "Sales recorded successfully")
			}

		} else {
			http.Error(w, "incorrect method", http.StatusMethodNotAllowed)
			return
		}
	}
}

func GetSales(db *db.Db) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			startDateStr := r.URL.Query().Get("start_date")
			endDateStr := r.URL.Query().Get("end_date")
			limitStr := r.URL.Query().Get("limit")
			storageIDStr := r.URL.Query().Get("storage_id")

			if startDateStr == "" || endDateStr == "" {
				http.Error(w, "start_date and end_date are required fields", http.StatusBadRequest)
				return
			}

			startDate, err := time.Parse("2006-01-02", startDateStr)
			if err != nil {
				http.Error(w, "invalid start_date format", http.StatusBadRequest)
				return
			}
			endDate, err := time.Parse("2006-01-02", endDateStr)
			if err != nil {
				http.Error(w, "invalid end_date format", http.StatusBadRequest)
				return
			}

			limit := 3
			if limitStr != "" {
				limit, err = strconv.Atoi(limitStr)
				if err != nil {
					http.Error(w, "invalid limit format", http.StatusBadRequest)
					return
				}
			}

			query := `
				select s.sales_id, p.name, pv.weight, pv.unit, pp.price, s.quantity, s.total_price
				from sales s
				join product_variants pv ON s.variant_id = pv.variant_id
				join products p ON pv.product_id = p.product_id
				join product_prices pp ON pv.variant_id = pp.variant_id
				where s.sold_at >= $1 AND s.sold_at <= $2
			`

			params := []interface{}{startDate, endDate}

			if storageIDStr != "" {
				storageID, err := strconv.Atoi(storageIDStr)
				if err != nil {
					http.Error(w, "invalid storage_id format", http.StatusBadRequest)
					return
				}
				query += fmt.Sprintf("AND s.storage_id = %d ", storageID)
			}

			query += "ORDER BY s.sold_at DESC LIMIT $"
			query += strconv.Itoa(len(params) + 1)
			params = append(params, limit)

			rows, err := db.Query(query, params...)
			if err != nil {
				http.Error(w, "database query error", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var sales []entities.Sale
			for rows.Next() {
				var sale entities.Sale
				err := rows.Scan(&sale.ID, &sale.ProductName, &sale.VariantWeight, &sale.VariantUnit, &sale.Price, &sale.Quantity, &sale.TotalPrice)
				if err != nil {
					http.Error(w, "error while scanning rows", http.StatusInternalServerError)
					return
				}
				sales = append(sales, sale)
			}
			if err := rows.Err(); err != nil {
				http.Error(w, "error iterating over rows", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(sales); err != nil {
				http.Error(w, "error encoding JSON", http.StatusInternalServerError)
				return
			}

		} else {
			http.Error(w, "method not alowed", http.StatusMethodNotAllowed)
			return
		}
	}
}
