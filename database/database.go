package database

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

func Connect() error {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPass, dbName)
	db, err := gorm.Open("postgres", dbinfo)

	if err != nil {
		fmt.Printf("Failed to open connection to postgres database: %+v \n", err)
		return err
	}

	log.Printf("Cross Pollinators Service connected to DB!!")

	user := User{
		Email:       "goat.man@gmail.com",
		Password:    "bcryptthisshit",
		Description: "some description",
		Interests:   nil,
		Projects:    nil,
		Location:    "wat",
	}

	log.Printf("User: %+v", db.NewRecord(user))
	db.Create(&user)
	log.Printf("User: %+v", db.NewRecord(user))

	defer db.Close()

	return nil
}
