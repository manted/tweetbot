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
  Time Time
  Disclaimer string
  Bpi Bpi
}

type Time struct {
  Updated string
  UpdatedISO string
  Updateduk string
}

type Bpi struct {
  USD Price
}

type Price struct {
  Code string
  Rate string
  Description string
  Rate_float float64
}

func getPrice() ([]byte, error) {
  var client = &http.Client{Timeout: 10 * time.Second}
  resp, err := client.Get("https://api.coindesk.com/v1/bpi/currentprice/USD.json")
  if err != nil {
    fmt.Printf("Request Error!\n")
  }
  defer resp.Body.Close()

  bodyBytes, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Printf("Read Error!\n")
  }
  return bodyBytes, err
}

func getPriceJson(body []byte) (*APIResponse, error) {
  apiResponse := APIResponse{}
  err := json.Unmarshal(body, &apiResponse)
  if err != nil {
    fmt.Printf("Parse Error!\n")
  }
  //fmt.Println(apiResponse)
  return &apiResponse, err
}

func tweet(price string) {
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
  client.Statuses.Update("Current Bitcoin price is $" + price, nil)
}

func main() {
  priceData, requestError := getPrice()
  if requestError != nil {
    return
  }
  priceJson, jsonError := getPriceJson(priceData)
  if jsonError != nil {
    return
  }
  fmt.Printf("Price is: $%s\n", priceJson.Bpi.USD.Rate)
  tweet(priceJson.Bpi.USD.Rate)
}
