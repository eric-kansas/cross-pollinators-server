package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/eric-kansas/cross-pollinators-server/database/models"
)

const (
	dbUser = "kansas"
	dbPass = "pass1234"
	dbName = "cross-pollinators-db"
)

func main() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPass, dbName)

	db, err := gorm.Open("postgres", dbinfo)

	if err != nil {
		log.Fatalf("Failed to open connection to postgres database: %+v \n", err)
	}

	log.Printf("Connected to Cross Pollinators DB!")

	db.AutoMigrate(&models.User{}, &models.Interest{}, &models.Project{}, &models.Comment{})

	defer db.Close()
}
