# amazonProductScraper
It's REST APi Which scrapes Product Name, Product image urls, Product description, Product price and Product total number of reviews for given Amazon web page URL and Stores details into a Mysql Database. 

## deploy
Verify Mysql Env
Mysql Env:
   - MYSQL_ROOT_PASSWORD=""
   - MYSQL_USER=root
   - MYSQL_DATABASE=test
```shell
git clone https://github.com/shani1998/amazonProductScraper.git
docker-compose up -d 
```

## test
```shell
curl -i -XGET -H "Content-Type: application/json" -d '{ "url": "https://www.amazon.com/PlayStation-4-Pro-1TB-Console/dp/B01LOP8EZC/"}' http://localhost:8080/scrape
curl -v -XGET -H "Content-Type: application/json" http://localhost:8081/products
```
