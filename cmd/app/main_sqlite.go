package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitgub.com/jarnsida/X-Tech_TestCase/models"
	"gitgub.com/jarnsida/X-Tech_TestCase/service"
	"gitgub.com/jarnsida/X-Tech_TestCase/store"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm/clause"
)

// Mat Ryer main() hack to catch panics
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

//run() function executes main program
func run() error {
	var ratesItem = []models.Income{}
	//ready := make(chan bool)
	//var symbolList []string
	//	var ratesList []models.RatesIn

	store.ConnectDb()
	//Start server
	app := fiber.New()
	app.Use(logger.New())

	//Goroutine for listener
	go func() {
		//Laddr
		log.Fatal(app.Listen(":3000"))
	}()

	//Goroutine gets data from request func every 30 seconds and puts into DB
	go func() {
		for {

			incomeBytes := service.Request() //get data from target JSON application

			err := json.Unmarshal(incomeBytes, &ratesItem) //Unmarshal
			if err != nil {
				log.Fatal(err)
			}

			for _, val := range ratesItem { //Update columns on conflict or create row if new

				store.DBconn.Db.Clauses(clause.OnConflict{
					Columns:   []clause.Column{{Name: "Symbol"}},
					DoUpdates: clause.AssignmentColumns([]string{"Price_24h", "Volume_24h", "Last_Trade"}),
				}).Create(&val)
			}

			time.Sleep(time.Second * time.Duration(30))
		}

	}()

	//Handler "/" uses func

	app.Get("/", service.GetRates)

	//Stop channel for Graceful Shutdown with data backup
	stop := make(chan os.Signal, 1)

	// Catch Ctrl+C, Ctrl+Z commands to stop server
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP)
	<-stop
	_ = app.Shutdown()

	return nil
}
