package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Product struct {
	Name         string `json:"name"`
	ImageURL     string `json:"image_url"`
	Description  string `json:"description"`
	Price        string `json:"price"`
	TotalReviews string `json:"total_reviews"`
}
type ProductDetails struct {
	Url             string  `json:"url"`
	Product         Product `json:"product"`
	CreationTime    string  `json:"creation_time"`
	LastUpdatedTime string  `json:"last_updated_time"`
}

var db *sql.DB

//-------------check whether given table is present ot not---------
func isTablePresent(tableName string) bool {
	_, err := db.Query("select * from " + tableName + ";")
	if err == nil {
		return true
	}
	log.Println(err)
	return false
}

//---------------insert product data into the database-------------------------
func insertProduct(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var productDetails ProductDetails

	//read payload
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeResponse(w, "Oops....unable to read body!")
		return
	}

	//convert body into our product struct
	err = json.Unmarshal(reqBody, &productDetails)
	if err != nil {
		writeResponse(w, "Oops....unable to unmarshal body to product struct")
		return
	}

	//if product table doesn't exist then create one
	if !isTablePresent("product") {
		_, err := db.Exec("CREATE TABLE `test`.`product` (`name` VARCHAR(50) NOT NULL ,`description` LONGTEXT NOT NULL ,`url` VARCHAR(250) NOT NULL ,`imageUrl` TEXT NOT NULL ,`price` VARCHAR(20) NOT NULL ,`totalReviews` VARCHAR(20) NOT NULL ,`lastUpdatedTime` VARCHAR(250) NOT NULL ,`creationTime` VARCHAR(250) NOT NULL ,PRIMARY KEY (`url`));")
		if err != nil {
			log.Println("unable to create table product, Reason:", err)
			writeResponse(w, fmt.Sprintf("%s",err))
			return
		}
		log.Println("product table created successfully!")
	}

	//insert data into product table
	insForm, err := db.Prepare("INSERT INTO `product` (`name`, `description`, `url`, `imageUrl`, `price`, `totalReviews`, `creationTime`) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("unable to prepare insert query, Reason:", err)
		return
	}

	//check whether body is empty or not to avoid run time error
	if (productDetails != ProductDetails{}) && (productDetails.Product != Product{}) {
		url := productDetails.Url
		name := productDetails.Product.Name
		desc := productDetails.Product.Description
		price := productDetails.Product.Price
		reviews := productDetails.Product.TotalReviews
		imageUrl := productDetails.Product.ImageURL
		creationTime := productDetails.CreationTime
		_, err := insForm.Exec(name, desc, url, imageUrl, price, reviews, creationTime)
		if err != nil {
			log.Println("unable to insert data into table product, Reason:", err)
		} else {
			log.Println("Data inserted successfully into product table!")
		}
	}

}

//-----------------list all products available in data store--------------------
func listProducts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var product ProductDetails
	var allProducts []ProductDetails

	//check whether product table present or not
	if !isTablePresent("product") {
		log.Println("product table not present!")
		writeResponse(w, "Oops....No data available, please try again later!")
		return
	}

	//fetch all rows from table product.
	products, err := db.Query("SELECT * FROM `product`")
	if err != nil {
		log.Println("product table not present!")
		writeResponse(w, "Oops....Something went wrong!, please try again later!")
		return
	}

	//traverse row by row and append to array of products.
	for products.Next() {
		var name, description, url, imageUrl, price, totalReviews, creationTime, lastUpdateTime string
		err = products.Scan(&name, &description, &url, &imageUrl, &price, &totalReviews, &creationTime, &lastUpdateTime)
		if err != nil {
			panic(err.Error())
		}
		product = ProductDetails{
			Url: url,
			Product: Product{
				Name:         name,
				ImageURL:     imageUrl,
				Description:  description,
				Price:        price,
				TotalReviews: totalReviews,
			},
			CreationTime:    creationTime,
			LastUpdatedTime: lastUpdateTime,
		}
		allProducts = append(allProducts, product)
	}

	//convert array of products struct into array of json and write response back to browser
	responseBody, err := json.Marshal(allProducts)
	if err != nil {
		log.Println("unable to convert struct to json, Reason: ", err)
		writeResponse(w, "Oops....Something went wrong!, please try again later!")
		return
	}
	_, err = w.Write(responseBody)
	if err != nil {
		log.Println("error occurred while writing response for get products, Reason: ", err)
		writeResponse(w, "Oops....Something went wrong!, please try again later!")
		return
	}
}

//---------------------write response back to api server--------------------------
func writeResponse(w http.ResponseWriter, msg string) {
	_, err := w.Write([]byte(msg))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("unable to write response, Reason:", err)
	}
}

//---initialize db and set reference variable to it---------------
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Set an Environment Variables
	_ = os.Setenv("DB_HOST", "localhost")
	_ = os.Setenv("DB_USERNAME", "root")
	_ = os.Setenv("DB_PASSWORD", "")
	_ = os.Setenv("DB_NAME", "test")

	// Get the value of an Environment Variable
	host := os.Getenv("DB_HOST")
	dbUserName := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dataSource := fmt.Sprintf("%s:%s@(%s:3306)/%s", dbUserName, dbPass, host, dbName)
	dBConnection, err := sql.Open("mysql", dataSource)
	if err != nil {
		log.Println("Db Connection Failed!!, Reason:", err)
		return
	}

	err = dBConnection.Ping()
	if err != nil {
		log.Println("Ping Failed!!, Reason:", err)
		return
	}

	log.Printf("Connected to DB %s successfully\n", dbName)
	db = dBConnection
	dBConnection.SetMaxOpenConns(10)
	dBConnection.SetMaxIdleConns(5)
	dBConnection.SetConnMaxLifetime(time.Second * 10)
}

func main() {
	defer db.Close()
	mux := httprouter.New()
	mux.GET("/products", listProducts)
	mux.POST("/insertProduct", insertProduct)
	log.Fatalln(http.ListenAndServe(":8081", mux))
}

