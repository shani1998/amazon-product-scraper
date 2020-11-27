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
)

type Product struct {
	Name         string  `json:"name"`
	ImageURL     string  `json:"image_url"`
	Description  string  `json:"description"`
	Price        string  `json:"price"`
	TotalReviews string   `json:"total_reviews"`
}
type ProductDetails struct {
	Url             string   `json:"url"`
	Product         Product   `json:"product"`
	CreationTime    string    `json:"creation_time"`
	LastUpdatedTime string   `json:"last_updated_time"`
}

type products []ProductDetails


//connect with mysql-database and return reference to it.
func dbConn() (db *sql.DB) {
	user := os.Getenv("MYSQL_USER")
	pass := os.Getenv("MYSQL_ROOT_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	dbDriver := "mysql"
	dbUser := user
	dbPass := pass
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}


//insert payload data into the database
func insertProduct(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var productDetails ProductDetails

	//get connection reference of database
	db := dbConn()
	defer db.Close()

	//read payload
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Unable read requested payload, Reason:%v",err)
	}

	//convert body into our desire struct
	err = json.Unmarshal(reqBody,&productDetails)
	if err != nil {
		fmt.Fprintf(w,"Unable to convert data to required struct.")
	}

	//check whether table exist or not
	_, tableCheck := db.Query("select * from product;")
	if tableCheck == nil {
		fmt.Println("product table is there.")
	} else {
		fmt.Println("creating table product.")
		//Create table if doesn't exist
		stmt, err := db.Prepare("CREATE TABLE `test`.`product` (`name` VARCHAR(50) NOT NULL ,`description` LONGTEXT NOT NULL ,`url` VARCHAR(250) NOT NULL ,`imageUrl` TEXT NOT NULL ,`price` VARCHAR(20) NOT NULL ,`totalReviews` VARCHAR(20) NOT NULL ,`lastUpdatedTime` VARCHAR(250) NOT NULL ,`creationTime` VARCHAR(250) NOT NULL ,PRIMARY KEY (`url`));")
		if err != nil {
			panic(err.Error())
		}
		_, err = stmt.Exec()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Table created successfully..")
		}

	}

	//insert data into product table
	insForm, err := db.Prepare("INSERT INTO `product` (`name`, `description`, `url`, `imageUrl`, `price`, `totalReviews`, `creationTime`) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("unable to insert into database, Reason:",err)
		panic(err.Error())
	}

	//check whether body is empty or not to avoid run time error
	if (productDetails != ProductDetails{}) && (productDetails.Product != Product{}) {
		url              := productDetails.Url
		name             := productDetails.Product.Name
		desc             :=  productDetails.Product.Description
		price            := productDetails.Product.Price
		reviews      := productDetails.Product.TotalReviews
		imageUrl         := productDetails.Product.ImageURL
		creationTime := productDetails.CreationTime
		_, err := insForm.Exec(name,desc,url,imageUrl,price,reviews,creationTime)
		if err != nil {
			log.Println("unable to insert data into table test, Reason:",err)
		} else {
			log.Println("Data inserted successfully!")
		}

	} else {
		log.Println("Unable to read payload successfully!")
	}
}


//display scrapped details by url
func listProducts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var product ProductDetails
	var allProducts []ProductDetails

	//get db reference
	db := dbConn()
	defer db.Close()

	//fetch all rows from table product.
	products, err := db.Query("SELECT * FROM `product`")
	if err != nil {
		log.Println("unable to list products")
		panic(err.Error())
	}

	//traverse row by row and append to array of products.
	for products.Next() {
		var name, description, url, imageUrl, price, totalReviews, creationTime,lastUpdateTime string
		err = products.Scan(&name, &description, &url,&imageUrl,&price,&totalReviews,&creationTime,&lastUpdateTime)
		if err != nil {
			panic(err.Error())
		}
		product= ProductDetails{
			Url:             url,
			Product:         Product{
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
		log.Println("unable to make response body, Reason: ",err)
	}
	_, err = w.Write(responseBody)
	if err != nil {
		log.Println("error occurred while writing response, Reason: ",err)
	}

}


func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	mux := httprouter.New()
	mux.GET("/products", listProducts)
	mux.POST("/insertScrapedData",insertProduct)
	log.Fatalln(http.ListenAndServe(":8081",mux))

}

