package db

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	dbUser = "kansas"
	dbPass = "pass1234"
	dbName = "cross-pollinators-db"
)

func ConnectDB() error {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPass, dbName)
	db, err := gorm.Open("postgres", dbinfo)

	if err != nil {
		fmt.Printf("Failed to open connection to postgres database: %+v \n", err)
		return err
	}

	log.Printf("Cross Pollinators Service connected to DB!")

	defer db.Close()

	return nil
}
