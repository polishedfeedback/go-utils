package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	from := flag.String("from", "USD", "Currency to send (e.g., USD)")
	to := flag.String("to", "GBP", "Currency to receive (e.g., GBP)")
	amount := flag.String("amount", "1", "Rate to send (e.g., 1)")
	flag.Parse()

	url := fmt.Sprintf("https://xe.com/currencyconverter/convert/?Amount=%s&From=%s&To=%s", *amount, *from, *to)

	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("HTTP Request failed: %v", err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Fatalf("Failed to fetch data: %s", res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}
	var sendingAmount, receivingAmount string
	doc.Find("input[aria-label='Receiving amount']").Each(func(index int, item *goquery.Selection) {
		if value, exists := item.Attr("value"); exists {
			switch index {
			case 0:
				sendingAmount = value
			case 1:
				receivingAmount = value
			}
		}
	})

	fmt.Printf("%s%s = %s%s\n", sendingAmount, *from, receivingAmount, *to)
}
