package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Types declaration for JSON
type RatesIn struct {
	Price_24h  float64 `json:"price_24h"`
	Volume_24h float64 `json:"volume_24h"`
	Last_trade float64 `json:"last_trade_price"`
}

type Outcome struct {
	Symbol string `json:"symbol"`
	RatesIn
}

// Mat Ryer main() hack to catch panics
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

//run() function executes main program
func run() error {
	var si = []Outcome{}
	var symbolList []string
	var ratesList []RatesIn

	//Start server
	app := fiber.New()
	app.Use(logger.New())

	//Goroutine for listener
	go func() {
		//Laddr
		log.Fatal(app.Listen(":3000"))
	}()

	//Goroutine gets data from request func every 30 seconds
	go func() {
		for {

			incomeBytes := request()
			err := json.Unmarshal(incomeBytes, &si)
			if err != nil {
				log.Fatal(err)
			}

			//fmt.Println(string(incomeBytes))
			//fmt.Println(si)

			for i := 0; i < len(si); i++ {
				symbolList = append(symbolList, si[i].Symbol)
				ratesList = append(ratesList, si[i].RatesIn)
			}
			r1 := symbolList
			r2 := si[0].RatesIn.Last_trade
			fmt.Printf("\n Si = %s, Last trade = %f \n", r1, r2)
			time.Sleep(time.Second * time.Duration(30))
		}

	}()

	//Handler "/" uses func
	app.Get("/", func(c *fiber.Ctx) error {

		//Marshal []Outcome struct. ??? symbol lost after RatesIn reLabling ???
		b, err := json.Marshal(si)
		if err != nil {
			log.Fatal(err)
		}

		//Insert Symbol into JSON marshaled RatesIn []bytes
		out := string(b)
		counter := 0

		for i := 0; i < len(out); i++ {

			if out[i] == '[' {
				out = out[:i+2] + strconv.Quote(string(symbolList[counter])) + ": " + out[i+1:]
				fmt.Println("found 1")
				counter++
				i = i + 6

			}

			if string(out[i]) == "}" && string(out[i+1]) == "," {
				out = out[:i+3] + strconv.Quote(string(symbolList[counter])) + ": " + out[i+2:]
				counter++
				i = i + 6

			}
		}

		return c.SendString(out) //Send right JSON to client
	})

	//Stop creates channel for Graceful Shutdown with data backup
	stop := make(chan os.Signal, 1)

	// Catch Ctrl+C, Ctrl+Z commands to stop server
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP)
	<-stop
	_ = app.Shutdown()

	return nil
}

//Redeclare JSON tags for RatesIn struct
func (r *RatesIn) MarshalJSON() ([]byte, error) {
	type alias struct {
		Price_24h  float64 `json:"price"`
		Volume_24h float64 `json:"volume"`
		Last_trade float64 `json:"last_trade"`
	}

	var a alias = alias(*r)
	return json.Marshal(&a)
}

//request func grabs body from target page
func request() []byte {
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
