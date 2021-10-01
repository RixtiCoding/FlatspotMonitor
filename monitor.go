package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	mu sync.Mutex
	wg sync.WaitGroup
)

func (t *Task) beginMonitor(db *gorm.DB) {
	t.setProxy()

	for {
		productInfo, err := t.getProductInfo()

		if err != nil {
			log.Printf("[ERROR] - %v", err)
			log.Println("[WARN] Swapping proxy...")
			t.setProxy()
			time.Sleep(1700 * time.Millisecond)
			continue
		}

		t.checkStock(productInfo,db)
		time.Sleep(4000 * time.Millisecond)
	}
}

func (t *Task) getProductInfo() (Reply, error) {
	var reply Reply
	log.Println("[INFO] Fetching Product info...")
	req, err := http.NewRequest("GET","https://releases.flatspot.com/products.json?limit=8", nil) //https://mocki.io/v1/e07fd655-e3a0-402e-83e3-53de856ecd89

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")

	if err != nil {
		log.Printf("[ERROR] could not fetch products - %v", err.Error())
	}

	resp, err := t.Client.Do(req)

	if err != nil {
		log.Printf("[ERROR] - %v", err.Error())
	}

	err = json.NewDecoder(resp.Body).Decode(&reply)

	if err != nil {
		log.Printf("[ERROR] %v",err.Error())
	}

	defer resp.Body.Close()
	return reply, err
}


func (t *Task) checkStock(reply Reply, db *gorm.DB)  {

	for _, product := range reply.Products {

		if strings.Contains(product.Title, "coming-soon") {
			var prod DBProduct

			mu.Lock()
			db.Where("id = ?", product.Title).First(&prod)
			mu.Unlock()

			if !prod.ID.Valid {
				dbproduct := DBProduct{
					ID: sql.NullString{
						String: product.Title,
						Valid:  true,
					},
				}

				mu.Lock()
				db.Create(&dbproduct)
				log.Printf("%v FOUND! Added to db", dbproduct.ID)
				bot.Send(group, FormatPingsoon(product), tb.ModeMarkdown)
				mu.Unlock()
			}
		} else {
				var ping PingedStock
				mu.Lock()
				db.Where("id = ?", product.Title).First(&ping)
				mu.Unlock()
				for _, size := range product.Variants {
					if size.Available && !ping.ID.Valid  {
						pingstock := PingedStock {
							ID: sql.NullString{
								String: product.Title,
								Valid: true,
							},
						}
						mu.Lock()
						db.Create(&pingstock)
						mu.Unlock()
						log.Printf("%v is IN-STOCK! added to secondary db",pingstock.ID)
						bot.Send(group, FormatPingstock(product), tb.ModeMarkdown)
						break
					}
				}
			}
		}
	}




func FormatPingsoon(prod Product) string {
	var sbFinal strings.Builder

	sbFinal.WriteString(fmt.Sprintf(`
					[%v Coming Soon !](%v)

	Price: £ %v            	Region: UK


`, prod.Name, fmt.Sprintf("https://releases.flatspot.com/products/"+prod.Title), prod.Variants[0].Price))

	return sbFinal.String()

}

func FormatPingstock(prod Product) string {
	var sbFinal strings.Builder

	sbFinal.WriteString(fmt.Sprintf(`
					[%v is IN-STOCK !](%v)

	Price: £ %v            	Region: UK


`, prod.Name, fmt.Sprintf("https://releases.flatspot.com/products/"+prod.Title), prod.Variants[0].Price))

	return sbFinal.String()

}
