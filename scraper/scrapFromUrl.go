package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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

type products []ProductDetails

//scrape mentioned details from given url.
func scrapeFromGivenUrl(url string) (*ProductDetails, error) {
	var productName, productDesc, price, imageUrls, totalNumberOfReviews string
	var images map[string]interface{}

	log.Println("Scraping from url: ", url)

	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnHTML("#productTitle", func(e *colly.HTMLElement) {
		productName = strings.TrimSpace(e.Text)
		log.Println("name:", productName)
	})

	c.OnHTML("span.a-size-base.a-color-price", func(e *colly.HTMLElement) {
		price = strings.TrimSpace(e.Text)
		log.Println("price:", price)
	})

	c.OnHTML("#productDescription", func(e *colly.HTMLElement) {
		productDesc = strings.TrimSpace(e.Text)
		log.Println("description:", productDesc)
	})

	c.OnHTML("#acrCustomerReviewText", func(e *colly.HTMLElement) {
		totalNumberOfReviews = strings.TrimSpace(e.Text)
		log.Println("nums of reviews:", totalNumberOfReviews)
	})

	c.OnHTML("div.imgTagWrapper img", func(e *colly.HTMLElement) {
		s := e.Attr("data-a-dynamic-image")
		err := json.Unmarshal([]byte(s), &images)
		if err != nil {
			log.Println("Unable to scrape product url, Reason:", err)
		}
		fmt.Println(images)
		//join all map keys, ie image URLs
		keys := make([]string, 0, len(images))
		for k := range images {
			keys = append(keys, k)
		}
		imageUrls = "[" + strings.Join(keys, ", ") + "]"
	})
	c.Visit(url)

	productDetails := &ProductDetails{
		Url: url,
		Product: Product{
			Name:         productName,
			Description:  productDesc,
			Price:        price,
			ImageURL:     imageUrls,
			TotalReviews: totalNumberOfReviews,
		},
		CreationTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	return productDetails, nil
}

//process post request at /scrape, read payload and call scrapper function.
func processScraper(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input map[string]string
	var url string
	var ok bool

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the url field to scrape.")
	}

	err = json.Unmarshal(reqBody, &input)
	if err != nil {
		fmt.Fprintf(w, "Unable to convert data to required struct.")
	}
	if url, ok = input["url"]; !ok {
		fmt.Fprintf(w, "Unable to read url from payload, please set body like, eg: {'url':'https://myexample.com/'}")
	}
	productDetails, err := scrapeFromGivenUrl(url)

	//write response back to browser
	responseBody, err := json.Marshal(productDetails)
	if err != nil {
		log.Println("unable to make response body, Reason: ", err)
	}
	_, err = w.Write(responseBody)
	if err != nil {
		log.Println("error occurred while writing response, Reason: ", err)
	}

	//store all scraped data into mysql table using.
	url = os.Getenv("DATA_STORE_URL")
	resp, err := http.Post(url+"/insertScrapedData", "application/json", bytes.NewBuffer(responseBody))
	if err != nil {
		log.Println("unable to make request to db service, Reason:", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))

}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	mux := httprouter.New()
	mux.POST("/scrape", processScraper)
	log.Fatalln(http.ListenAndServe(":8080", mux))

}
