package main

import (
	client "alien/internals"
	"alien/internals/logger"
	"alien/internals/utils"
	"encoding/json"
	"os"
	"sync"
)

func main() {
	var mutex sync.Mutex
	var wg sync.WaitGroup
	var config struct {
		Threads int    `json:"Threads"`
		Proxy   string `json:"Proxy"`
		CapKey  string `json:"CapKey"`
	}
	configData, err := os.ReadFile("./assets/config.json")
	if err != nil {
		logger.ErrorLogger.Fatal("configDataErr: ", err)
	}
	if err := json.Unmarshal(configData, &config); err != nil {
		logger.ErrorLogger.Fatal("configMarshalErr: ", err)
	}
	threadLimit := make(chan struct{}, config.Threads)
	for {
		threadLimit <- struct{}{}
		wg.Add(1)
		go func() {
			defer func() {
				<-threadLimit
				wg.Done()
			}()

			client, err := client.NewClient(config.Proxy, config.CapKey)
			if err != nil {
				return
			}
			if err := client.CreateAccount(); err != nil {
				return
			}
			if err := client.GetPromo(); err != nil {
				return
			}
			utils.Write("./assets/promos.txt", client.Promo, &mutex)
		}()
	}
}
