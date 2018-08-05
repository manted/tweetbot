package main

import (
	"encoding/json"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type APIResponse struct {
	Bpi Bpi
}

type Bpi struct {
	USD Price
}

type Price struct {
	Rate string
}

type PreviousStock struct {
	Change        float32
	ChangePercent float32
}

func getBitcoinPrice(client *http.Client) ([]byte, error) {
	resp, err := client.Get("https://api.coindesk.com/v1/bpi/currentprice/USD.json")
	if err != nil {
		fmt.Printf("Bitcoin Request Error!\n")
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Bitcoin Read Error!\n")
	}
	return bodyBytes, err
}

func getBitcoinPriceJson(body []byte) (*APIResponse, error) {
	apiResponse := APIResponse{}
	err := json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Printf("Bitcoin Response Parse Error!\n")
	}
	//fmt.Println(apiResponse)
	return &apiResponse, err
}

func getZenStockPrice(client *http.Client) (string, error) {
	resp, err := client.Get("https://api.iextrading.com/1.0/stock/zen/price")
	if err != nil {
		fmt.Printf("Stock Request Error!\n")
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Stock Read Error!\n")
	}
	return string(bodyBytes), err
}

func getZenPreviousStockPrice(client *http.Client) ([]byte, error) {
	resp, err := client.Get("https://api.iextrading.com/1.0/stock/zen/previous")
	if err != nil {
		fmt.Printf("Previous Stock Request Error!\n")
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Previous Stock Read Error!\n")
	}
	return bodyBytes, err
}

func getPreviousStockJson(body []byte) (*PreviousStock, error) {
	previousStock := PreviousStock{}
	err := json.Unmarshal(body, &previousStock)
	if err != nil {
		fmt.Printf("Previous Stock Response Parse Error!\n")
	}
	return &previousStock, err
}

func tweet(bitcoinPrice string, zenStockPrice string, change string, changePercent string) {
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessToken := os.Getenv("ACCESS_TOKEN")
	accessSecret := os.Getenv("ACCESS_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Send a Tweet
	client.Statuses.Update("Bitcoin: $"+bitcoinPrice+". ZEN: $"+zenStockPrice+" "+change+" ("+changePercent+"%)", nil)
}

func main() {
	var client = &http.Client{Timeout: 10 * time.Second}
	priceData, requestError := getBitcoinPrice(client)
	if requestError != nil {
		return
	}
	priceJson, jsonError := getBitcoinPriceJson(priceData)
	if jsonError != nil {
		return
	}
	zenStockPrice, stockRequestError := getZenStockPrice(client)
	if stockRequestError != nil {
		return
	}
	zenPreviousStockData, previousStockequestError := getZenPreviousStockPrice(client)
	if previousStockequestError != nil {
		return
	}
	zenPreviousStockJson, jsonError := getPreviousStockJson(zenPreviousStockData)
	if jsonError != nil {
		return
	}
	fmt.Printf("Bitcoin: $%s. ZEN: $%s %.2f (%.2f%%)", priceJson.Bpi.USD.Rate, zenStockPrice, zenPreviousStockJson.Change, zenPreviousStockJson.ChangePercent)
	tweet(priceJson.Bpi.USD.Rate, zenStockPrice, fmt.Sprintf("%.2f", zenPreviousStockJson.Change), fmt.Sprintf("%.2f", zenPreviousStockJson.ChangePercent))
}
