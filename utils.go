package main

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"math/rand"
	"net/http"
	"net/url"
)

func (t *Task) setProxy() {
	if len(config.ProxyArray) > 0 {
		proxy := config.ProxyArray[rand.Intn(len(config.ProxyArray))]

		proxyURL, err := url.Parse(proxy)
		fmt.Println(proxyURL)

		if err != nil {
			log.Printf("Error %v", err.Error())
			log.Println("[WARN] Running Proxyless")
			return
		}

		t.Client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}

		log.Printf("[INFO] Running Proxy (%v)", proxyURL.String())
	} else {
		log.Println("[WARN] Running Proxyless")
	}
}

func DeployTbBot() *tb.Bot {
	b, err := tb.NewBot(tb.Settings{
		Token: "2024776575:AAGlLRV8gbRUJbvZEpYvtIbrHXqeNqBoeP0",
	})

	if err != nil {
		log.Fatalf("Could not connect to Telegram - %v", err.Error())
		return nil
	} else {
		return b
	}
}