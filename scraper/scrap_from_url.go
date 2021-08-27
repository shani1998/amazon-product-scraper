package scraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"

	"github.com/shani1998/amazon-product-scraper/datastore"
)

const (
	URL = "url"
)

//---------------scrape product details from given url----------------
func scrapeFromGivenUrl(url string) *datastore.ProductDetails {
	var productName, productDesc, price, imageUrls, totalNumberOfReviews string
	var images map[string]interface{}

	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
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
		log.Println(images)

		// join all images corresponding to this product, ie image URLs
		keys := make([]string, 0, len(images))
		for k := range images {
			keys = append(keys, k)
		}
		imageUrls = "[" + strings.Join(keys, ", ") + "]"
	})

	c.Visit(url)

	productDetails := &datastore.ProductDetails{
		URL: url,
		Product: &datastore.Product{
			Name:         productName,
			Description:  productDesc,
			Price:        price,
			ImageURL:     imageUrls,
			TotalReviews: totalNumberOfReviews,
		},
		CreationTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	return productDetails

}

// ProcessScraper read amazon URL from req and return corresponding product details.
func ProcessScraper(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		log.Warnf("method %s not implemented for url %s", req.Method, req.URL.String())
		http.Error(rw, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var input map[string]string
	var url string
	var ok bool
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("unable to reead req body, %v", err)
		http.Error(rw, "unable to read req body", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(reqBody, &input)
	if err != nil {
		log.Errorf("unable to unmarshal req body to map[string]string , %v", err)
		http.Error(rw, fmt.Sprintf("unable to unmarshal req body to map[string]string , %v", err), http.StatusBadRequest)
		return
	}
	log.Infof("req body %v", input)
	if url, ok = input[URL]; !ok || len(url) == 0 {
		log.Warn("missing field 'url' in req body")
		http.Error(rw, "field 'url' must be provided", http.StatusBadRequest)
		return
	}
	productDetails := scrapeFromGivenUrl(url)
	responseBody, err := json.Marshal(productDetails)
	if err != nil {
		log.Errorf("unable to marshal product details %v ", err)
		http.Error(rw, fmt.Sprintf("unable to marshal product details %v", err), http.StatusInternalServerError)
	}
	_, err = rw.Write(responseBody)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		log.Errorf("error occurred while writing response, %v ", err)
	}

	// store all scraped data into mysql table using.
	datastore.InsertProduct(productDetails)
}
