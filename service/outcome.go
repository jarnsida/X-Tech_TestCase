package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gitgub.com/jarnsida/X-Tech_TestCase/models"
	"gitgub.com/jarnsida/X-Tech_TestCase/store"
	"github.com/gofiber/fiber/v2"
)

//Outcome serializer

type RatesOut struct {
	Price     float64 `json:"price"`
	Volume    float64 `json:"volume"`
	LastTrade float64 `json:"last_trade"`
}

type Outcome struct {
	Label    string `json:"Symbol"`
	RatesOut struct {
		Price     float64 `json:"price"`
		Volume    float64 `json:"volume"`
		LastTrade float64 `json:"last_trade"`
	}
}

func GetRates(c *fiber.Ctx) error {

	rates := []models.Income{}

	store.DBconn.Db.Find(&rates)

	responseItems := []Outcome{}
	fmt.Println(c.JSON(rates))
	for _, rate := range rates {
		responseItem := CreateResresponseItems(rate)
		responseItems = append(responseItems, responseItem)

	}
	//Marshal []Outcome struct. ??? symbol lost after RatesIn reLabling ???
	responseBytes, err := json.Marshal(responseItems)
	if err != nil {
		log.Fatal(err)
	}
	responseString := tagNameDelete(string(responseBytes))

	return c.Status(200).SendString(responseString)

}

func CreateResresponseItems(rates models.Income) Outcome {

	ratesOutItem := RatesOut{Price: rates.Price_24h, Volume: rates.Volume_24h,
		LastTrade: rates.Volume_24h}

	return Outcome{Label: rates.Symbol, RatesOut: ratesOutItem}
}

//Request func execute request to target site
func Request() []byte {
	res, err := http.Get("https://api.blockchain.com/v3/exchange/tickers")
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	return body
}

//tagNameDelete remove Symbol and RetesIn tags from JSON
func tagNameDelete(json string) string {
	s := strings.NewReplacer(strconv.Quote("Symbol")+":", "", ","+strconv.Quote("RatesOut"), "")
	json = s.Replace(json)

	return json
}
