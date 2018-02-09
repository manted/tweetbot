package main
import (
  "fmt"
  "net/http"
  "io/ioutil"
  "time"
  "encoding/json"
  "github.com/dghubble/go-twitter/twitter"
  "github.com/dghubble/oauth1"
  "os"
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

func tweet(bitcoinPrice string, zenStockPrice string) {
  consumerKey := os.Getenv("CONSUMER_KEY")
  consumerSecret  := os.Getenv("CONSUMER_SECRET")
  accessToken  := os.Getenv("ACCESS_TOKEN")
  accessSecret := os.Getenv("ACCESS_SECRET")

  config := oauth1.NewConfig(consumerKey, consumerSecret)
  token := oauth1.NewToken(accessToken, accessSecret)
  httpClient := config.Client(oauth1.NoContext, token)

  // Twitter client
  client := twitter.NewClient(httpClient)

  // Send a Tweet
  client.Statuses.Update("Bitcoin price is $" + bitcoinPrice + ". ZEN stock price is $" + zenStockPrice + ".", nil)
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
  if stockRequestError  != nil {
    return
  }
  fmt.Printf("Bitcoin price is $%s, ZEN stock price is $%s.\n", priceJson.Bpi.USD.Rate, zenStockPrice)
  tweet(priceJson.Bpi.USD.Rate, zenStockPrice)
}
