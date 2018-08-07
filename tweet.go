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

type Quote struct {
	Close         float32
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

func getZenQuote(client *http.Client) ([]byte, error) {
	resp, err := client.Get("https://api.iextrading.com/1.0/stock/zen/quote?displayPercent=true")
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

func getQuoteJson(body []byte) (*Quote, error) {
	previousStock := Quote{}
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

	zenQuoteData, previousStockequestError := getZenQuote(client)
	if previousStockequestError != nil {
		return
	}
	zenQuoteJson, jsonError := getQuoteJson(zenQuoteData)
	if jsonError != nil {
		return
	}
	close := zenQuoteJson.Close
	closeStr := fmt.Sprintf("%.2f", close)

	change := zenQuoteJson.Change
	changeStr := fmt.Sprintf("%.2f", change)
	if change >= 0 {
		changeStr = "+" + changeStr
	}

	changePercent := zenQuoteJson.ChangePercent
	changePercentStr := fmt.Sprintf("%.2f", changePercent)
	if changePercent >= 0 {
		changePercentStr = "+" + changePercentStr
	}
	fmt.Printf("Bitcoin: $%s. ZEN: $%s %s (%s%%)", priceJson.Bpi.USD.Rate, closeStr, changeStr, changePercentStr)
	tweet(priceJson.Bpi.USD.Rate, closeStr, changeStr, changePercentStr)
}
