package main

import (
	"encoding/json"
	tb "gopkg.in/tucnak/telebot.v2"
	"io/ioutil"
	"log"
	"os"
)

// GLOBAL VARIABLES
var (
	config ProxyConfig
	bot = DeployTbBot()
	group = tb.ChatID(-1001571222914) // 579673940 //532026863 // 1001571222914

)

const IN_STOCK = true
const OUT_OF_STOCK = false

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	configFile, err := os.Open("ProxyConfig.json")

	if err != nil {
		log.Printf("[ERROR] [CONFIG] - %v", err.Error())
	}

	defer configFile.Close()
	configBytes, err := ioutil.ReadAll(configFile)

	if err != nil {
		log.Printf("[ERROR] [CONFIG] - %v", err.Error())
	}

	json.Unmarshal(configBytes, &config)

	log.Printf("[INFO] Loaded %v Proxies into config", len(config.ProxyArray))
}

func main() {
	db := Connect()
	t := createTask()
	t.beginMonitor(db)

}
