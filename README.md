# Amazon-product-scraper
It's REST APi Which scrapes Product Name, Product image urls, Product description, Product price and Product total number of reviews for given Amazon web page URL and Stores details into a Mysql Database. 

## Deploy
Verify Mysql Env
Mysql Env:
   - MYSQL_ROOT_PASSWORD=""
   - MYSQL_USER=root
   - MYSQL_DATABASE=test
```shell
git clone https://github.com/shani1998/amazon-product-scraper.git
docker-compose up -d 
```

## Test
```shell
curl -i -XPOST -H "Content-Type: application/json" -d '{ "url": "https://www.amazon.com/PlayStation-4-Pro-1TB-Console/dp/B01LOP8EZC/"}' http://localhost:8080/scrape
```
``` json
{
    "url": "https://www.amazon.com/PlayStation-4-Pro-1TB-Console/dp/B01LOP8EZC/",
    "product": {
        "name": "PlayStation 4 Pro 1TB Console",
        "image_url": "[https://m.media-amazon.com/images/I/6118ctEjpoL._AC_SY450_.jpg, https://m.media-amazon.com/images/I/6118ctEjpoL._AC_SX679_.jpg, https://m.media-amazon.com/images/I/6118ctEjpoL._AC_SY355_.jpg, https://m.media-amazon.com/images/I/6118ctEjpoL._AC_SX522_.jpg, https://m.media-amazon.com/images/I/6118ctEjpoL._AC_SX425_.jpg, https://m.media-amazon.com/images/I/6118ctEjpoL._AC_SX569_.jpg, https://m.media-amazon.com/images/I/6118ctEjpoL._AC_SX466_.jpg]",
        "description": "Edition:Pro 1TB\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\nPS4 Pro 4K TV GAMING & MORE The most advanced PlayStation system ever. PS4 Pro is designed to take your favorite PS4 games and add to them with more power for graphics, performance, or features for your 4K HDR TV, or 1080p HD TV. Ready to level up?  4K TV Gaming – PS4 Pro outputs gameplay to your 4K TV. Many games, like Call of Duty: WWII, Gran Turismo Sport, and more, are optimized to look stunningly sharp and detailed when played on a 4K TV with PS4 Pro. More HD Power – Turn on Boost Mode to give PS4 games access to the increased power of PS4 Pro. For HD TV Enhanced games, players can benefit from increased image clarity, faster frame rates, or more. HDR Technology – With an HDR TV, compatible PS4 games display an unbelievably vibrant and lifelike range of colors. 4K Entertainment – Stream 4K videos, movies, and shows to your PS4 Pro.  GREATNESS AWAITS 4K Entertainment requires access to a 4K compatible content streaming service, a robust internet connection, and a compatible 4K display. Enhanced for PS4 Pro Many of the biggest and best PS4 games get an additional boost from PS4 Pro enhancements that fine tune the game’s performance. From the stunning Manhattan skyline of Marvel’s Spider Man and the towering Norse mountains of God of War, to the vast plains of Red Dead Redemption 2 and the intense battlegrounds of Call of Duty: Black Ops 4, you’ll feel the power of your games unleashed wherever you see the Enhanced for PS4 Pro badge.",
        "price": "$549.99",
        "total_reviews": "11,561 ratings"
    },
    "creation_time": "2021-08-27 06:33:30",
    "last_updated_time": ""
}
```