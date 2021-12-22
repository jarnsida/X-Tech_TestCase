package store

import (
	"log"
	"os"

	"gitgub.com/jarnsida/X-Tech_TestCase/models"
	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var (
	DBconn DbInstance
)

func ConnectDb() {
	db, err := gorm.Open(sqlite.Open("rates.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Faied to connect to DB \n", err.Error())
		os.Exit(2)
	}
	log.Println("DB connected success")
	db.Logger = logger.Default.LogMode(logger.Info)

	//Start migration
	log.Println("Running migration")
	db.AutoMigrate(&models.Income{})

	DBconn = DbInstance{Db: db}
}
