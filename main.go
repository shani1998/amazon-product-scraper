package main

import (
	"net"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/shani1998/amazon-product-scraper/datastore"
	"github.com/shani1998/amazon-product-scraper/scraper"
)

func main() {
	conn, err := datastore.DBConn()
	if err != nil {
		log.Fatalf("unable to initialize DB client")
		return
	}
	defer conn.Close() // TODO add readiness and liveness check for database

	scraperMux := http.NewServeMux()
	scraperMux.HandleFunc("/scrape", scraper.ProcessScraper)
	scraperMux.HandleFunc("/products", datastore.ListProducts)
	address := net.JoinHostPort("", os.Getenv("SCRAPER_HOST_PORT"))
	log.Infof("Starting scraper server at address [%s]", address)
	if err := http.ListenAndServe(address, scraperMux); err != nil {
		log.Fatalf("Failed to serve metrics at [%s]: %v", address, err)
	}
}
