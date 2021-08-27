package datastore

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

const (
	tableName     = "product"
	defaultDBHost = "localhost"
	defaultDBUser = "root"
	defaultDBName = "test"
)

var db *sql.DB

type ProductDetails struct {
	URL             string   `json:"url"`
	Product         *Product `json:"product"`
	CreationTime    string   `json:"creation_time"`
	LastUpdatedTime string   `json:"last_updated_time"`
}

type Product struct {
	Name         string `json:"name"`
	ImageURL     string `json:"image_url"`
	Description  string `json:"description"`
	Price        string `json:"price"`
	TotalReviews string `json:"total_reviews"`
}

func DBConn() (*sql.DB, error) {
	var err error
	var host, dbUser, dbPass, dbName string

	// Get db config variables
	if host = os.Getenv("MYSQL_HOST"); host == "" {
		host = defaultDBHost
	}
	if dbUser = os.Getenv("MYSQL_USER"); dbUser == "" {
		dbUser = defaultDBUser
	}
	if dbName = os.Getenv("MYSQL_DATABASE"); dbName == "" {
		dbName = defaultDBName
	}
	dbPass = os.Getenv("MYSQL_ROOT_PASSWORD")
	log.Infof("Trying to connect with DB with info, host=%s, user=%s, db=%s", host, dbUser, dbName)
	dataSource := fmt.Sprintf("%s:%s@(%s:3306)/%s", dbUser, dbPass, host, dbName)
	retryCount := 1
	for {
		db, err = sql.Open("mysql", dataSource)
		if err == nil {
			break
		}
		time.Sleep((1 << retryCount) * 5 * time.Second) // wait 5s, 10s, 20ms, 40s ...
		log.Warnf("db Connection Failed!!, error: %v", err)
		if retryCount > 20 {
			retryCount = 1
		}
	}
	err = db.Ping()
	if err != nil {
		log.Warnf("Ping Failed!!, error: %v", err)
	}
	log.Printf("Successfully connected to db %s ", dbName)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Second * 10)
	log.Printf("Successfully connected to Postgres db")
	return db, nil
}

//-------------check whether given table is present ot not---------
func isTablePresent(tableName string) bool {
	_, err := db.Query("select * from " + tableName + ";")
	if err == nil {
		return true
	}
	log.Println(err)
	return false
}

// InsertProduct insert product data into the database
func InsertProduct(product *ProductDetails) {
	// if product table doesn't exist then create one
	if !isTablePresent(tableName) {
		_, err := db.Exec("CREATE TABLE `test`.`product` (`name` VARCHAR(50) NOT NULL ,`description` LONGTEXT NOT NULL ,`url` VARCHAR(250) NOT NULL ,`imageUrl` TEXT NOT NULL ,`price` VARCHAR(20) NOT NULL ,`totalReviews` VARCHAR(20) NOT NULL ,`lastUpdatedTime` VARCHAR(250) NOT NULL ,`creationTime` VARCHAR(250) NOT NULL ,PRIMARY KEY (`url`));")
		if err != nil {
			log.Errorf("unable to create table %s, %v", tableName, err)
			return
		}
		log.Infof("table %s created successfully!", tableName)
	}

	// insert data into product table
	insForm, err := db.Prepare("INSERT INTO `product` (`name`, `description`, `url`, `imageUrl`, `price`, `totalReviews`, `creationTime`) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Errorf("unable to prepare insert query, %v", err)
		return
	}
	_, err = insForm.Exec(product.Product.Name, product.Product.Description, product.URL, product.Product.ImageURL, product.Product.Price, product.Product.TotalReviews, product.CreationTime)
	if err != nil {
		log.Errorf("unable to insert data into table %s, %v", tableName, err)
		return
	}
	log.Println("Data inserted successfully into product table!")
}

// ListProducts list all products available in data store.
func ListProducts(rw http.ResponseWriter, req *http.Request) {
	var product ProductDetails
	var allProducts []ProductDetails

	// check whether product table present or not
	if !isTablePresent(tableName) {
		log.Warnf("table %s not present in DB", tableName)
		http.Error(rw, "product data not found", http.StatusNotFound)
		return
	}

	// fetch all rows from table product.
	products, err := db.Query("SELECT * FROM `product`")
	if err != nil {
		log.Warnf("unable to query table %s", tableName)
		http.Error(rw, "unable to query product data", http.StatusInternalServerError)
		return
	}

	// traverse row by row and append to array of products.
	for products.Next() {
		var name, description, url, imageUrl, price, totalReviews, creationTime, lastUpdateTime string
		err = products.Scan(&name, &description, &url, &imageUrl, &price, &totalReviews, &creationTime, &lastUpdateTime)
		if err != nil {
			panic(err.Error())
		}
		product = ProductDetails{
			URL: url,
			Product: &Product{
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

	// convert array of products struct into array of json and write response back to browser
	responseBody, err := json.Marshal(allProducts)
	if err != nil {
		log.Println("unable to convert struct to json, Reason: ", err)
		http.Error(rw, fmt.Sprintf("unable to marshal product details %v", err), http.StatusInternalServerError)
		return
	}
	_, err = rw.Write(responseBody)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Errorf("error occurred while writing response, %v ", err)
	}
}
